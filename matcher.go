package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"strangerbot/repository"
	"strangerbot/repository/model"
	"strangerbot/service"
	"strangerbot/vars"
)

func matchUsers(chatIDs <-chan int64) {

	ctx := context.TODO()

	for c := range chatIDs {

		user, err := retrieveUser(c)

		if err != nil {
			log.Printf("Error in matcher: %s", err)
			continue
		}

		if !user.Available || user.MatchChatID.Valid {
			log.Println("User already assigned")
			continue
		}

		matchUser, err := service.ServiceMatch(ctx, user.ChatID, user.IsVerify)
		if err != nil {
			log.Printf("Error retrieving available users: %s", err)
			continue
		}

		if matchUser == nil {
			continue
		}

		// start get user profile and matched user profile

		repo := repository.GetRepository()

		// find all question
		questions, err := repo.GetAllQuestion(ctx)
		if err != nil {
			continue
		}

		// find all question option
		options, err := repo.GetAllOption(ctx)
		if err != nil {
			continue
		}

		// find user all user question data
		userQuestionData, err := repo.GetUserQuestionDataByUser(ctx, user.ChatID)
		if err != nil {
			continue
		}

		matchedUserQuestionData, err := repo.GetUserQuestionDataByUser(ctx, matchUser.ChatID)
		if err != nil {
			continue
		}

		// user and matched user profile
		profileQuestion := model.Questions(questions).GetProfileQuestion()
		profileOptions := model.Options(options).GetQuestionOptions(ctx, profileQuestion)
		userProfile := model.UserQuestionDataList(userQuestionData).GetUserQuestionDataByOptions(ctx, profileOptions)
		matchedUserProfile := model.UserQuestionDataList(matchedUserQuestionData).GetUserQuestionDataByOptions(ctx, profileOptions)
		userProfileStr := strings.Join(profileOptions.GetOptionsByIds(userProfile.GetOptionIds()), ",")
		matchedUserProfileStr := strings.Join(profileOptions.GetOptionsByIds(matchedUserProfile.GetOptionIds()), ",")

		// user and matched user goals
		var userGoals, matchedUserGoals string
		if vars.GoalsQuestionId > 0 {
			goalsOptions := model.Options(options).GetOptionsByQuestionId(vars.GoalsQuestionId)
			userGoalsOptions := model.UserQuestionDataList(userQuestionData).GetUserQuestionDataByOptions(ctx, goalsOptions)
			userGoals = strings.Join(model.Options(options).GetOptionsByIds(userGoalsOptions.GetOptionIds()), ",")
			matchedUserGoalsOptions := model.UserQuestionDataList(matchedUserQuestionData).GetUserQuestionDataByOptions(ctx, goalsOptions)
			matchedUserGoals = strings.Join(model.Options(options).GetOptionsByIds(matchedUserGoalsOptions.GetOptionIds()), ",")
		}

		createMatch(user.ChatID, user.ID, matchUser.ChatID, matchUser.ID, userProfileStr, userGoals, matchedUserProfileStr, matchedUserGoals)

		time.Sleep(100 * time.Millisecond)
	}

}

func createMatch(userChatId, userId, matchUserChatId, matchUserId int64, userProfile, userGoals, matchedUserProfile, matchedUserGoals string) {

	query := "UPDATE users SET match_chat_id = ? WHERE id = ?"

	db.Exec(query, userChatId, matchUserId)
	db.Exec(query, matchUserChatId, userId)

	// record
	_ = service.ServiceMatchedDetailRecord(context.Background(), userChatId, matchUserChatId)
	_ = service.ServiceMatchedDetailRecord(context.Background(), matchUserChatId, userChatId)

	if len(userProfile) == 0 {
		userProfile = "(!NOT SETTING)"
	}

	if len(userGoals) == 0 {
		userGoals = "(!NOT SETTING)"
	}

	if len(matchedUserProfile) == 0 {
		matchedUserProfile = "(!NOT SETTING)"
	}

	if len(matchedUserGoals) == 0 {
		matchedUserGoals = "(!NOT SETTING)"
	}

	_, _ = telegram.SendMessage(matchUserChatId, fmt.Sprintf(vars.MatchedMessage, userProfile, userGoals), emptyOpts)
	_, _ = telegram.SendMessage(userChatId, fmt.Sprintf(vars.MatchedMessage, matchedUserProfile, matchedUserGoals), emptyOpts)

}
