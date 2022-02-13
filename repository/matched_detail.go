package repository

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"
	"strangerbot/repository/model"
)

func (r *Repository) MatchedDetailAdd(ctx context.Context, po *model.MatchedDetail) error {

	if err := r.db.Create(&po).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil
		}
		return err
	}

	return nil
}

func (r *Repository) MatchedCount(ctx context.Context, chatId int64, startTime int64, endTime int64) (int64, error) {
	var cnt struct{ Cnt int64 }

	err := r.db.Select("COUNT(*) AS cnt").Model(model.MatchedDetail{}).Where("chat_id = ? AND create_time >= ? AND create_time <= ? AND is_del = 0", chatId, startTime, endTime).Scan(&cnt).Error

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return cnt.Cnt, nil
}
