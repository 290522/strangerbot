package service

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strangerbot/keyboard"
	"strangerbot/repository"
	"strangerbot/repository/model"
	"strangerbot/vars"
)

func ServiceQuestionOption(ctx context.Context, chatId int64, data *keyboard.KeyboardCallbackDataPlus) ([]*tgbotapi.MessageConfig, []tgbotapi.CallbackConfig, error) {

	repo := repository.GetRepository()

	// get option
	option, err := repo.GetOptionById(ctx, data.ButtonRelId)
	if err != nil {
		return nil, nil, nil
	}

	// change gender
	if option.ID == vars.FemaleMatchRateLimit.OptionId || option.ID == vars.MaleMatchRateLimit.OptionId {
		if !vars.ChangeGenderEnabled {
			userGender, err := repo.GetUserQuestionDataByUserQuestion(ctx, chatId, option.QuestionId)
			if err != nil {
				return nil, nil, err
			}
			if userGender != nil && len(userGender) > 0 {
				if userGender[0].OptionId != option.ID {
					return nil, []tgbotapi.CallbackConfig{tgbotapi.NewCallback("", vars.ChangeGenderErrorMessage)}, nil
				}
			}
		}
	}

	// get question
	question, err := repo.GetQuestionById(ctx, option.QuestionId)
	if err != nil {
		return nil, nil, nil
	}

	switch question.FrontendType {
	case model.FRONTEND_TYPE_SELECT:

		// delete question old option
		err := repo.DeleteUserQuestionDataByQuestion(ctx, chatId, question.ID)
		if err != nil {
			return nil, nil, nil
		}

	case model.FRONTEND_TYPE_MULTI_SELECT:

		userQuestionData, err := repo.GetUserQuestionDataByOptionAndChat(ctx, option.ID, chatId)
		if err != nil {
			return nil, nil, err
		}

		// delete
		if userQuestionData != nil {
			if err := repo.DeleteUserQuestionData(ctx, userQuestionData); err != nil {
				return nil, nil, nil
			}

			// back to menu
			menu, err := repo.GetMenuByQuestionId(ctx, option.QuestionId)
			if err != nil {
				return nil, nil, nil
			}

			msgs, err := serviceBackMenu(ctx, chatId, repo, menu)
			return msgs, nil, err
		}

		// add len
		if question.MaxMultiLen > 0 {
			uq, err := repo.GetUserQuestionDataByUserQuestion(ctx, chatId, question.ID)
			if err != nil {
				return nil, nil, err
			}

			if len(uq) >= question.MaxMultiLen {
				return nil, []tgbotapi.CallbackConfig{tgbotapi.NewCallback("", fmt.Sprintf(vars.QuestionCallbackTest, question.MaxMultiLen))}, nil
			}
		}
	}

	// add
	po := &model.UserQuestionData{
		ChatId:     chatId,
		QuestionId: option.QuestionId,
		OptionId:   option.ID,
		Value:      option.Value,
	}

	if err := repo.UserQuestionDataAdd(ctx, po); err != nil {
		return nil, nil, err
	}

	// back to menu
	menu, err := repo.GetMenuByQuestionId(ctx, option.QuestionId)
	if err != nil {
		return nil, nil, nil
	}

	msgs, err := serviceBackMenu(ctx, chatId, repo, menu)
	return msgs, nil, err

}
