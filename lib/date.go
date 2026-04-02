package lib

import (
	"fmt"
	"sort"
	"time"
)

func NearestWeekdays(reference time.Time, weekdayStr string) ([]string, error) {
	/*startDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}*/
	startDate := reference

	seen := make(map[int]bool)
	var weekdays []int

	for _, ch := range weekdayStr {
		day := int(ch - '0')
		if day < 0 || day > 6 {
			return nil, fmt.Errorf("invalid weekday: %c", ch)
		}
		if !seen[day] {
			seen[day] = true
			weekdays = append(weekdays, day)
		}
	}

	sort.Ints(weekdays)

	result := make([]string, 0, len(weekdays))

	for _, wantedDay := range weekdays {
		currentDay := int(startDate.Weekday()) // Go: Sunday=0 ... Saturday=6

		diff := wantedDay - currentDay
		if diff < 0 {
			diff += 7
		}

		targetDate := startDate.AddDate(0, 0, diff)
		result = append(result, targetDate.Format("2006-01-02"))
	}

	return result, nil
}

func InNextMonth(reference time.Time, month, week, day int) (string, error) {
	if month != 0 && month != 1 {
		return "", fmt.Errorf("invalid month value: %d, must be 0 or 1", month)
	}
	if week < 0 || week > 3 {
		return "", fmt.Errorf("invalid week value: %d, must be between 0 and 3", week)
	}
	if day < 0 || day > 6 {
		return "", fmt.Errorf("invalid day value: %d, must be between 0 and 6", day)
	}

	// Move to the first day of the target month
	firstOfMonth := time.Date(
		reference.Year(),
		reference.Month(),
		1, 0, 0, 0, 0,
		reference.Location(),
	).AddDate(0, month, 0)

	// Weekday of the first day of target month
	firstWeekday := int(firstOfMonth.Weekday()) // Sunday=0 ... Saturday=6

	// How many days to add to reach the first requested weekday
	offset := day - firstWeekday
	if offset < 0 {
		offset += 7
	}

	// First occurrence of requested weekday in the target month
	firstMatch := firstOfMonth.AddDate(0, 0, offset)

	// Add 7 days per requested week
	resultDate := firstMatch.AddDate(0, 0, week*7)

	// Make sure result is still inside the target month
	if resultDate.Month() != firstOfMonth.Month() {
		return "", fmt.Errorf("requested week/day combination does not exist in target month")
	}

	return resultDate.Format("2006-01-02"), nil
}
