package dates

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("нет правила повтороения")
	}

	parseDate, err := time.Parse(TimeFormat, date)
	if err != nil {
		return "", err
	}

	switch {
	case repeat == "y":
		return addYear(now, parseDate)

	case strings.HasPrefix(repeat, "d"):
		if len(repeat) < 2 {
			return "", fmt.Errorf("неправильный формат")
		}
		repSplit := strings.Split(repeat, " ")
		parseAddDays, err := strconv.Atoi(repSplit[1])
		if err != nil {
			return "", err
		}
		return addDays(now, parseDate, parseAddDays)

	case strings.HasPrefix(repeat, "w"):
		if len(repeat) < 2 {
			return "", fmt.Errorf("неправильный формат")
		}
		repSplit := strings.Split(repeat, " ")
		parseWeekDays := strings.Split(repSplit[1], ",")
		return addWeekDay(now, parseDate, parseWeekDays)

	case strings.HasPrefix(repeat, "m"):
		return addMonthDay(now, parseDate, strings.TrimPrefix(repeat, "m "))
	}
	return "", fmt.Errorf("неподдерживаемый формат")
}

func addYear(currDate time.Time, date time.Time) (string, error) {
	date = date.AddDate(1, 0, 0)
	for date.Before(currDate) || date.Equal(currDate) {
		date = date.AddDate(1, 0, 0)
	}
	return date.Format(TimeFormat), nil
}

func addDays(currDate time.Time, date time.Time, days int) (string, error) {
	if 0 < days && days < 401 {
		date = date.AddDate(0, 0, days)
		for date.Before(currDate) || date.Equal(currDate) {
			date = date.AddDate(0, 0, days)
		}
		return date.Format(TimeFormat), nil
	} else {
		return "", fmt.Errorf("значение дней не входит в допустимый интервал")
	}
}

func addWeekDay(currDate time.Time, date time.Time, daysOfWeek []string) (string, error) {
	comparDate := latestDate(currDate, date)
	nextDate := comparDate
	for _, day := range daysOfWeek {
		parseDay, err := strconv.Atoi(day)
		if err != nil || parseDay < 1 || parseDay > 7 {
			return "", fmt.Errorf("недопустимый формат дня недели")
		}
		targetDay := time.Weekday(parseDay % 7)
		for nextDate.Weekday() != targetDay {
			nextDate = nextDate.AddDate(0, 0, 1)
		}
		if nextDate.After(comparDate) {
			return nextDate.Format(TimeFormat), nil
		}
	}
	return "", fmt.Errorf("нет подходящей даты")
}

func addMonthDay(currDate time.Time, date time.Time, monthRule string) (string, error) {
	rules := strings.Split(monthRule, " ")
	if len(rules) < 1 || len(rules) > 2 {
		return "", fmt.Errorf("неверный формат правила")
	}

	comparDate := latestDate(currDate, date)
	nextDate := comparDate

	days, err := parseDays(rules[0])
	if err != nil {
		return "", err
	}

	months := make([]time.Month, 12)
	for i := range months {
		months[i] = time.Month(i)
	}

	if len(rules) > 1 {
		months, err = parseMonths(rules[1])
		if err != nil {
			return "", err
		}
	}
	for i := range 2 {
		for _, month := range months {
			for _, day := range days {
				lastDay := time.Date(comparDate.Year()+i, month+1, 1, 0, 0, 0, 0, comparDate.Location()).AddDate(0, 0, -1)
				switch {
				case 0 < day && day <= lastDay.Day():
					nextDate = time.Date(comparDate.Year()+i, month, day, 0, 0, 0, 0, comparDate.Location())
				case day == -1:
					nextDate = lastDay
				case day == -2:
					nextDate = lastDay.AddDate(0, 0, -1)
				}
				if nextDate.After(comparDate) {
					return nextDate.Format(TimeFormat), nil
				}
			}
		}

	}
	return nextDate.Format(TimeFormat), nil
}

func parseDays(s string) ([]int, error) {
	daysRules := strings.Split(s, ",")
	var days []int

	for _, dayRule := range daysRules {
		parseDay, err := strconv.Atoi(dayRule)
		if err != nil {
			return nil, fmt.Errorf("неверный формат дня: %s", dayRule)
		}

		if (parseDay < -2 || parseDay > 31) || parseDay == 0 {
			return nil, fmt.Errorf("недопустимое значение дня: %d", parseDay)
		}
		days = append(days, parseDay)
	}
	sort.Slice(days, func(i, j int) bool {
		a, b := days[i], days[j]
		if a > 0 && b > 0 {
			return a < b
		}
		if a < 0 && b < 0 {
			return a < b
		}
		return a > 0
	})
	return days, nil
}

func parseMonths(s string) ([]time.Month, error) {
	monthsRules := strings.Split(s, ",")
	var months []time.Month

	for _, monthRule := range monthsRules {
		month, err := strconv.Atoi(monthRule)
		if err != nil {
			return nil, fmt.Errorf("неверный формат месяца: %s", monthRule)
		}

		if month < 1 || month > 12 {
			return nil, fmt.Errorf("недопустимое значение месяца: %d", month)
		}
		months = append(months, time.Month(month))
	}
	slices.Sort(months)
	return months, nil
}

func latestDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
