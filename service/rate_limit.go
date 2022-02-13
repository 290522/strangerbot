package service

import (
	"context"
	"errors"

	"strangerbot/repository"
	"strangerbot/repository/model"
	"strangerbot/utils"
	"strangerbot/vars"
)

var (
	ErrRateLimit = errors.New("Maximum number of matches reached")
)

func RateLimit(ctx context.Context, user *model.User, list model.UserQuestionDataList) (bool, error) {

	if vars.FemaleMatchRateLimit.RateLimitEnabled == false && vars.MaleMatchRateLimit.RateLimitEnabled == false {
		return true, nil
	}

	// default female
	rateLimitRule := vars.FemaleMatchRateLimit

	// check is male?
	isMale := list.CheckExistsOption(vars.MaleMatchRateLimit.OptionId)
	if isMale {
		rateLimitRule = vars.MaleMatchRateLimit
	}

	// check user custom?
	if rateLimitRule.RateLimitEnabled {
		if user.CustomRateLimitEnabled {
			rateLimitRule.RateLimitUnit = user.RateLimitUnit
			rateLimitRule.RateLimitUnitPeriod = user.RateLimitUnitPeriod
			rateLimitRule.MatchPerRate = user.MatchPerRate
		}
	}

	if rateLimitRule.RateLimitEnabled == false {
		return true, nil
	}

	repo := repository.GetRepository()

	startTime, endTime := utils.DayRange(int(rateLimitRule.RateLimitUnitPeriod))

	cnt, err := repo.MatchedCount(ctx, user.ChatID, startTime, endTime)
	if err != nil {
		return false, err
	}

	if cnt >= rateLimitRule.MatchPerRate {
		return false, ErrRateLimit
	}

	return true, nil
}
