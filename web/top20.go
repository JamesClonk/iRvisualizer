package web

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/JamesClonk/iRvisualizer/image/top20"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/gorilla/mux"
)

var top20Mutex = &sync.Mutex{}

func (h *Handler) weeklyTop20(rw http.ResponseWriter, req *http.Request) {
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
	if week < 0 || week > 12 {
		week = 0
	}

	// was there a forceOverwrite given?
	forceOverwrite := false
	value := req.URL.Query().Get("forceOverwrite")
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
	if !forceOverwrite && top20.IsAvailable(seasonID, week) {
		http.ServeFile(rw, req, top20.Filename(seasonID, week))
		return
	}
	// lock global mutex
	top20Mutex.Lock()
	defer top20Mutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top20.IsAvailable(seasonID, week) {
		http.ServeFile(rw, req, top20.Filename(seasonID, week))
		return
	}

	// create/update top20 image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, results, err := h.getWeek(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek results: %v", err)
		h.failure(rw, req, err)
		return
	}
	hm := top20.New(season, raceweek, track, results)
	if err := hm.Draw(); err != nil {
		log.Errorf("could not create top20: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, top20.Filename(seasonID, week))
}
