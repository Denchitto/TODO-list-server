package service

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	DateFormat       = "20060102"
	DateFormatSearch = "02.01.2006"
)

func convertToMask(values []string, min, max int, entity string) (uint32, error) {
	// Так как у нас дни недели и месяца - уникальные значения, мы создаем маску
	// и присваиваем каждой позиции соответствующий день. Например 31 20 1
	// 10000000000000000000000000000000 - 31
	// 00000000000100000000000000000000 - 20
	// 00000000000000000000000000000010 - 1
	// 10000000000100000000000000000010 - 31 | 20 | 1
	var mask uint32
	for _, v := range values {

		num, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}

		if num < min || num > max {
			return 0, fmt.Errorf("%s must be between %d and %d", entity, min, max)
		}

		switch entity {
		case "monthday":
			switch num {
			case 0:
				return 0, fmt.Errorf("%s cannot be zero", entity)
			case -1, -2:
				continue
			}
		case "weekday":
			if num == 7 {
				num = 0
			}
		}
		// Проверяем что на позиции num в маске пусто, если значение уже записано возвращаем ошибку об избыточности
		if mask&(1<<num) != 0 {
			return 0, fmt.Errorf("redundant %s: %d", entity, num)
		}
		mask |= 1 << num
	}

	return mask, nil
}

func parseMonth(rule []string) (uint32, uint32, [2]bool, error) {

	if len(rule) < 2 || len(rule) > 3 {
		return 0, 0, [2]bool{}, fmt.Errorf("month rule may have from one up to two arguments")
	}

	monthdays := strings.Split(rule[1], ",")
	lastDays := [2]bool{}
	for _, v := range monthdays {
		if v == "-1" {
			lastDays[0] = true
		}
		if v == "-2" {
			lastDays[1] = true
		}
	}

	monthDaysMask, err := convertToMask(monthdays, -2, 31, "monthday")
	if err != nil {
		return 0, 0, [2]bool{}, err
	}

	if len(rule) == 3 {
		monthsMask, err := convertToMask(strings.Split(rule[2], ","), 1, 12, "month")
		if err != nil {
			return 0, 0, [2]bool{}, err
		}
		return monthDaysMask, monthsMask, lastDays, nil
	}
	return monthDaysMask, 0, lastDays, nil
}

func parseWeek(rule []string) (uint32, error) {

	if len(rule) != 2 {
		return 0, fmt.Errorf("weeks rule requires one argument")
	}

	weekDaysMask, err := convertToMask(strings.Split(rule[1], ","), 1, 7, "weekday")
	if err != nil {
		return 0, err
	}

	return weekDaysMask, nil
}

func parseYear(rule []string) (int, error) {

	if len(rule) != 1 {
		return 0, fmt.Errorf("redundant value; year rule requires only a flag")
	}
	return 1, nil
}

func parseDay(rule []string) (int, error) {

	if len(rule) != 2 {
		return 0, fmt.Errorf("days rule requires one argument")
	}

	days, err := strconv.Atoi(rule[1])
	if err != nil {
		return 0, fmt.Errorf("conversion error: %w", err)
	}

	if days < 1 || days > 400 {
		return 0, fmt.Errorf("day is out of range (1-400), %d", days)
	}
	return days, nil
}
