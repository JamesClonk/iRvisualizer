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
	// collect champ & TT points for all weeks
	var weeks int
	ccPoints := make(map[database.Driver][]float64)
	ttPoints := make(map[database.Driver][]int)
	for week := 0; week < 13; week++ { // allow for leap seasons with 13 official weeks, like 2020S3
		weeklyCcPoints, err := h.getChampPoints(seasonID, week)
		if err != nil {
			log.Errorf("could not get championship points for week [%d]: %v", week+1, err)
			h.failure(rw, req, err)
			return
		}
		weeklyTtResults, err := h.getTTStandings(seasonID, week)
		if err != nil {
			log.Errorf("could not get TT results for week [%d]: %v", week+1, err)
			h.failure(rw, req, err)
			return
		}

		// do we have data for this week?
		if len(weeklyCcPoints) > 0 || len(weeklyTtResults) > 0 {
			weeks++
		}

		// collect champpoints for all drivers
		if len(weeklyCcPoints) > 0 {
			drivers := make(map[database.Driver][]int)
			for _, p := range weeklyCcPoints {
				if _, ok := drivers[p.Driver]; !ok {
					drivers[p.Driver] = make([]int, 0)
				}
				drivers[p.Driver] = append(drivers[p.Driver], p.ChampPoints)
			}
			// figure out points for each driver this week
			for driver, values := range drivers {
				resultCount := int(math.Ceil(float64(len(values)) / 4))
				var result float64
				for i := 0; i < resultCount; i++ {
					result += float64(values[i])
				}
				// final result / average
				if _, ok := ccPoints[driver]; !ok {
					ccPoints[driver] = make([]float64, 0)
				}
				ccPoints[driver] = append(ccPoints[driver], (result / float64(resultCount)))
			}
		}

		// collect TT points for all drivers
		if len(weeklyTtResults) > 0 {
			for _, tt := range weeklyTtResults {
				if _, ok := ttPoints[tt.Driver]; !ok {
					ttPoints[tt.Driver] = make([]int, 0)
				}
				ttPoints[tt.Driver] = append(ttPoints[tt.Driver], tt.Points)
			}
		}
	}
	bestN := weeks - int(math.Floor(float64(weeks)/3)) // how many weeks to count so far? (removes dropweeks)

	// total bestN values
	champData := make([]ranking.DataRow, 0)
	for driver, values := range ccPoints {
		sort.Slice(values, func(i, j int) bool {
			return values[i] > values[j]
		})
		var total float64
		for n := 0; n < bestN && n < len(values); n++ {
			total += values[n]
		}
		champData = append(champData, ranking.DataRow{
			Driver: driver.Name,
			Value:  fmt.Sprintf("%d", int(math.Floor(total))),
		})
	}
	// sort by values
	sort.Slice(champData, func(i, j int) bool {
		return champData[i].Driver < champData[j].Driver
	})
	sort.Slice(champData, func(i, j int) bool {
		a, _ := strconv.Atoi(champData[i].Value)
		b, _ := strconv.Atoi(champData[j].Value)
		return a > b
	})
	ttData := make([]ranking.DataRow, 0)
	for driver, values := range ttPoints {
		sort.Slice(values, func(i, j int) bool {
			return values[i] > values[j]
		})
		var total int
		for n := 0; n < bestN && n < len(values); n++ {
			total += values[n]
		}
		ttData = append(ttData, ranking.DataRow{
			Driver: driver.Name,
			Value:  fmt.Sprintf("%d", total),
		})
	}
	// sort by values
	sort.Slice(ttData, func(i, j int) bool {
		return ttData[i].Driver < ttData[j].Driver
	})
	sort.Slice(ttData, func(i, j int) bool {
		a, _ := strconv.Atoi(ttData[i].Value)
		b, _ := strconv.Atoi(ttData[j].Value)
		return a > b
	})

	r := ranking.New(season, champData, ttData)
	if err := r.Draw(req.URL.Query().Get("colorScheme"), bestN, weeks); err != nil {
		log.Errorf("could not create season ranking: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, ranking.Filename(seasonID))
}
