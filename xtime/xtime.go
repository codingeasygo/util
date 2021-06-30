package xtime

import "time"

func Now() int64 {
	return time.Now().Local().UnixNano() / 1e6
}

// TimeStartOfToday return 00:00:00 of today
func TimeStartOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TimeStartOfWeek return 00:00:00 of week
func TimeStartOfWeek() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := today.Add(-24 * time.Hour * time.Duration(today.Weekday()))
	return start
}

// TimeStartOfMonth return 00:00:00 of month
func TimeStartOfMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
}
