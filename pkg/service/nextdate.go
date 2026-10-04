package service

import (
	"fmt"
	"strings"
	"time"
)

func nextWeekday(now time.Time, date time.Time, weeks []int) time.Time {

	//находим ближайший следующий день недели из правила относительно начальной даты
	var (
		start, k int
		shift    bool
	)

	start = weeks[0]
	for k = range weeks {
		if weeks[k] >= int(date.Weekday()) {
			start = weeks[k]
			date = date.AddDate(0, 0, start-int(date.Weekday()))
			if AfterNow(date, now) {
				return date
			}
			shift = true
			break
		}
	}
	if !shift {
		date = date.AddDate(0, 0, (7%int(date.Weekday()))+start)
		k = 0
		if AfterNow(date, now) {
			return date
		}
	}

	for ; ; k++ {
		if weeks[k] == weeks[len(weeks)-1] {
			k = 0
			date = date.AddDate(0, 0, (7-int(date.Weekday()))+weeks[k])
			if AfterNow(date, now) {
				return date
			}
			continue
		}
		date = date.AddDate(0, 0, weeks[k+1]-weeks[k])
		if AfterNow(date, now) {
			return date
		}

	}
}

func AfterNow(date time.Time, now time.Time) bool {
	if date.Truncate(time.Hour * 24).After(now.Truncate(time.Hour * 24)) {
		return true
	}
	return false
}

func forAddDate(date time.Time, now time.Time, y, m, d int) time.Time {
	for {
		date = date.AddDate(y, m, d)
		if AfterNow(date, now) {
			break
		}
	}
	return date
}
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("conversion error:%w", err)
	}

	rule := strings.Split(repeat, " ")

	var years, days int

	switch rule[0] {
	case "d":
		if days, err = parseDay(rule); err != nil {
			return "", err
		}
		date = forAddDate(date, now, 0, 0, days)
	case "y":
		if years, err = parseYear(rule); err != nil {
			return "", err
		}
		date = forAddDate(date, now, years, 0, 0)
	case "w":
		weeks, err := parseWeek(rule)
		if err != nil {
			return "", err
		}
		date = nextWeekday(now, date, weeks)
	case "m":
		parseMonth(rule)
	default:
		return "", fmt.Errorf("incorrect rule format")
	}

	return date.Format(DateFormat), nil
}
