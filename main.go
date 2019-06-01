package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/env"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

const (
	borderSize     = float64(7)
	imageHeight    = float64(480)
	imageLength    = float64(1024)
	headerHeight   = float64(45)
	timeslotHeight = float64(50)
	dayLength      = float64(160)
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

	r.HandleFunc("/heatmap/season/{seasonID}/week/{week}", h.heatmap).Methods("GET")
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
	if !h.verifyBasicAuth(rw, req) {
		return
	}

	vars := mux.Vars(req)
	seasonID, err := strconv.Atoi(vars["seasonID"])
	if err != nil {
		log.Errorf("could not convert seasonID [%s] to int: %v", vars["seasonID"], err)
		h.failure(rw, req, err)
		return
	}
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		log.Errorf("could not convert week [%s] to int: %v", vars["week"], err)
		h.failure(rw, req, err)
		return
	}

	// check if file already exists
	filename := fmt.Sprintf("public/heatmaps/season_%d_week_%d.png", seasonID, week)
	if !fileExists(filename) {
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

		maxSOF := 2500
		if err := drawHeatmap(season, raceweek, track, results, maxSOF); err != nil {
			log.Errorf("could not create heatmap: %v", err)
			h.failure(rw, req, err)
			return
		}
	}
	http.ServeFile(rw, req, filename)
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

func getResult(slot time.Time, results []database.RaceWeekResult) database.RaceWeekResult {
	sessions := make([]database.RaceWeekResult, 0)
	for _, result := range results {
		if result.StartTime.UTC() == slot.UTC() {
			sessions = append(sessions, result)
		}
	}

	// summarize splits
	result := database.RaceWeekResult{
		SizeOfField:     0,
		StrengthOfField: 0,
	}
	for _, session := range sessions {
		result.Official = session.Official
		result.SizeOfField += session.SizeOfField
		if session.StrengthOfField > result.StrengthOfField {
			result.StrengthOfField = session.StrengthOfField
		}
	}
	return result
}

