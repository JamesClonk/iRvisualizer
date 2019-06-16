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
	if week < 1 || week > 12 {
		week = 1
	}

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
	maxSOF := 2500
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
	if !forceOverwrite && heatmap.IsAvailable(seasonID, week) {
		http.ServeFile(rw, req, heatmap.Filename(seasonID, week))
		return
	}
	// lock global mutex
	heatmapMutex.Lock()
	defer heatmapMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && heatmap.IsAvailable(seasonID, week) {
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
		log.Errorf("could not get raceweek: %v", err)
		h.failure(rw, req, err)
		return
	}
	results, err := h.getRaceWeekResults(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek results: %v", err)
		h.failure(rw, req, err)
		return
	}
	hm := heatmap.New(season, raceweek, track, results)
	if err := hm.Draw(minSOF, maxSOF, true); err != nil {
		log.Errorf("could not create heatmap: %v", err)
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
	maxSOF := 2500
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
	if !forceOverwrite && heatmap.IsAvailable(seasonID, -1) {
		http.ServeFile(rw, req, heatmap.Filename(seasonID, -1))
		return
	}
	// lock global mutex
	heatmapMutex.Lock()
	defer heatmapMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && heatmap.IsAvailable(seasonID, -1) {
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
	sessionIdx := 0
	results := make([]database.RaceWeekResult, 0)
	var weeksNotFound int
	for week := 0; week < 12; week++ {
		rs, err := h.getRaceWeekResults(seasonID, week)
		if err != nil {
			weeksNotFound++
			log.Warnf("could not get raceweek results: %v", err)
			// h.failure(rw, req, err)
			// return
		}

		// sum/avg splits first
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

		// collect all needed timeslots per week
		// and make sure all timeslots are present, even if they were a 0-show
		start := database.WeekStart(season.StartDate.UTC().AddDate(0, 0, (week+1)*7)).Add(-1 * time.Minute)
		timeslots := make([]time.Time, 0)
		next := schedule.Next(start)                             // get first timeslot
		for next.Before(schedule.Next(start.AddDate(0, 0, 7))) { // collect all timeslots of 1 week
			timeslots = append(timeslots, next)
			next = schedule.Next(next)
		}
		for _, timeslot := range timeslots {
			var found bool
			for _, session := range sessions {
				if timeslot.UTC().Weekday() == session.StartTime.UTC().Weekday() &&
					timeslot.UTC().Hour() == session.StartTime.UTC().Hour() &&
					timeslot.UTC().Minute() == session.StartTime.UTC().Minute() {
					found = true
					break
				}
			}
			if !found {
				log.Debugf("timeslot [%s] was missing for week [%d], adding it ...", timeslot, week+1)
				sessionIdx++
				sessions[sessionIdx] = database.RaceWeekResult{
					SessionID:       sessionIdx,
					StartTime:       timeslot,
					Official:        false,
					SizeOfField:     0,
					StrengthOfField: 0,
				}
			}
		}

		// merge weekly results into main results list
		for _, session := range sessions {
			var found bool
			for i, result := range results {
				if session.StartTime.UTC().Weekday() == result.StartTime.UTC().Weekday() &&
					session.StartTime.UTC().Hour() == result.StartTime.UTC().Hour() &&
					session.StartTime.UTC().Minute() == result.StartTime.UTC().Minute() {
					found = true
					if session.Official {
						results[i].Official = true
					}
					results[i].SizeOfField = result.SizeOfField + session.SizeOfField
					results[i].StrengthOfField = result.StrengthOfField + session.StrengthOfField
					break
				}
			}
			if !found {
				results = append(results, session)
			}
		}
	}
	// divide by 12 weeks (minus weeksNotFound)
	for i := range results {
		results[i].StrengthOfField = results[i].StrengthOfField / (12 - weeksNotFound)
		results[i].SizeOfField = results[i].SizeOfField / (12 - weeksNotFound)
	}
	// sort again to be on the safe side..
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.Before(results[j].StartTime)
	})

	hm := heatmap.New(season, database.RaceWeek{RaceWeek: -1}, database.Track{}, results)
	if err := hm.Draw(minSOF, maxSOF, false); err != nil {
		log.Errorf("could not create seasonal heatmap: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, heatmap.Filename(seasonID, -1))
}
