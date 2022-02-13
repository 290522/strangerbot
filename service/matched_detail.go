package service

import (
	"context"
	"time"

	"strangerbot/repository"
	"strangerbot/repository/gorm_global"
	"strangerbot/repository/model"
)

func ServiceMatchedDetailRecord(ctx context.Context, chatId int64, matchedChatId int64) error {

	repo := repository.GetRepository()

	// add
	po := &model.MatchedDetail{
		ColumnCreateModifyDeleteTime: gorm_global.ColumnCreateModifyDeleteTime{
			CreateTime: time.Now().Unix(),
			ModifyTime: time.Now().Unix(),
		},
		ChatId:      chatId,
		MatchChatId: matchedChatId,
	}

	if err := repo.MatchedDetailAdd(ctx, po); err != nil {
		return err
	}

	return nil
}
