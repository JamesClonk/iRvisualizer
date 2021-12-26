package web

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image/summary"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/gorilla/mux"
)

var summaryMutex = &sync.Mutex{}

func (h *Handler) weeklySummary(rw http.ResponseWriter, req *http.Request) {
	image := "summary"

	vars := mux.Vars(req)
	seasonID, err := strconv.Atoi(vars["seasonID"])
	if err != nil {
		log.Errorf("summary: could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
		h.failure(rw, req, err)
		return
	}
	if seasonID < 2000 || seasonID > 9999 {
		seasonID = 2377
	}
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		log.Errorf("summary: could not convert week [%s] to int: %v", vars["week"], err)
		h.failure(rw, req, err)
		return
	}
	if week < 1 || week > 13 {
		week = 1
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a topN given?
	topN := 30
	value := req.URL.Query().Get("topN")
	if len(value) > 0 {
		topN, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("summary: could not convert topN [%s] to int: %v", value, err)
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
			log.Errorf("summary: could not convert forceOverwrite [%s] to bool: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// are there any individually marked drivers given?
	drivers := strings.Split(req.URL.Query().Get("drivers"), ",")

	// is there a team given?
	team := req.URL.Query().Get("team")

	// do we need to update the image file?
	// check if file already exists and is up-to-date, serve it immediately if yes
	if !forceOverwrite && summary.IsAvailable(colorScheme, seasonID, week, team) {
		http.ServeFile(rw, req, summary.Filename(seasonID, week, team))
		return
	}
	// lock global mutex
	summaryMutex.Lock()
	defer summaryMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && summary.IsAvailable(colorScheme, seasonID, week, team) {
		http.ServeFile(rw, req, summary.Filename(seasonID, week, team))
		return
	}

	// create/update summary image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("summary: could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("summary: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		raceweek.LastUpdate = time.Now()
		track.Name = "starting soon..."
	}
	var summaries []database.Summary
	if len(team) > 0 {
		summaries, err = h.getRaceWeekSummariesByTeam(seasonID, week-1, team)
		if err != nil {
			log.Errorf("summary: could not get raceweek summaries for season[%d], week[%d], team[%s]: %v", seasonID, week-1, team, err)
			h.failure(rw, req, err)
			return
		}
	} else {
		summaries, err = h.getRaceWeekSummaries(seasonID, week-1)
		if err != nil {
			log.Errorf("summary: could not get raceweek summaries for season[%d], week[%d]: %v", seasonID, week-1, err)
			h.failure(rw, req, err)
			return
		}
	}

	data := make([]summary.DataSet, 0)
	// sort by champ points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].HighestChampPoints > summaries[j].HighestChampPoints
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		data = append(data, summary.DataSet{
			Summary: summaries[i],
			Marked:  isDriverMarked(drivers, summaries[i].Driver.DriverID),
		})
	}

	hm := summary.New(colorScheme, team, season, raceweek, track, data)
	if err := hm.Draw(); err != nil {
		log.Errorf("summary: could not create weekly summary [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, summary.Filename(seasonID, week, team))
}

func (h *Handler) seasonSummary(rw http.ResponseWriter, req *http.Request) {
	image := "summary"

	vars := mux.Vars(req)
	seasonID, err := strconv.Atoi(vars["seasonID"])
	if err != nil {
		log.Errorf("summary: could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
		h.failure(rw, req, err)
		return
	}
	if seasonID < 2000 || seasonID > 9999 {
		seasonID = 2377
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a topN given?
	topN := 30
	value := req.URL.Query().Get("topN")
	if len(value) > 0 {
		topN, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("summary: could not convert topN [%s] to int: %v", value, err)
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
			log.Errorf("summary: could not convert forceOverwrite [%s] to bool: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// are there any individually marked drivers given?
	drivers := strings.Split(req.URL.Query().Get("drivers"), ",")

	// is there a team given?
	team := req.URL.Query().Get("team")
	if len(team) == 0 {
		team = "TNT Racing"
	}

	// do we need to update the image file?
	// check if file already exists and is up-to-date, serve it immediately if yes
	if !forceOverwrite && summary.IsAvailable(colorScheme, seasonID, -1, team) {
		http.ServeFile(rw, req, summary.Filename(seasonID, -1, team))
		return
	}
	// lock global mutex
	summaryMutex.Lock()
	defer summaryMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && summary.IsAvailable(colorScheme, seasonID, -1, team) {
		http.ServeFile(rw, req, summary.Filename(seasonID, -1, team))
		return
	}

	// create/update summary image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("summary: could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	var summaries []database.Summary
	summaries, err = h.getSeasonSummariesByTeam(seasonID, team)
	if err != nil {
		log.Errorf("summary: could not get season summaries for season[%d], team[%s]: %v", seasonID, team, err)
		h.failure(rw, req, err)
		return
	}

	data := make([]summary.DataSet, 0)
	// sort by champ points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].AverageChampPoints > summaries[j].AverageChampPoints
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		data = append(data, summary.DataSet{
			Summary: summaries[i],
			Marked:  isDriverMarked(drivers, summaries[i].Driver.DriverID),
		})
	}

	hm := summary.New(colorScheme, team, season, database.RaceWeek{RaceWeek: -1, LastUpdate: time.Now()}, database.Track{}, data)
	if err := hm.Draw(); err != nil {
		log.Errorf("summary: could not create season summary [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, summary.Filename(seasonID, -1, team))
}
