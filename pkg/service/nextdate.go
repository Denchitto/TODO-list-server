package service

import (
	"fmt"
	"strings"
	"time"
)

func AfterNow(date time.Time, now time.Time) bool {

	if date.Truncate(time.Hour * 24).After(now.Truncate(time.Hour * 24)) {
		return true
	}
	return false
}

func nextMonthday(now time.Time, date time.Time, day, month uint32, lastDays [2]bool) time.Time {
	// Если месяцев нет, выставляем всю маску месяца на единицы
	if month == 0 {
		month--
	}

	for ; ; date = date.AddDate(0, 0, 1) {
		var lastDay, preLastDay bool
		if lastDays[0] {
			lastDay = date.AddDate(0, 0, 1).Month() != date.Month()
		}
		if lastDays[1] {
			preLastDay = date.AddDate(0, 0, 2).Month() != date.Month()
		}
		// Побитово сравниваем маску и текущий день, либо если у нас текущий день последний или предпоследний
		if (day&(1<<date.Day()) != 0 || lastDay || preLastDay) && (month&(1<<date.Month())) != 0 {
			if AfterNow(date, now) {
				return date
			}
		}
	}
}

func nextWeekday(now time.Time, date time.Time, weekDaysMask uint32) time.Time {

	for ; ; date = date.AddDate(0, 0, 1) {
		if weekDaysMask&(1<<date.Weekday()) != 0 {
			if AfterNow(date, now) {
				return date
			}
		}
	}

}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("conversion error:%w", err)
	}

	rule := strings.Split(repeat, " ")

	switch rule[0] {
	case "d":
		days, err := parseDay(rule)
		if err != nil {
			return "", err
		}

		for {
			date = date.AddDate(0, 0, days)
			if AfterNow(date, now) {
				break
			}

		}
	case "y":
		years, err := parseYear(rule)
		if err != nil {
			return "", err
		}

		for {
			date = date.AddDate(years, 0, 0)
			if AfterNow(date, now) {
				break
			}

		}
	case "w":
		weekDaysMask, err := parseWeek(rule)
		if err != nil {
			return "", err
		}
		date = nextWeekday(now, date, weekDaysMask)
	case "m":
		monthDaysMask, monthsMask, lastDays, err := parseMonth(rule)
		if err != nil {
			return "", err
		}
		date = nextMonthday(now, date, monthDaysMask, monthsMask, lastDays)
	default:
		return "", fmt.Errorf("incorrect rule format")
	}

	return date.Format(DateFormat), nil
}
