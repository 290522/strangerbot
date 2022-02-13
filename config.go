package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigParser struct {
	Telegram             Telegram
	WhiteList            WhiteList
	MysqlDB              Database
	RedisConf            RedisConf
	OTPConf              OTPConf
	EmailOTP             EmailOTP
	VerifyProfileConf    VerifyProfileConf
	VerifyMatchingConf   VerifyMatchingConf
	FemaleMatchRateLimit DefaultMatchRateLimitConf
	MaleMatchRateLimit   DefaultMatchRateLimitConf
}

type WhiteList struct {
	WhiteDomainEnabled bool
	WhiteDomain        string
	WhiteEmailEnabled  bool
}

type Telegram struct {
	Namespace           string
	Secret              string
	BotKey              string
	ChangeGenderEnabled bool
}

type RedisConf struct {
	Host           string
	Port           int
	Username       string
	Password       string
	MaxActive      int
	MaxIdle        int
	TimeoutSeconds int
	KeyPrefix      string
}

type OTPConf struct {
	OTPTTL         int64
	OTPMaxAttempts int
	OTPMaxLen      int
}

type EmailOTP struct {
	Subject      string
	Template     string
	TemplateType string
	Config       string
}

type Database struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	Charset         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
}

type VerifyProfileConf struct {
	ProfileQuestionId    int64
	VerifyOptionId       int64
	RemainVerifyOptionId int64
}

type VerifyMatchingConf struct {
	MatchingQuestionId int64
	VerifiedOptionId   int64
	UnverifiedOptionId int64
	AnyOptionId        int64
}

type DefaultMatchRateLimitConf struct {
	OptionId            int64
	RateLimitEnabled    bool
	RateLimitUnit       string
	RateLimitUnitPeriod int64
	MatchPerRate        int64
}

// 加载配置
func WatchConfig(changeConfig chan struct{}, filename string) error {

	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(AppPath, "config"))
	viper.SetConfigName(filename)
	viper.SetConfigType("toml")

	if err = viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		changeConfig <- struct{}{}
	})

	return nil
}

// 重新从viper载入与校验
func LoadConfig() (*ConfigParser, error) {

	var config ConfigParser
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
