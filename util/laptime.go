package util

import (
	"fmt"
	"time"

	"github.com/JamesClonk/iRcollector/database"
)

func ConvertLaptime(laptime database.Laptime) string {
	duration := time.Duration(laptime*100) * time.Microsecond

	seconds := int64(duration.Seconds()) % 60
	minutes := int64(duration.Minutes()) % 60
	hours := int64(duration.Hours()) % 24
	days := int64(duration/(24*time.Hour)) % 365 % 7
	leftYearDays := int64(duration/(24*time.Hour)) % 365
	weeks := leftYearDays / 7
	if leftYearDays >= 364 && leftYearDays < 365 {
		weeks = 52
	}
	years := int64(duration/(24*time.Hour)) / 365
	milliseconds := int64(duration/time.Millisecond) -
		(seconds * 1000) - (minutes * 60000) - (hours * 3600000) -
		(days * 86400000) - (weeks * 604800000) - (years * 31536000000)

	return fmt.Sprintf("%d:%02d.%03d", minutes, seconds, milliseconds)
}
