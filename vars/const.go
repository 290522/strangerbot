package vars

var (
	// debug or prod
	RUN_MODE = "prod"
)

const (
	StartMatchingMessage    = `Looking for another cool user to match you with... Hold on! (This may take a while, maybe even days as we strive for a perfect match! Keep your notifications on!) **NOTE: If you send anything illegal here, your data will be handed over to the authorities. You're anonymous only until you break the rules. Enter /report (with a reason) to report a user; e.g. /report advertising . If a chat with a user you want to report has already ended, do not start a new chat—immediately contact the admin.`
	NotProfileFinishMessage = "Press /setup to ensure everything has been configured in your profile and match settings before pressing /start"
	TopMenuMessage          = `Set up both your profile and match settings here! You'll need to fully fill both before you can enter /start to join the queue. 💬`
	HelpMessage             = `Help:

Use /start to start looking for a conversation partner; once you're matched you can use /end to end the conversation.

Use /report to report a user, use it as follows:
/report <reason>

Use /nopics to disable receiving photos, and /nopics if you want to enable it again.

HEAD OVER to @unichatbotchannel for rules, updates, announcements or info on how you can support the bot!

Note that voice messages and tele bubbles are currently disabled.

If you require any help, feel free to contact @aaldentnay !`
	MatchedMessage           = `You have been matched, have fun! Your match is %s and is open to %s!`
	QuestionCallbackTest     = `This question can only have %d options at most`
	BanMessage               = `You are banned until %s`
	RegTipMessage            = `If you have a student email, you can get verified here! This is meaningful for students who want to match only with students and also prevents abuse. Otherwise, feel free to select 'no student email' :)`
	NeedInputEmailMessage    = `What is your student email address? (enter this in lower-case); not all student emails are recognised yet, but give yours a shot!`
	SendEmailCodeMessage     = `We've e-mailed you a code. Please check your e-mail (check your junk/spam box too! 📪) and enter the code here to complete the verification.`
	OTPNoExistsMessage       = `Code does not exist, please re-enter email or code.`
	OTPFailMessage           = `Code does not match, please re-enter email or code.`
	OTPSuccessMessage        = `Successfully registered! Please click /setup if you are not done setting up for the other options (including match settings), or /start to match.`
	NotAccessRegisterMessage = `Sorry, you are not authorised to register.`
	RetryInputEmailMessage   = `Mail service send failed, Please re-enter email.`
	RateLimitMessage         = `You have hit your weekly limit! Wait a few days before trying again. Make a donation and contact @aaldentnay to be a paid user (you’ll be  able to increase the number of times you can press /start !)`
	InternalErrorMessage     = `An error occurred inside the robot`
	ChangeGenderErrorMessage = `Gender change is not supported`
)

// emoji UTF-8 from https://apps.timwhitlock.info/emoji/tables/unicode
const (
	ChooseMark = "\xE2\x9C\x85"
)

// set goals question id, it work for MatchedMessage.
const (
	GoalsQuestionId int64 = 7
)
