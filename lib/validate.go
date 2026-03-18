package lib

import "time"

func IsValidDate(dateStr string) bool {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return false
	}

	// Ensure the formatted output matches exactly (prevents cases like 2026-02-30 → Mar 2)
	return t.Format(layout) == dateStr
}
