package utils

import "time"

func WeekIntervalTime(week int) (startTime, endTime string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
	endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

	return startTime, endTime
}

func DayRangeZero(day int) (startTime, endTime int64) {

	now := time.Now()
	t := now.AddDate(0, 0, -(day - 1))
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	return start.Unix(), end.Unix()
}

func DayRange(day int) (startTime, endTime int64) {

	now := time.Now()
	t := now.AddDate(0, 0, -(day - 1))
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	return t.Unix(), end.Unix()
}
