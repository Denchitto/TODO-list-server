package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	DateFormat = "20060102"
)

func convertSort(values []string, min, max int, entity string) ([]int, error) {

	valuesConverted := make([]int, len(values))

	for k, v := range values {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		if num < min || num > max {
			return nil, fmt.Errorf("%s must be between %d and %d", entity, min, max)
		}
		switch entity {
		case "monthday":
			if num == 0 {
				return nil, fmt.Errorf("%s cannot be zero", entity)
			}
		case "weekday":
			if num == 7 {
				num = 0
			}
		}

		valuesConverted[k] = num
	}

	sort.Ints(valuesConverted)

	for i := 1; i < len(valuesConverted); i++ {
		if valuesConverted[i] == valuesConverted[i-1] {
			return nil, fmt.Errorf("redundant value %d", valuesConverted[i])
		}
	}
	return valuesConverted, nil
}

func parseMonth(data []string) ([]int, []int, error) {

	if len(data) < 2 || len(data) > 3 {
		return nil, nil, fmt.Errorf("month may have from ont to two arguments")
	}

	monthDays, err := convertSort(strings.Split(data[1], ","), -2, 31, "monthday")
	if err != nil {
		return nil, nil, err
	}

	if len(data) == 3 {
		months, err := convertSort(strings.Split(data[2], ","), 1, 12, "month")
		if err != nil {
			return nil, nil, err
		}
		return monthDays, months, nil
	}
	return monthDays, nil, nil
}

func parseWeek(data []string) ([]int, error) {

	if len(data) != 2 {
		return nil, fmt.Errorf("days must have one argument")
	}

	weekDays, err := convertSort(strings.Split(data[1], ","), 1, 7, "weekday")
	if err != nil {
		return nil, err
	}

	return weekDays, nil
}

func parseYear(data []string) (int, error) {

	if len(data) != 1 {
		return 0, fmt.Errorf("year requires only a flag")
	}
	return 1, nil
}

func parseDay(data []string) (int, error) {

	if len(data) != 2 {
		return 0, fmt.Errorf("days must have one argument")
	}

	days, err := strconv.Atoi(data[1])
	if err != nil {
		return 0, fmt.Errorf("conversion error: %w", err)
	}

	if days < 1 || days > 400 {
		return 0, fmt.Errorf("days is out of range (1-400), %d", days)
	}

	return days, nil
}
