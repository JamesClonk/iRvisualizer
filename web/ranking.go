package web

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image/ranking"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/gorilla/mux"
)

var rankingMutex = &sync.Mutex{}

func (h *Handler) ranking(rw http.ResponseWriter, req *http.Request) {
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
	if !forceOverwrite && ranking.IsAvailable(seasonID) {
		http.ServeFile(rw, req, ranking.Filename(seasonID))
		return
	}
	// lock global mutex
	rankingMutex.Lock()
	defer rankingMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && ranking.IsAvailable(seasonID) {
		http.ServeFile(rw, req, ranking.Filename(seasonID))
		return
	}

	// create/update ranking image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	// collect all weeks
	points := make(map[database.Driver]int)
	for week := 0; week < 12; week++ {
		weeklyPoints, err := h.getPoints(seasonID, week)
		if err != nil {
			log.Errorf("could not get points for week [%d]: %v", week+1, err)
			h.failure(rw, req, err)
			return
		}

		// // collect points for all drivers
		drivers := make(map[database.Driver][]int)
		for _, p := range weeklyPoints {
			if _, ok := drivers[p.Driver]; !ok {
				drivers[p.Driver] = make([]int, 0)
			}
			drivers[p.Driver] = append(drivers[p.Driver], p.ChampPoints)
		}
		// figure out points for each driver this week
		for driver, values := range drivers {
			resultCount := int(math.Ceil(float64(len(values)) / 4))
			var result int
			for i := 0; i < resultCount; i++ {
				result += values[i]
			}

			// final result / average
			if _, ok := points[driver]; !ok {
				points[driver] = 0
			}
			points[driver] = points[driver] + int(math.RoundToEven(float64(result)/float64(resultCount)))
		}
	}

	// sort by values
	data := make([]ranking.DataRow, 0)
	for driver, value := range points {
		data = append(data, ranking.DataRow{
			Driver: driver.Name,
			Value:  fmt.Sprintf("%d", value),
		})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Driver < data[j].Driver
	})
	sort.Slice(data, func(i, j int) bool {
		a, _ := strconv.Atoi(data[i].Value)
		b, _ := strconv.Atoi(data[j].Value)
		return a > b
	})

	hm := ranking.New(season, data)
	if err := hm.Draw(req.URL.Query().Get("colorScheme"), 0, 99); err != nil {
		log.Errorf("could not create season ranking: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, ranking.Filename(seasonID))
}
