package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/env"
	"github.com/JamesClonk/iRvisualizer/heatmap"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
	"github.com/gorilla/mux"
)

type Handler struct {
	Username string
	Password string
	DB       database.Database
	Mutex    *sync.Mutex
}

func main() {
	port := env.Get("PORT", "8080")
	level := env.Get("LOG_LEVEL", "info")
	username := env.MustGet("AUTH_USERNAME")
	password := env.MustGet("AUTH_PASSWORD")

	log.Infoln("port:", port)
	log.Infoln("log level:", level)
	log.Infoln("auth username:", username)

	// setup database connection
	db := database.NewDatabase(database.NewAdapter())

	// global handler
	h := &Handler{
		Username: username,
		Password: password,
		DB:       db,
		Mutex:    &sync.Mutex{},
	}

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, router(h)))
}

func router(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/health").HandlerFunc(h.health)

	r.HandleFunc("/season/{seasonID}/week/{week}/heatmap.png", h.heatmap).Methods("GET")
	return r
}

func (h *Handler) failure(rw http.ResponseWriter, req *http.Request, err error) {
	rw.WriteHeader(500)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(fmt.Sprintf(`{ "error": "%v" }`, err.Error())))
}

func (h *Handler) health(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{ "status": "ok" }`))
}

func (h *Handler) heatmap(rw http.ResponseWriter, req *http.Request) {
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

	// was there a maxSOF given?
	maxSOF := 2500
	value := req.URL.Query().Get("maxSOF")
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

	// check if file already exists
	heatmapFilename := heatmap.HeatmapFilename(seasonID, week)
	metaFilename := heatmap.MetadataFilename(seasonID, week)
	if !forceOverwrite && util.FileExists(metaFilename) && util.FileExists(heatmapFilename) {
		metadata := heatmap.GetMetadata(metaFilename)
		// if it's older than 2 hours
		if (time.Now().Sub(metadata.LastUpdated) < time.Hour*2) ||
			// or if it's from a week longer than 10 days ago and updated somewhere within 10 days after weekstart
			(time.Now().Sub(metadata.StartDate.AddDate(0, 0, metadata.Week*7)) > time.Hour*24*10 &&
				metadata.LastUpdated.Sub(metadata.StartDate.AddDate(0, 0, metadata.Week*7)) > time.Hour*24*10) {
			// serve image immediately
			http.ServeFile(rw, req, heatmapFilename)
			return
		}
	}

	// create/update heatmap image
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
	hm := heatmap.New(season, raceweek, track, results)
	if err := hm.Draw(maxSOF); err != nil {
		log.Errorf("could not create heatmap: %v", err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, heatmapFilename)
}

func (h *Handler) verifyBasicAuth(rw http.ResponseWriter, req *http.Request) bool {
	user, pw, ok := req.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(h.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pw), []byte(h.Password)) != 1 {
		rw.Header().Set("WWW-Authenticate", `Basic realm="iRcollector"`)
		rw.WriteHeader(401)
		rw.Write([]byte("Unauthorized"))
		return false
	}
	return true
}

func (h *Handler) getSeason(seasonID int) (database.Season, error) {
	log.Infof("collect season [%d]", seasonID)
	return h.DB.GetSeasonByID(seasonID)
}

func (h *Handler) getWeek(seasonID, week int) (database.RaceWeek, database.Track, []database.RaceWeekResult, error) {
	log.Infof("collect results for season [%d], week [%d]", seasonID, week)

	raceweek, err := h.DB.GetRaceWeekBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return database.RaceWeek{}, database.Track{}, nil, err
	}

	track, err := h.DB.GetTrackByID(raceweek.TrackID)
	if err != nil {
		return raceweek, database.Track{}, nil, err
	}

	results, err := h.DB.GetRaceWeekResultsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return raceweek, track, nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.Before(results[j].StartTime)
	})

	return raceweek, track, results, nil
}
