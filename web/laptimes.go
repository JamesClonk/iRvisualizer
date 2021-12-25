package web

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image/laptime"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
	"github.com/gorilla/mux"
)

var laptimeMutex = &sync.Mutex{}

func (h *Handler) weeklyLaptimes(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seasonID, err := strconv.Atoi(vars["seasonID"])
	if err != nil {
		log.Errorf("laptimes: could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
		h.failure(rw, req, err)
		return
	}
	if seasonID < 2000 || seasonID > 9999 {
		seasonID = 2377
	}
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		log.Errorf("laptimes: could not convert week [%s] to int: %v", vars["week"], err)
		h.failure(rw, req, err)
		return
	}
	if week < 1 || week > 13 {
		week = 1
	}

	// was there a reference lap given?
	var refLap int
	lap := req.URL.Query().Get("laptime")
	if strings.Contains(lap, "s") { // 1m23s456ms format
		refLap = int(util.ParseLaptime(lap))
	} else { // int milliseconds
		refLap, _ = strconv.Atoi(lap)
		refLap = refLap * 10
	}
	if refLap < 1 {
		refLap = 0
	}
	refName := req.URL.Query().Get("reference")
	if len(refName) == 0 {
		refName = "Reference"
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a forceOverwrite given?
	forceOverwrite := false
	value := req.URL.Query().Get("forceOverwrite")
	if len(value) > 0 {
		forceOverwrite, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("laptimes: could not convert forceOverwrite [%s] to bool: %v", value, err)
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
	if !forceOverwrite && laptime.IsAvailable(colorScheme, seasonID, week, team) {
		http.ServeFile(rw, req, laptime.Filename(seasonID, week, team))
		return
	}
	// lock global mutex
	laptimeMutex.Lock()
	defer laptimeMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && laptime.IsAvailable(colorScheme, seasonID, week, team) {
		http.ServeFile(rw, req, laptime.Filename(seasonID, week, team))
		return
	}

	// create/update ranking image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("laptimes: could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("laptimes: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		raceweek.LastUpdate = time.Now()
		track.Name = "starting soon..."
	}
	timeRankings, err := h.getRaceWeekTimeRankings(seasonID, week-1)
	if err != nil {
		log.Errorf("laptimes: could not get raceweek timerankings: %v", err)
		h.failure(rw, req, err)
		return
	}

	// sort by fastest race laptimes
	sort.Slice(timeRankings, func(i, j int) bool {
		return timeRankings[i].Race < timeRankings[j].Race
	})

	// collect first/fastest driver for each division, 1-5
	laptimes := make([]laptime.DataSet, 0)
	if refLap > 0 && len(refName) > 0 {
		laptimes = append(laptimes, laptime.DataSet{
			Division: "-",
			Driver:   refName,
			Laptime:  database.Laptime(refLap),
		})
	}
	for division := 1; division <= 5; division++ {
		for _, timeRanking := range timeRankings {
			if timeRanking.Driver.Division == division && timeRanking.Race > 100 {
				laptimes = append(laptimes, laptime.DataSet{
					Division: fmt.Sprintf("%v", timeRanking.Driver.Division),
					Driver:   timeRanking.Driver.Name,
					Laptime:  timeRanking.Race,
					Marked:   isDriverMarked(drivers, timeRanking.Driver.DriverID) || (timeRanking.Driver.Team == team && len(team) > 0),
				})
				break
			}
		}
	}

	l := laptime.New(colorScheme, team, season, raceweek, track, laptimes)
	if err := l.Draw(); err != nil {
		log.Errorf("laptimes: could not create weekly laptime chart: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, laptime.Filename(seasonID, week, team))
}
