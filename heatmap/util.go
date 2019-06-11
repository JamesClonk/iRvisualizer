package heatmap

import (
	"time"

	"github.com/JamesClonk/iRcollector/database"
)

func (h *Heatmap) GetResult(slot time.Time) database.RaceWeekResult {
	sessions := make([]database.RaceWeekResult, 0)
	for _, result := range h.Results {
		if result.StartTime.UTC().Weekday() == slot.UTC().Weekday() &&
			result.StartTime.UTC().Hour() == slot.UTC().Hour() &&
			result.StartTime.UTC().Minute() == slot.UTC().Minute() {
			sessions = append(sessions, result)
		}
	}

	// summarize splits
	result := database.RaceWeekResult{
		SizeOfField:     0,
		StrengthOfField: 0,
	}
	for _, session := range sessions {
		result.Official = session.Official
		result.SizeOfField += session.SizeOfField
		if session.StrengthOfField > result.StrengthOfField {
			result.StrengthOfField = session.StrengthOfField
		}
	}
	return result
}

func mapValueIntoRange(rangeStart, rangeEnd, min, max, value int) int {
	if value <= min {
		value = min + 1
	}
	if value >= max {
		return rangeEnd
	}
	rangeSize := rangeEnd - rangeStart
	return rangeStart + int((float64(value-min)/float64(max-min))*float64(rangeSize))
}