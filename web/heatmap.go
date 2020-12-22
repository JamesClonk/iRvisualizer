package web

import (
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image/heatmap"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

var heatmapMutex = &sync.Mutex{}

func (h *Handler) weeklyHeatmap(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seasonID, err := strconv.Atoi(vars["seasonID"])
	if err != nil {
		log.Errorf("could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
		h.failure(rw, req, err)
		return
	}
	if seasonID < 2000 || seasonID > 9999 {
		seasonID = 2377
	}
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		log.Errorf("could not convert week [%s] to int: %v", vars["week"], err)
		h.failure(rw, req, err)
		return
	}
	if week < 1 || week > 13 { // allow leap weeks
		week = 1
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a minSOF given?
	minSOF := 1000
	value := req.URL.Query().Get("minSOF")
	if len(value) > 0 {
		minSOF, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert minSOF [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a maxSOF given?
	maxSOF := 2700
	value = req.URL.Query().Get("maxSOF")
	if len(value) > 0 {
		maxSOF, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert maxSOF [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a forceOverwrite given?
	forceOverwrite := false
	value = req.URL.Query().Get("forceOverwrite")
	if len(value) > 0 {
		forceOverwrite, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("could not convert forceOverwrite [%s] to bool: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// do we need to update the image file?
	// check if file already exists and is up-to-date, serve it immediately if yes
	if !forceOverwrite && heatmap.IsAvailable(colorScheme, seasonID, week) {
		http.ServeFile(rw, req, heatmap.Filename(seasonID, week))
		return
	}
	// lock global mutex
	heatmapMutex.Lock()
	defer heatmapMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && heatmap.IsAvailable(colorScheme, seasonID, week) {
		http.ServeFile(rw, req, heatmap.Filename(seasonID, week))
		return
	}

	// create/update heatmap image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("heatmap: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		track.Name = "starting soon..."
	}
	results, err := h.getRaceWeekResults(seasonID, week-1)
	if err != nil {
		log.Errorf("heatmap: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		h.failure(rw, req, err)
		return
	}
	hm := heatmap.New(colorScheme, season, raceweek, track, results)
	if err := hm.Draw(minSOF, maxSOF, true); err != nil {
		log.Errorf("could not create heatmap season[%d], week[%d]: %v", seasonID, week-1, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, heatmap.Filename(seasonID, week))
}

func (h *Handler) seasonalHeatmap(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seasonID, err := strconv.Atoi(vars["seasonID"])
	if err != nil {
		log.Errorf("could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
		h.failure(rw, req, err)
		return
	}
	if seasonID < 2000 || seasonID > 9999 {
		seasonID = 2377
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a minSOF given?
	minSOF := 900
	value := req.URL.Query().Get("minSOF")
	if len(value) > 0 {
		minSOF, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert minSOF [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a maxSOF given?
	maxSOF := 2700
	value = req.URL.Query().Get("maxSOF")
	if len(value) > 0 {
		maxSOF, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert maxSOF [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a forceOverwrite given?
	forceOverwrite := false
	value = req.URL.Query().Get("forceOverwrite")
	if len(value) > 0 {
		forceOverwrite, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("could not convert forceOverwrite [%s] to bool: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// do we need to update the image file?
	// check if file already exists and is up-to-date, serve it immediately if yes
	if !forceOverwrite && heatmap.IsAvailable(colorScheme, seasonID, -1) {
		http.ServeFile(rw, req, heatmap.Filename(seasonID, -1))
		return
	}
	// lock global mutex
	heatmapMutex.Lock()
	defer heatmapMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && heatmap.IsAvailable(colorScheme, seasonID, -1) {
		http.ServeFile(rw, req, heatmap.Filename(seasonID, -1))
		return
	}

	// create/update heatmap image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}

	// figure out timeslots schedule
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := p.Parse(season.Timeslots)
	if err != nil {
		log.Errorf("could not parse timeslot [%s] to crontab format: %v", season.Timeslots, err)
		h.failure(rw, req, err)
		return
	}

	// sum/avg all weeks together
	var weeksFound int
	allResults := make([][]database.RaceWeekResult, 0)
	for week := 0; week < 12; week++ {
		rs, err := h.getRaceWeekResults(seasonID, week)
		if err != nil {
			log.Debugf("seasonal heatmap: could not get raceweek results for season[%d], week[%d]: %v", seasonID, week, err)
		}
		if len(rs) > 0 {
			weeksFound++
		}

		// sum splits
		sessions := make(map[int]database.RaceWeekResult, 0)
		for _, r := range rs {
			if session, found := sessions[r.SessionID]; found {
				session.SizeOfField += r.SizeOfField
				if r.StrengthOfField > session.StrengthOfField {
					session.StrengthOfField = r.StrengthOfField
				}
				sessions[r.SessionID] = session
			} else {
				sessions[r.SessionID] = r
			}
		}

		// add to allResults collection
		weekResults := make([]database.RaceWeekResult, 0)
		for _, session := range sessions {
			weekResults = append(weekResults, session)
		}
		allResults = append(allResults, weekResults)
	}

	// go through all timeslots and calculate final result
	finalResults := make([]database.RaceWeekResult, 0)
	start := database.WeekStart(season.StartDate.UTC().AddDate(0, 0, 7)).Add(-1 * time.Minute)
	timeslots := make([]time.Time, 0)
	next := schedule.Next(start)                             // get first timeslot
	for next.Before(schedule.Next(start.AddDate(0, 0, 7))) { // collect all timeslots of 1 week
		timeslots = append(timeslots, next)
		next = schedule.Next(next)
	}
	sessionIdx := 0
	for _, timeslot := range timeslots {
		sessionIdx++
		finalResult := database.RaceWeekResult{
			SessionID:       sessionIdx,
			StartTime:       timeslot,
			Official:        false,
			SizeOfField:     0,
			StrengthOfField: 0,
		}

		// see how many results exist for that timeslot
		raceWeeksFound := make(map[int]bool, 0)
		for _, results := range allResults {
			for _, result := range results {
				if timeslot.UTC().Weekday() == result.StartTime.UTC().Weekday() &&
					timeslot.UTC().Hour() == result.StartTime.UTC().Hour() &&
					timeslot.UTC().Minute() == result.StartTime.UTC().Minute() {
					if result.Official {
						finalResult.Official = true
					}

					finalResult.SizeOfField += result.SizeOfField
					finalResult.StrengthOfField += result.StrengthOfField
					raceWeeksFound[result.RaceWeekID] = true
				}
			}
		}

		// average size and sof
		if len(raceWeeksFound) > 0 {
			finalResult.SizeOfField = finalResult.SizeOfField / len(raceWeeksFound)
			finalResult.StrengthOfField = finalResult.StrengthOfField / len(raceWeeksFound)
		}

		finalResults = append(finalResults, finalResult)
	}

	// sort again to be on the safe side..
	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i].StartTime.Before(finalResults[j].StartTime)
	})

	hm := heatmap.New(colorScheme, season, database.RaceWeek{RaceWeek: -1, LastUpdate: time.Now()}, database.Track{}, finalResults)
	if err := hm.Draw(minSOF, maxSOF, false); err != nil {
		log.Errorf("could not create seasonal heatmap: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, heatmap.Filename(seasonID, -1))
}
