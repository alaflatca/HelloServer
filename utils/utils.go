package utils

import "time"

func Now() time.Time {
	return time.Now()
}

func NowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
