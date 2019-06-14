package web

// import (
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/JamesClonk/iRvisualizer/heatmap"
// 	"github.com/JamesClonk/iRvisualizer/log"
// 	"github.com/JamesClonk/iRvisualizer/util"
// 	"github.com/gorilla/mux"
// )

// func (h *Handler) weeklyTop20(rw http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	seasonID, err := strconv.Atoi(vars["seasonID"])
// 	if err != nil {
// 		log.Errorf("could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
// 		h.failure(rw, req, err)
// 		return
// 	}
// 	if seasonID < 2000 || seasonID > 9999 {
// 		seasonID = 2377
// 	}
// 	week, err := strconv.Atoi(vars["week"])
// 	if err != nil {
// 		log.Errorf("could not convert week [%s] to int: %v", vars["week"], err)
// 		h.failure(rw, req, err)
// 		return
// 	}
// 	if week < 0 || week > 12 {
// 		week = 0
// 	}

// 	// was there a forceOverwrite given?
// 	forceOverwrite := false
// 	value = req.URL.Query().Get("forceOverwrite")
// 	if len(value) > 0 {
// 		forceOverwrite, err = strconv.ParseBool(value)
// 		if err != nil {
// 			log.Errorf("could not convert forceOverwrite [%s] to bool: %v", value, err)
// 			h.failure(rw, req, err)
// 			return
// 		}
// 	}

// 	// check if file already exists
// 	heatmapFilename := heatmap.HeatmapFilename(seasonID, week)
// 	metaFilename := heatmap.MetadataFilename(seasonID, week)
// 	if !forceOverwrite && util.FileExists(metaFilename) && util.FileExists(heatmapFilename) {
// 		metadata := heatmap.GetMetadata(metaFilename)
// 		// if it's older than 2 hours
// 		if (time.Now().Sub(metadata.LastUpdated) < time.Hour*2) ||
// 			// or if it's from a week longer than 10 days ago and updated somewhere within 10 days after weekstart
// 			(time.Now().Sub(metadata.StartDate.AddDate(0, 0, metadata.Week*7)) > time.Hour*24*10 &&
// 				metadata.LastUpdated.Sub(metadata.StartDate.AddDate(0, 0, metadata.Week*7)) > time.Hour*24*10) {
// 			// serve image immediately
// 			http.ServeFile(rw, req, heatmapFilename)
// 			return
// 		}
// 	}

// 	// create/update heatmap image
// 	season, err := h.getSeason(seasonID)
// 	if err != nil {
// 		log.Errorf("could not get season: %v", err)
// 		h.failure(rw, req, err)
// 		return
// 	}
// 	raceweek, track, results, err := h.getWeek(seasonID, week-1)
// 	if err != nil {
// 		log.Errorf("could not get raceweek results: %v", err)
// 		h.failure(rw, req, err)
// 		return
// 	}
// 	hm := heatmap.New(season, raceweek, track, results)
// 	if err := hm.Draw(minSOF, maxSOF, true); err != nil {
// 		log.Errorf("could not create heatmap: %v", err)
// 		h.failure(rw, req, err)
// 		return
// 	}

// 	// serve new/updated image
// 	http.ServeFile(rw, req, heatmapFilename)
// }
