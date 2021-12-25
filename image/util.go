package image

import (
	"fmt"
	"strings"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
)

func IsAvailable(colorScheme, image string, seasonID, week int, team string) bool {
	// check if file already exists
	imageFilename := ImageFilename(image, seasonID, week, team)
	metaFilename := MetadataFilename(image, seasonID, week, team)
	if util.FileExists(metaFilename) && util.FileExists(imageFilename) {
		metadata := GetMetadata(metaFilename)
		if metadata.ColorScheme != colorScheme && len(colorScheme) > 0 {
			log.Debugf("file [%s] has a different colorScheme, needs to be regenerated", imageFilename)
			return false // cached image has a different colorscheme, needs to be regenerated
		}

		if week <= 0 {
			metadata.Week = 12 // set to 12 if we want to calculate a seasonal image file from last season ago
		}

		// if it's from a week longer than 10 days ago and updated somewhere within 10 days after weekstart
		if time.Since(metadata.StartDate.AddDate(0, 0, metadata.Week*7)) > time.Hour*24*10 &&
			metadata.LastUpdated.Sub(metadata.StartDate.AddDate(0, 0, metadata.Week*7)) > time.Hour*24*10 {
			//log.Debugf("file [%s] already exists", imageFilename)
			return true
		}

		// if it's a team file and older than 15 minutes
		if len(metadata.Team) > 0 && time.Since(metadata.LastUpdated) > time.Minute*15 {
			//log.Debugf("file [%s] is older than 15 minutes, needs to be regenerated", imageFilename)
			return false // cached image needs to be regenerated
		}

		// if it's younger than 2 hours
		if time.Since(metadata.LastUpdated) < time.Hour*2 {
			//log.Debugf("file [%s] already exists", imageFilename)
			return true
		}
	}
	return false // cached image needs to be regenerated
}

func ImageFilename(image string, seasonID, week int, team string) string {
	if len(team) > 0 {
		team = strings.ReplaceAll(strings.ToLower(team), " ", "_")
		if week <= 0 {
			return fmt.Sprintf("public/%s/season_%d_%s.png", image, seasonID, team)
		}
		return fmt.Sprintf("public/%s/season_%d_week_%d_%s.png", image, seasonID, week, team)
	}

	if week <= 0 {
		return fmt.Sprintf("public/%s/season_%d.png", image, seasonID)
	}
	return fmt.Sprintf("public/%s/season_%d_week_%d.png", image, seasonID, week)
}

func GetResult(slot time.Time, results []database.RaceWeekResult) database.RaceWeekResult {
	sessions := make([]database.RaceWeekResult, 0)
	for _, result := range results {
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

func MapValueIntoRange(rangeStart, rangeEnd, min, max, value int) int {
	if value <= min {
		value = min + 1
	}
	if value >= max {
		return rangeEnd
	}
	rangeSize := rangeEnd - rangeStart
	return rangeStart + int((float64(value-min)/float64(max-min))*float64(rangeSize))
}
