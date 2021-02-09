package util

import "time"

func YearToUnixTimestamp(year int) int64 {
	date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	return date.Unix()
}

func UnixTimestampToDate(timestamp int64) time.Time {
	date := time.Unix(timestamp, 0)
	return date
}
