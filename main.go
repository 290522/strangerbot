package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Machiel/telegrambot"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strangerbot/otpgateway"
	"strangerbot/otpgateway/smtp"
	"strangerbot/repository"
	"strangerbot/vars"
)

var (
	telegram        telegrambot.TelegramBot
	telegramBot     *tgbotapi.BotAPI
	db              *sqlx.DB
	emptyOpts       = telegrambot.SendMessageOptions{}
	commandHandlers = []commandHandler{
		commandDisablePictures,
		commandHelp,
		commandStart,
		commandStop,
		commandReport,
		commandSetup,
		commandMessage,
	}
	startJobs            = make(chan int64, 10000)
	updatesQueue         = make(chan *tgbotapi.Update, 10000)
	endConversationQueue = make(chan EndConversationEvent, 10000)
	updateMap            = NewIdRecorder()
	messageMap           = NewIdRecorder()
	stopped              = false
)

func main() {

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// flag params
	configFilename := flag.String("cfg", "app.toml", "config filename,didn't need input file extend name.")
	if len(*configFilename) == 0 {
		panic(errors.New("config filename is empty"))
	}

	log.Println("config file: ", *configFilename)

	var err error

	log.Println("Starting...")

	// config init
	upCfg := make(chan struct{})
	err = WatchConfig(upCfg, *configFilename)
	if err != nil {
		panic(err)
	}
	defer close(upCfg)

	// load config
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	// init white cfg
	vars.ChangeGenderEnabled = cfg.Telegram.ChangeGenderEnabled
	vars.WhiteDomainEnabled = cfg.WhiteList.WhiteDomainEnabled
	vars.WhiteDomain = cfg.WhiteList.WhiteDomain
	vars.WhiteEmailEnabled = cfg.WhiteList.WhiteEmailEnabled

	vars.VerifyProfileQuestionId = cfg.VerifyProfileConf.ProfileQuestionId
	vars.VerifyOptionId = cfg.VerifyProfileConf.VerifyOptionId

	vars.MatchingQuestionId = cfg.VerifyMatchingConf.MatchingQuestionId
	vars.MatchingVerifiedOptionId = cfg.VerifyMatchingConf.VerifiedOptionId
	vars.MatchingUnverifiedOptionId = cfg.VerifyMatchingConf.UnverifiedOptionId
	vars.MatchingAnyOptionId = cfg.VerifyMatchingConf.AnyOptionId

	vars.FemaleMatchRateLimit = vars.MatchRateLimit{
		OptionId:            cfg.FemaleMatchRateLimit.OptionId,
		RateLimitEnabled:    cfg.FemaleMatchRateLimit.RateLimitEnabled,
		RateLimitUnit:       cfg.FemaleMatchRateLimit.RateLimitUnit,
		RateLimitUnitPeriod: cfg.FemaleMatchRateLimit.RateLimitUnitPeriod,
		MatchPerRate:        cfg.FemaleMatchRateLimit.MatchPerRate,
	}
	vars.MaleMatchRateLimit = vars.MatchRateLimit{
		OptionId:            cfg.MaleMatchRateLimit.OptionId,
		RateLimitEnabled:    cfg.MaleMatchRateLimit.RateLimitEnabled,
		RateLimitUnit:       cfg.MaleMatchRateLimit.RateLimitUnit,
		RateLimitUnitPeriod: cfg.MaleMatchRateLimit.RateLimitUnitPeriod,
		MatchPerRate:        cfg.MaleMatchRateLimit.MatchPerRate,
	}

	// init gorm db
	if err := InitDB(cfg.MysqlDB); err != nil {
		panic(err)
	}

	_ = repository.InitRepository(DB)

	// init sqlx db
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", cfg.MysqlDB.Username, cfg.MysqlDB.Password, cfg.MysqlDB.Host, cfg.MysqlDB.Port, cfg.MysqlDB.DBName)
	db, err = sqlx.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	// init telegram bot
	telegram = telegrambot.New(cfg.Telegram.BotKey)
	telegramBot, err = tgbotapi.NewBotAPI(cfg.Telegram.BotKey)
	if err != nil {
		panic(err)
	}

	// load tpl
	otpTpl, err := otpgateway.LoadProviderTemplates(cfg.EmailOTP.Template, cfg.EmailOTP.Subject)
	if err != nil {
		panic(err)
	}
	_ = otpTpl

	// init smtp
	_, err = smtp.InitEmailer([]byte(cfg.EmailOTP.Config))
	if err != nil {
		panic(err)
	}

	// init store
	store := otpgateway.NewRedisStore(otpgateway.RedisConf{
		Host:      cfg.RedisConf.Host,
		Port:      cfg.RedisConf.Port,
		Username:  cfg.RedisConf.Username,
		Password:  cfg.RedisConf.Password,
		MaxActive: cfg.RedisConf.MaxActive,
		MaxIdle:   cfg.RedisConf.MaxIdle,
		Timeout:   time.Duration(cfg.RedisConf.TimeoutSeconds) * time.Second,
		KeyPrefix: cfg.RedisConf.KeyPrefix,
	})

	// init otp master
	OTPMasterIns = NewOTPMaster(
		cfg.Telegram.Namespace,
		store,
		otpTpl,
		smtp.GetEmailer(),
		time.Duration(cfg.OTPConf.OTPTTL)*time.Second,
		cfg.OTPConf.OTPMaxAttempts,
		cfg.OTPConf.OTPMaxLen,
	)

	if OTPMasterIns == nil {
		panic(errors.New("OTP Master is nil"))
	}

	// start open goroutine to listen
	var wg sync.WaitGroup

	wg.Add(1)
	go func(jobs <-chan int64) {
		defer wg.Done()
		log.Println("Starting match user job")
		matchUsers(jobs)
	}(startJobs)

	for j := 0; j < 1; j++ {
		wg.Add(1)
		go func(jobs chan<- int64) {
			defer wg.Done()
			log.Println("Started load available user job")
			loadAvailableUsers(jobs)
		}(startJobs)
	}

	var workerWg sync.WaitGroup
	for i := 0; i < 3; i++ {
		workerWg.Add(1)
		go func(queue <-chan *tgbotapi.Update) {
			defer workerWg.Done()
			log.Println("Started a message worker...")
			updateWorker(queue)
		}(updatesQueue)
	}

	for x := 0; x < 1; x++ {
		wg.Add(1)

		go func(queue <-chan EndConversationEvent) {
			defer wg.Done()
			log.Println("Started end convo worker...")
			endConversationWorker(queue)
		}(endConversationQueue)
	}

	var receiverWg sync.WaitGroup
	receiverWg.Add(1)
	go func() {
		defer receiverWg.Done()
		log.Println("Started update worker")

		var offset int

		for {
			//log.Println("Requesting updates")
			offset = processUpdates(offset)
			//log.Println("Request completed")
			time.Sleep(1 * time.Second)

			if stopped {
				break
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		done <- true
	}()

	<-done

	log.Printf("Stopping...")

	stopped = true

	receiverWg.Wait()

	close(updatesQueue)

	workerWg.Wait()

	close(startJobs)
	close(endConversationQueue)

	log.Printf("Waiting for goroutines to stop...")

	wg.Wait()

	log.Printf("Closed...")
}

func loadAvailableUsers(startJobs chan<- int64) {

	for {

		u, err := retrieveAllAvailableUsers()

		if err != nil {
			log.Printf("Error retrieving everyone available: %s", err)
		} else {
			for _, x := range u {
				startJobs <- x.ChatID
			}
		}

		time.Sleep(10 * time.Second)
	}

}

// User holds user data
type User struct {
	ID                     int64         `db:"id"`
	ChatID                 int64         `db:"chat_id"`
	Available              bool          `db:"available"`
	LastActivity           time.Time     `db:"last_activity"`
	MatchChatID            sql.NullInt64 `db:"match_chat_id"`
	RegisterDate           time.Time     `db:"register_date"`
	PreviousMatch          sql.NullInt64 `db:"previous_match"`
	AllowPictures          bool          `db:"allow_pictures"`
	BannedUntil            NullTime      `db:"banned_until"`
	Gender                 int           `db:"gender"`
	Tags                   string        `db:"tags"`
	MatchMode              int           `db:"match_mode"`
	Email                  string        `db:"email"`
	IsVerify               bool          `db:"is_verify"`
	IsWaitInputEmail       bool          `db:"is_wait_input_email"`
	CustomRateLimitEnabled bool          `db:"custom_rate_limit_enabled"`
	RateLimitUnit          string        `db:"rate_limit_unit"`
	RateLimitUnitPeriod    int64         `db:"rate_limit_unit_period"`
	MatchPerRate           int64         `db:"match_per_rate"`
}

func (u User) IsProfileFinish() bool {
	if u.Gender == 0 || u.MatchMode < 0 || u.Tags == "" {
		return false
	}
	return true
}

func (u User) GetNeedFinishProfile() string {

	field := make([]string, 0, 3)

	if u.Gender == 0 {
		field = append(field, "Gender")
	}
	if u.MatchMode < 0 {
		field = append(field, "Gender Preference")
	}
	if u.Tags == "" {
		field = append(field, "Goal")
	}

	return strings.Join(field, ",")
}

func retrieveUser(chatID int64) (User, error) {
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE chat_id = ?", chatID)
	return u, err
}

func retrieveOrCreateUser(chatID int64) (User, error) {
	var u User
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE chat_id = ?", chatID)

	if err != nil {
		return u, err
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO users(chat_id, available, last_activity, register_date, allow_pictures) VALUES (?, ?, ?, ?, 1)", chatID, false, time.Now(), time.Now())

		if err != nil {
			return u, err
		}

		telegram.SendMessage(chatID, `Welcome to the Cupid SG Bot! :D

                To configure your profile:

                /setup

                To start a chat, enter:

                /start

                If you feel like ending the conversation, enter:

                /end

                If you want another chat partner, type /start again after typing /end!

                Have fun!`, emptyOpts)
	}

	return retrieveUser(chatID)
}

func updateLastActivity(id int64) {
	db.Exec("UPDATE users SET last_activity = ? WHERE id = ?", time.Now(), id)
}

func updateGender(id int64, gender int) {
	_, err := db.Exec("UPDATE users SET gender = ? WHERE id = ?", gender, id)
	if err != nil {
		log.Println("update gender info error", err.Error())
	}
}

func updateMathMode(id int64, mathMode int) {
	_, err := db.Exec("UPDATE users SET match_mode = ? WHERE id = ?", mathMode, id)
	if err != nil {
		log.Println("update match mode error", err.Error())
	}
}

func updateTags(id int64, tags string) {
	_, err := db.Exec("UPDATE users SET tags = ? WHERE id = ?", tags, id)
	if err != nil {
		log.Println("update tags error", err.Error())
	}
}

func retrieveAllAvailableUsers() ([]User, error) {
	var u []User
	err := db.Select(&u, "SELECT * FROM users WHERE available = 1 AND match_chat_id IS NULL")
	return u, err
}

func retrieveAvailableUsers(c int64, user User) ([]User, error) {
	var u []User

	sql := `SELECT * FROM users WHERE (gender > 0 AND tags!="" AND match_mode > -1) AND chat_id != ? AND available = 1 AND match_chat_id IS NULL`

	switch user.MatchMode {
	case 1:
		sql = sql + fmt.Sprintf(" AND gender = 1 AND (match_mode = 0 OR match_mode = %d)", user.Gender)
	case 2:
		sql = sql + fmt.Sprintf(" AND gender = 2 AND (match_mode = 0 OR match_mode = %d)", user.Gender)
	default:
		sql = sql + fmt.Sprintf(" AND (match_mode = 0 OR match_mode = %d)", user.Gender)
	}

	if user.Tags != "" {
		sql = sql + ` AND tags = "` + user.Tags + `"`
	}

	err := db.Select(&u, sql, c)
	//err := db.Select(&u, "SELECT * FROM users WHERE chat_id != ? AND available = 1 AND match_chat_id IS NULL", c)
	return u, err
}

func shuffle(a []User) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func handleMessage(message *tgbotapi.Message) {

	u, err := retrieveOrCreateUser(message.Chat.ID)

	if err != nil {
		log.Printf("retrieveOrCreateUser err: %s", err.Error())
		return
	}

	if u.BannedUntil.Valid && time.Now().Before(u.BannedUntil.Time) {
		date := u.BannedUntil.Time.Format("02 January 2006")
		response := fmt.Sprintf("You are banned until %s", date)
		_, err := telegram.SendMessage(message.Chat.ID, response, emptyOpts)
		if err != nil {
			log.Printf("handleMessage telegram.SendMessage err: %s", err.Error())
			return
		}
		return
	}

	sendToHandler(u, message)

	updateLastActivity(u.ID)

}

func sendToHandler(u User, message *tgbotapi.Message) {

	log.Printf("msg_id: %d sendToHandler", message.MessageID)

	for _, handler := range commandHandlers {

		res := handler(u, message)
		if res {
			return
		}

	}

}

func processUpdates(offset int) int {

	log.Printf("Fetching with offset %d", offset)

	updates, err := telegramBot.GetUpdates(tgbotapi.UpdateConfig{
		Offset:  offset,
		Limit:   100,
		Timeout: 20,
	})

	if err != nil {
		log.Printf("GetUpdates err: %s", err.Error())
		return offset
	}

	return handleUpdates(updates, offset)

}

func handleUpdate(update *tgbotapi.Update) {

	if update.Message != nil {

		if !messageMap.IsSent(update.Message.MessageID) {
			if messageMap.SetSent(update.Message.MessageID) {
				handleMessage(update.Message)
			} else {
				log.Printf("message id:%d is handled", update.Message.MessageID)
			}

		} else {
			log.Printf("message id:%d is handled", update.Message.MessageID)
		}

	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update.CallbackQuery)
	}

}

func updateWorker(updates <-chan *tgbotapi.Update) {

	for update := range updates {

		if !updateMap.IsSent(update.UpdateID) {
			if updateMap.SetSent(update.UpdateID) {
				handleUpdate(update)
			} else {
				log.Printf("update id:%d is handled", update.UpdateID)
			}
		} else {
			log.Printf("update id:%d is handled", update.UpdateID)
		}

	}

}

func handleUpdates(updates []tgbotapi.Update, offset int) int {

	for i, update := range updates {

		if update.UpdateID >= offset {
			if update.UpdateID%1000 == 0 {
				log.Printf("Update ID: %d", update.UpdateID)
			}
			offset = update.UpdateID + 1
		}

		updatesQueue <- &updates[i]
	}

	return offset
}