func mapValueIntoRange(rangeStart, rangeEnd, min, max, value int) int {
	if value <= min {
		value = min + 1
	}
	if value >= max {
		return rangeEnd
	}
	rangeSize := rangeEnd - rangeStart
	return rangeStart + int((float64(value-min)/float64(max-min))*float64(rangeSize))
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func drawHeatmap(season database.Season, week database.RaceWeek, track database.Track, results []database.RaceWeekResult, maxSOF int) error {
	// heatmap titles, season + track
	heatmapTitle := fmt.Sprintf("%s - Week %d", season.SeasonName, week.RaceWeek+1)
	heatmap2ndTitle := track.Name
	// if len(track.Config) > 0 {
	// 	heatmap2ndTitle = fmt.Sprintf("%s - %s", track.Name, track.Config)
	// }

	log.Infof("draw heatmap for [%s] - [%s]", heatmapTitle, heatmap2ndTitle)

	// figure out timeslots
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := p.Parse(season.Timeslots)
	if err != nil {
		return fmt.Errorf("could not parse timeslot [%s] to crontab format: %v", season.Timeslots, err)
	}
	days := 7 // pretty sure that's never gonna change..
	// start -1 minute to previous day, to make sure schedule.Next will catch a midnight start (00:00)
	start := database.WeekStart(season.StartDate.UTC().AddDate(0, 0, (week.RaceWeek+1)*days).Add(-1 * time.Minute))
	timeslots := make([]time.Time, 0)
	next := schedule.Next(start)                             // get first timeslot
	weekStart := next                                        // first timeslot is our week start
	for next.Before(schedule.Next(start.AddDate(0, 0, 1))) { // collect all timeslots of 1 day
		timeslots = append(timeslots, next)
		next = schedule.Next(next)
	}
	// figure out dynamic SOF
	minSOF := 1000
	if maxSOF == 0 {
		maxSOF = minSOF * 2
		for _, result := range results {
			if result.StrengthOfField > maxSOF {
				maxSOF = result.StrengthOfField
			}
		}
	}

	// create canvas
	dc := gg.NewContext(int(imageLength), int(imageHeight))

	// background
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, imageLength, headerHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()
	dc.DrawRectangle(imageLength/2+dayLength/2, 0, imageLength/2, headerHeight)
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.Fill()

	// draw season name
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 20); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(heatmapTitle, dayLength/4, headerHeight/2, 0, 0.5)
	// draw track config
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 20); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(heatmap2ndTitle, imageLength-dayLength/4, headerHeight/2, 1, 0.5)

	// timeslots
	dc.DrawRectangle(0, headerHeight, dayLength, timeslotHeight)
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.Fill()
	if err := dc.LoadFontFace("public/fonts/roboto-mono_thin.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(0, 0, 0) // black
	dc.DrawStringAnchored("UTC / GMT+0", dayLength/2, headerHeight+timeslotHeight/2, 0.5, 0.5)
	if err := dc.LoadFontFace("public/fonts/roboto-mono_medium.ttf", 16); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	timeslotLength := ((imageLength - dayLength) / float64(len(timeslots))) - 1
	for slot := 0; slot < len(timeslots); slot++ {
		dc.DrawRectangle((float64(slot)*(timeslotLength+1))+(dayLength+1), headerHeight, timeslotLength, timeslotHeight)
		if slot%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
		// draw timeslot starting time
		dc.SetRGB255(0, 0, 0) // black
		dc.DrawStringAnchored(
			timeslots[slot].Format("15:04"),
			(float64(slot)*(timeslotLength+1))+(dayLength+1)+(timeslotLength/2),
			headerHeight+timeslotHeight/2,
			0.5, 0.5)
	}

	// weekdays
	if err := dc.LoadFontFace("public/fonts/RobotoCondensed-Regular.ttf", 20); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dayHeight := ((imageHeight - headerHeight - timeslotHeight) / float64(days)) - 1
	for day := 0; day < days; day++ {
		dc.DrawRectangle(0, (float64(day)*(dayHeight+1))+(headerHeight+timeslotHeight+1), dayLength, dayHeight)
		if day%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
		// draw weekday name
		dc.SetRGB255(0, 0, 0) // black
		dc.DrawStringAnchored(
			weekStart.AddDate(0, 0, day).Weekday().String(),
			dayLength/2,
			(float64(day)*(dayHeight+1))+(headerHeight+timeslotHeight+1)+dayHeight/2,
			0.5, 0.5)
	}

	// events
	eventHeight := ((imageHeight - headerHeight - timeslotHeight) / float64(days)) - 1
	eventLength := ((imageLength - dayLength) / float64(len(timeslots))) - 1
	for day := 0; day < days; day++ {
		for slot := 0; slot < len(timeslots); slot++ {
			dc.DrawRectangle(
				(float64(slot)*(eventLength+1))+(dayLength+1),
				(float64(day)*(eventHeight+1))+(headerHeight+timeslotHeight+1),
				eventLength, eventHeight)
			dc.SetRGB255(255, 255, 255) // white
			dc.Fill()

			// draw event values
			timeslot := weekStart.AddDate(0, 0, day).Add(time.Hour * time.Duration(timeslots[slot].Hour()))
			result := getResult(timeslot, results)
			// only draw event if a session actually happened already
			if timeslot.Before(time.Now().Add(time.Hour * -2)) {
				sof := 0
				if result.Official {
					sof = result.StrengthOfField
					// draw background color
					dc.DrawRectangle(
						(float64(slot)*(eventLength+1))+(dayLength+1),
						(float64(day)*(eventHeight+1))+(headerHeight+timeslotHeight+1),
						eventLength, eventHeight)
					dc.SetRGBA255(0, 0, 240-mapValueIntoRange(0, 120, minSOF, maxSOF, sof), mapValueIntoRange(10, 200, minSOF, maxSOF, sof)) // sof color
					dc.Fill()
				}
				dc.SetRGB255(0, 0, 0) // black
				if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 15); err != nil {
					return fmt.Errorf("could not load font: %v", err)
				}
				dc.DrawStringAnchored(
					fmt.Sprintf("%d", sof),
					(float64(slot)*(eventLength+1))+(dayLength+1)+eventLength/2,
					(float64(day)*(eventHeight+1))+(headerHeight+timeslotHeight+1)+eventHeight/3,
					0.5, 0.5)
				if err := dc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 13); err != nil {
					return fmt.Errorf("could not load font: %v", err)
				}
				dc.DrawStringAnchored(
					fmt.Sprintf("%d", result.SizeOfField),
					(float64(slot)*(eventLength+1))+(dayLength+1)+eventLength/2,
					(float64(day)*(eventHeight+1))+(headerHeight+timeslotHeight+1)+eventHeight/1.5,
					0.5, 0.5)
			}
		}
	}

	// add border to image
	fdc := gg.NewContext(int(imageLength+borderSize*2), int(imageHeight+borderSize*2))
	fdc.SetRGB255(39, 39, 39) // dark gray 1
	fdc.Clear()
	fdc.DrawImage(dc.Image(), int(borderSize), int(borderSize))

	filename := fmt.Sprintf("public/heatmaps/season_%d_week_%d.png", season.SeasonID, week.RaceWeek+1)
	return fdc.SavePNG(filename)
}

/*
	Colors:
	dc.SetRGB255(0, 0, 0) // black
	dc.SetRGB255(39, 39, 39) // dark gray 1
	dc.SetRGB255(255, 255, 255) // white
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
