package utils

import (
	"time"
)

func GetStartAndEndTime(interval string) (time.Time, time.Time) {
	now := time.Now()

	switch interval {
	case "5m":
		return now.Add(time.Duration(-5) * time.Minute), now
	case "15m":
		return now.Add(time.Duration(-15) * time.Minute), now
	case "30m":
		return now.Add(time.Duration(-30) * time.Minute), now
	case "1h":
		return now.Add(time.Duration(-1) * time.Hour), now
	case "3h":
		return now.Add(time.Duration(-3) * time.Hour), now
	case "6h":
		return now.Add(time.Duration(-6) * time.Hour), now
	case "12h":
		return now.Add(time.Duration(-12) * time.Hour), now
	case "24h":
		return now.Add(time.Duration(-24) * time.Hour), now
	case "2d":
		return now.Add(time.Duration(-48) * time.Hour), now
	case "7d":
		return now.Add(time.Duration(-168) * time.Hour), now
	case "30d":
		return now.Add(time.Duration(-720) * time.Hour), now
	default:
		return now.Add(time.Duration(-1) * time.Hour), now
	}
}
