package model

import "strangerbot/repository/gorm_global"

type MatchedDetail struct {
	gorm_global.ColumnCreateModifyDeleteTime
	ChatId      int64 `db:"chat_id"`
	MatchChatId int64 `db:"match_chat_id"`
}

func (m *MatchedDetail) TableName() string {
	return "matched_detail"
}
