# Overview
[![Twitter](https://img.shields.io/badge/author-%40MachielMolenaar-blue.svg)](https://twitter.com/MachielMolenaar)

This is the source of StrangerBot, the bot on Telegram that matches two random
users and allows them to chat with each other.

# License
StrangerBot is licensed under the Apache 2.0 license.

# Depends

* Golang
* Mysql
* Redis

# Installation
Currently there are no binaries available as direct download yet, you should
build it yourself using Go.

If you have go installed, you can install strangerbot like this:

`go get -u github.com/Machiel/strangerbot`

# How Run?

1. cd /root/go/src/strangerbot dir
2. run command "go build ."
3. edit config file, run command "vi ./config/app.toml"
4. run command "nohup ./strangerbot > /dev/null 2>&1 &"
5. run command "ps -ef | grep strangerbot" check process is running

# config/app.toml

```toml
[Telegram]
# Application Name
Namespace = "MyApp"
# Secret Token for OTP,Random 32-bit string e.g: blwPphsuyer3O1QgXe0sy1d3M0ZXzRZl
Secret = "mySecretToken"
# Telegram Bot Key Input Here
BotKey = "Telegram Bot Key"

[WhiteList]
# white email domain enabled
WhiteDomainEnabled = true
# white email domain list
WhiteDomain = "@qq.com,@hotmail.com"
# white email list enabled , email list read from db
WhiteEmailEnabled = true

# Bot Program DB Connect Info
[MysqlDB]
Host = "127.0.0.1"
Port = 3306
Username = "root"
Password = "root"
DBName = "strangerbot"
# The following configuration can be left unchanged
Charset = "utf8"
MaxOpenConns = 1000
MaxIdleConns = 1000
ConnMaxLifetime = 10

# OTP Redis Conf
[RedisConf]
Host = "127.0.0.1"
Port = "6379"
Username = ""
Password = ""
MaxActive = 100
MaxIdle = 10
TimeoutSeconds = 10
KeyPrefix = "OTP"

[OTPConf]
# otp time to live seconds
OTPTTL = 300
# otp max attempts
OTPMaxAttempts = 500
# otp len
OTPMaxLen = 6

[EmailOTP]
# email subject
Subject = "MyApp Email verification"
# email send template
Template = "static/smtp.tpl"
# email template type default:html
TemplateType = "html"
# email smtp config
Config = '''
    {
        "Host": "smtp.sendgrid.net",
        "Port": 25,
        "AuthProtocol": "",
        "AuthProtocol": "cram",
        "User": "smtp-user",
        "Password": "smtp-password",
        "FromEmail": "OTP verification <yoursite@yoursite.com>",
        "MaxConns": 10,
        "Sendtimeout": 5
    }
'''

[VerifyProfileConf]
# profile question id
ProfileQuestionId = 8
# verify option id
VerifyOptionId = 27

[VerifyMatchingConf]
# matching verify question id
MatchingQuestionId = 0
# verified option id
VerifiedOptionId = 0
# unverified option id
UnverifiedOptionId = 0
# anything option id
AnyOptionId = 0

[FemaleMatchRateLimit]
# female option id
OptionId = 2
# rate limit endabled
RateLimitEnabled = true
# rate limit unit only support "day"
RateLimitUnit = "day"
# rate limit unit period, don't change this value.
RateLimitUnitPeriod = 7
# match per rate
MatchPerRate = 1

[MaleMatchRateLimit]
# male option id
OptionId = 1
# rate limit endabled
RateLimitEnabled = true
# rate limit unit only support "day"
RateLimitUnit = "day"
# rate limit unit period, don't change this value.
RateLimitUnitPeriod = 7
# match per rate
MatchPerRate = 2
```

## Menu configuration

| Field | Default | Description |
| --- | --- | --- |
| parent_id | 0 | parent menu id 0.top menu >0.sub menu |
| target_type | 1 | target 1.menu to menu 2. menu to question |
| name | | menu name(title) |
| question_id | | if target_type is 2, this field must be set question id |
| sort | 0 | order by sort asc |
| row_index | 0 | same row index will be inline |
| helper_title | | menu helper title |
| helper_text | | menu helper text |
| is_back_enabled | 1 | show back button 0.false 1.true |
| back_button_text | Back | back button text |

## Question configuration

| Field | Default | Description |
| --- | --- | --- |
| scene_type | 0 | scene type 1.profile question 2.matching question |
| helper_title | | question helper title |
| title | | question title |
| helper_text | | helper text |
| frontend_type | 1 | frontend type 1.select 2.multi select |
| max_multi_len | 0 | multi select max choose option length, 0 is not limits. |
| sort | 0 | order by sort asc, this field not used. |
| matching_mode | 1 | order by sort asc, this field not used. |
| matching_question_id | 0 | if scene is matching , it is origin question id, if 0 it will not support matching |

## Option configuration

| Field | Default | Description |
| --- | --- | --- |
| question_id | 0 | question id it form form_question table |
| option_type | 1 | option type 1.value option, now only have one type. |
| matching_option_id | 0 | match option id |
| label | | option label. |
| value | | option value. |
| is_matching_any | 0 | 0.false 1.true |
| sort | 0 | order by sort asc |
| row_index | 0 | same row index will be inline |

# WhiteList

## White Domain

1. edit config file
```
vi ./config/app.toml
```
2. set WhiteDomainEnabled as true
```
WhiteDomainEnabled = true
```
3. edit white email domain
```
WhiteDomain = "@a.com,@b.com,@c.com"
```
> multi domain use Comma separated

## White Email

1. edit config file
```
vi ./config/app.toml
```
2. set WhiteEmailEnabled as true
```
WhiteEmailEnabled = false
```
3. add white email domain mysql table.
    1. login mysql console.
    2. choose db: ```use bot2;```
    3. run INSERT command
        ```
            INSERT INTO `bot2`.`e_email_whitelist` (`id`, `email`) VALUES (1, 'test@a.com');
        ```
       > test@a.com is white email.

# Match Rate Limit

1. set up config file

```toml
[FemaleMatchRateLimit]
# female option id
OptionId = 2
# rate limit endabled
RateLimitEnabled = true
# rate limit unit only support "day"
RateLimitUnit = "day"
# rate limit unit period, don't change this value.
RateLimitUnitPeriod = 7
# match per rate
MatchPerRate = 1

[MaleMatchRateLimit]
# male option id
OptionId = 1
# rate limit endabled
RateLimitEnabled = true
# rate limit unit only support "day"
RateLimitUnit = "day"
# rate limit unit period, don't change this value.
RateLimitUnitPeriod = 7
# match per rate
MatchPerRate = 2
```

2. how custom user match rate limit?

    1. open `users` database table
    2. set `custom_rate_limit_enabled` value as `1`
    3. set `match_per_rate`