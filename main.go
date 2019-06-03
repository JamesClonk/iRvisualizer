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
	"github.com/robfig/cron"
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

	r.HandleFunc("/season/{seasonID}/week/{week}/heatmap.png", h.weeklyHeatmap).Methods("GET")
	r.HandleFunc("/season/{seasonID}/heatmap.png", h.seasonalHeatmap).Methods("GET")
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
	heatmapFilename := heatmap.HeatmapFilename(seasonID, -1)
	metaFilename := heatmap.MetadataFilename(seasonID, -1)
	if !forceOverwrite && util.FileExists(metaFilename) && util.FileExists(heatmapFilename) {
		metadata := heatmap.GetMetadata(metaFilename)
		// if it's older than 2 hours
		if (time.Now().Sub(metadata.LastUpdated) < time.Hour*2) ||
			// or if it's from a week longer than 10 days ago and updated somewhere within 10 days after weekstart
			(time.Now().Sub(metadata.StartDate.AddDate(0, 0, 12*7)) > time.Hour*24*10 &&
				metadata.LastUpdated.Sub(metadata.StartDate.AddDate(0, 0, 12*7)) > time.Hour*24*10) {
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
	for week := 0; week < 12; week++ {
		_, _, rs, err := h.getWeek(seasonID, week)
		if err != nil {
			log.Errorf("could not get raceweek results: %v", err)
			h.failure(rw, req, err)
			return
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
		start := database.WeekStart(season.StartDate.UTC().AddDate(0, 0, (week+1)*7).Add(-1 * time.Minute))
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
	// divide by 12 weeks
	for i := range results {
		results[i].StrengthOfField = results[i].StrengthOfField / 12
		results[i].SizeOfField = results[i].SizeOfField / 12
	}
	// sort again to be on the safe side..
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.Before(results[j].StartTime)
	})

	hm := heatmap.New(season, database.RaceWeek{RaceWeek: -1}, database.Track{}, results)
	if err := hm.Draw(maxSOF); err != nil {
		log.Errorf("could not create seasonal heatmap: %v", err)
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
