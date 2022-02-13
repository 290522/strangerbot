package model

import (
	"database/sql"
)

type User struct {
	ID                     int64         `db:"id"`
	ChatID                 int64         `db:"chat_id"`
	Available              bool          `db:"available"`
	MatchChatID            sql.NullInt64 `db:"match_chat_id"`
	PreviousMatch          sql.NullInt64 `db:"previous_match"`
	AllowPictures          bool          `db:"allow_pictures"`
	CustomRateLimitEnabled bool          `db:"custom_rate_limit_enabled"`
	RateLimitUnit          string        `db:"rate_limit_unit"`
	RateLimitUnitPeriod    int64         `db:"rate_limit_unit_period"`
	MatchPerRate           int64         `db:"match_per_rate"`
}

func (p *User) TableName() string {
	return "users"
}
