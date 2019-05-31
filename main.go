package main

import (
	"fmt"
	"sort"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/env"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
	"github.com/robfig/cron"
)

const (
	imageHeight    = float64(480)
	imageLength    = float64(1024)
	headerHeight   = float64(45)
	timeslotHeight = float64(50)
	dayLength      = float64(160)
)

var (
	username, password string
	db                 database.Database
)

func main() {
	port := env.Get("PORT", "8080")
	level := env.Get("LOG_LEVEL", "info")
	username = env.MustGet("AUTH_USERNAME")
	password = env.MustGet("AUTH_PASSWORD")

	log.Infoln("port:", port)
	log.Infoln("log level:", level)
	log.Infoln("auth username:", username)

	// setup database connection
	db = database.NewDatabase(database.NewAdapter())

	seasonID := 2377
	week := 11

	season, err := getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		return
	}

	raceweek, track, results, err := getWeek(seasonID, week)
	if err != nil {
		log.Errorf("could not get raceweek results: %v", err)
		return
	}

	drawHeatmap(season, raceweek, track, results)
}

func getSeason(seasonID int) (database.Season, error) {
	log.Infof("collect season [%d]", seasonID)
	return db.GetSeasonByID(seasonID)
}

func getWeek(seasonID, week int) (database.RaceWeek, database.Track, []database.RaceWeekResult, error) {
	log.Infof("collect results for season [%d], week [%d]", seasonID, week)

	raceweek, err := db.GetRaceWeekBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return database.RaceWeek{}, database.Track{}, nil, err
	}

	track, err := db.GetTrackByID(raceweek.TrackID)
	if err != nil {
		return raceweek, database.Track{}, nil, err
	}

	results, err := db.GetRaceWeekResultsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return raceweek, track, nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.Before(results[j].StartTime)
	})

	return raceweek, track, results, nil
}

func drawHeatmap(season database.Season, week database.RaceWeek, track database.Track, results []database.RaceWeekResult) {
	heatmapTitle := fmt.Sprintf("%s - Week %d", season.SeasonName, week.RaceWeek+1)
	heatmap2ndTitle := track.Name
	if len(track.Config) > 0 {
		heatmap2ndTitle = fmt.Sprintf("%s - %s", track.Name, track.Config)
	}

	log.Infof("draw heatmap for [%s] - [%s]", heatmapTitle, heatmap2ndTitle)
	dc := gg.NewContext(int(imageLength), int(imageHeight))

	// background
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, imageLength, headerHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()
	dc.DrawRectangle(imageLength/2+dayLength/2, 0, imageLength/2, headerHeight)
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.Fill()

	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 20); err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(heatmapTitle, dayLength/4, headerHeight/2, 0, 0.5)

	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 20); err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(heatmap2ndTitle, imageLength-dayLength/4, headerHeight/2, 1, 0.5)

	// timeslots
	dc.DrawRectangle(0, headerHeight, dayLength, timeslotHeight)
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.Fill()
	slots := 12
	timeslotLength := ((imageLength - dayLength) / float64(slots)) - 1
	for slot := 0; slot < slots; slot++ {
		dc.DrawRectangle((float64(slot)*(timeslotLength+1))+(dayLength+1), headerHeight, timeslotLength, timeslotHeight)
		if slot%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
	}

	// weekdays
	days := 7
	dayHeight := ((imageHeight - headerHeight - timeslotHeight) / float64(days)) - 1
	for day := 0; day < days; day++ {
		dc.DrawRectangle(0, (float64(day)*(dayHeight+1))+(headerHeight+timeslotHeight+1), dayLength, dayHeight)
		if day%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
	}

	// empty events
	eventDays := 7
	eventSlots := 12

	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := p.Parse(season.Timeslots)
	if err != nil {
		log.Errorf("could not parse timeslot [%s] to crontab format: %v", season.Timeslots, err)
		return
	}

	start := database.WeekStart(season.StartDate.UTC().AddDate(0, 0, week.RaceWeek*7))
	for d := 0; d < 7; d++ {
		log.Infoln(schedule.Next(start.AddDate(0, 0, d)))
	}

	eventHeight := ((imageHeight - headerHeight - timeslotHeight) / float64(eventDays)) - 1
	eventLength := ((imageLength - dayLength) / float64(eventSlots)) - 1
	for day := 0; day < eventDays; day++ {
		for slot := 0; slot < eventSlots; slot++ {
			dc.DrawRectangle(
				(float64(slot)*(eventLength+1))+(dayLength+1),
				(float64(day)*(eventHeight+1))+(headerHeight+timeslotHeight+1),
				eventLength, eventHeight)
			dc.SetRGB255(255, 255, 255) // white
			dc.Fill()
		}
	}

	dc.SavePNG("public/test.png")
}

/*
	Colors:
	dc.SetRGB255(255, 255, 255) // white
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
