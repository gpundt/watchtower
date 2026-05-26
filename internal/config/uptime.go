package config

import (
	"fmt"
	"time"
)

var StartTime time.Time

func init() {
	StartTime = time.Now()
}

// Calculates delta between StartTime and now
func GetUptime() time.Duration {
	return time.Since(StartTime)
}

// Formats GetUptime() into string used for health_check api Response
func GetUptimeString() string {
	// Round to nearest second
	elapsed := GetUptime()
	elapsed = elapsed.Round(time.Second)

	days := elapsed / (24 * time.Hour)
	elapsed -= days * 24 * time.Hour

	hours := elapsed / time.Hour
	elapsed -= hours * time.Hour

	minutes := elapsed / time.Minute
	elapsed -= minutes * time.Minute

	seconds := elapsed / time.Second

	return fmt.Sprintf(
		"%d days %02d:%02d:%02d",
		days,
		hours,
		minutes,
		seconds,
	)
}

// Formats time.now into string
func GetTimestampString() string {
	now := time.Now()

	return now.Format(time.RFC3339)
}