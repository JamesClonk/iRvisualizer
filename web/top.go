package web

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image/top"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
	"github.com/gorilla/mux"
)

var topMutex = &sync.Mutex{}

func (h *Handler) weeklyTopScores(rw http.ResponseWriter, req *http.Request) {
	image := "scores"

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
	if week < 1 || week > 13 {
		week = 1
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a topN given?
	topN := 20
	value := req.URL.Query().Get("topN")
	if len(value) > 0 {
		topN, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert topN [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a headerless given?
	headerless := false
	value = req.URL.Query().Get("headerless")
	if len(value) > 0 {
		headerless, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("could not convert headerless [%s] to bool: %v", value, err)
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
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}

	// create/update top image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("top scores: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		raceweek.LastUpdate = time.Now()
		track.Name = "starting soon..."
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("top scores: could not get raceweek summaries for season[%d], week[%d]: %v", seasonID, week-1, err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// champ points
	champ := top.DataSet{
		Title: "Highest Championship Points",
		Icons: "star",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by champ points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].HighestChampPoints > summaries[j].HighestChampPoints
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		champ.Rows = append(champ.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].HighestChampPoints),
		})
	}
	data = append(data, champ)

	// club points
	club := top.DataSet{
		Title: "Total Club Points contributed",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by club points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalClubPoints > summaries[j].TotalClubPoints
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		club.Rows = append(club.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].TotalClubPoints),
		})
	}
	data = append(data, club)

	// podiums
	podiums := top.DataSet{
		Title: "Podium Positions",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by podiums
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Podiums > summaries[j].Podiums
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		podiums.Rows = append(podiums.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].Podiums),
		})
	}
	data = append(data, podiums)

	hm := top.New(colorScheme, image, season, raceweek, track, data)
	if err := hm.Draw(headerless); err != nil {
		log.Errorf("could not create weekly top [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, top.Filename(image, seasonID, week))
}

func (h *Handler) weeklyTopRacers(rw http.ResponseWriter, req *http.Request) {
	image := "racers"

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
	if week < 1 || week > 13 {
		week = 1
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a topN given?
	topN := 20
	value := req.URL.Query().Get("topN")
	if len(value) > 0 {
		topN, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert topN [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a headerless given?
	headerless := false
	value = req.URL.Query().Get("headerless")
	if len(value) > 0 {
		headerless, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("could not convert headerless [%s] to bool: %v", value, err)
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
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}

	// create/update top image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("top racers: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		raceweek.LastUpdate = time.Now()
		track.Name = "starting soon..."
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("top racers: could not get raceweek summaries for season[%d], week[%d]: %v", seasonID, week-1, err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// top5 positions
	top5 := top.DataSet{
		Title: "Top5 Hype (Finishing Positions)",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by top5
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Top5 > summaries[j].Top5
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		top5.Rows = append(top5.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].Top5),
		})
	}
	data = append(data, top5)

	// positions-gained
	positions := top.DataSet{
		Title: "Positions gained / Hard Charger",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by positions-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalPositionsGained > summaries[j].TotalPositionsGained
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		value := fmt.Sprintf("%d", summaries[i].TotalPositionsGained)
		if summaries[i].TotalPositionsGained > 0 {
			value = "+" + value
		}
		positions.Rows = append(positions.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  value,
		})
	}
	data = append(data, positions)

	// races driven
	races := top.DataSet{
		Title: "Most Races (+ Laps Completed)",
		Icons: "flag",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by races
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].NumberOfRaces > summaries[j].NumberOfRaces
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		races.Rows = append(races.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d (%d)", summaries[i].NumberOfRaces, summaries[i].LapsCompleted),
		})
	}
	data = append(data, races)

	hm := top.New(colorScheme, image, season, raceweek, track, data)
	if err := hm.Draw(headerless); err != nil {
		log.Errorf("could not create weekly top [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, top.Filename(image, seasonID, week))
}

func (h *Handler) weeklyTopLaps(rw http.ResponseWriter, req *http.Request) {
	image := "laps"

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
	if week < 1 || week > 13 {
		week = 1
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a topN given?
	topN := 20
	value := req.URL.Query().Get("topN")
	if len(value) > 0 {
		topN, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert topN [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a headerless given?
	headerless := false
	value = req.URL.Query().Get("headerless")
	if len(value) > 0 {
		headerless, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("could not convert headerless [%s] to bool: %v", value, err)
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
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}

	// create/update top image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("top laps: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		raceweek.LastUpdate = time.Now()
		track.Name = "starting soon..."
	}
	// summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	// if err != nil {
	// 	log.Errorf("top laps: could not get raceweek summaries for season[%d], week[%d]: %v", seasonID, week-1, err)
	// 	h.failure(rw, req, err)
	// 	return
	// }
	timeRankings, err := h.getRaceWeekTimeRankings(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek timerankings: %v", err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// tt lap
	tt := top.DataSet{
		Title: "Fastest Time Trial Session",
		Icons: "clock",
		Rows:  make([]top.DataSetRow, 0),
	}
	// filter by > 100
	filtered := make([]database.TimeRanking, 0)
	for _, ranking := range timeRankings {
		if ranking.TimeTrial > 100 {
			filtered = append(filtered, ranking)
		}
	}
	// sort by tt lap
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TimeTrial < filtered[j].TimeTrial
	})
	for i := 0; i < topN && i < len(filtered); i++ {
		icon := ""
		//if (((filtered[i].TimeTrial) - (filtered[i].TimeTrialFastestLap)) / 10) < 150 { // if smaller than 150ms
		if ((filtered[i].TimeTrial-filtered[i].TimeTrialFastestLap)/10) < (filtered[i].TimeTrial/7777) &&
			(filtered[i].TimeTrial-filtered[i].TimeTrialFastestLap) > 0 {
			icon = "fire"
		}
		tt.Rows = append(tt.Rows, top.DataSetRow{
			Driver: filtered[i].Driver.Name,
			Icon:   icon,
			Value:  util.ConvertLaptime(filtered[i].TimeTrial),
		})
	}
	data = append(data, tt)

	// tt fastest lap
	ttf := top.DataSet{
		Title: "Fastest Time Trial Lap",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by tt fastest lap
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].TimeTrialFastestLap < filtered[j].TimeTrialFastestLap
	})
	for i := 0; i < topN && i < len(filtered); i++ {
		icon := ""
		if i+1 < len(filtered) &&
			filtered[i+1].TimeTrialFastestLap-filtered[i].TimeTrialFastestLap > filtered[i].TimeTrialFastestLap/222 {
			icon = "green_arrow"
		}
		ttf.Rows = append(ttf.Rows, top.DataSetRow{
			Driver: filtered[i].Driver.Name,
			Icon:   icon,
			Value:  util.ConvertLaptime(filtered[i].TimeTrialFastestLap),
		})
	}
	data = append(data, ttf)

	// race lap
	race := top.DataSet{
		Title: "Fastest Race Lap",
		Rows:  make([]top.DataSetRow, 0),
	}
	// filter by > 100
	filtered = make([]database.TimeRanking, 0)
	for _, ranking := range timeRankings {
		if ranking.Race > 100 {
			filtered = append(filtered, ranking)
		}
	}
	// sort by race lap
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Race < filtered[j].Race
	})
	for i := 0; i < topN && i < len(filtered); i++ {
		icon := ""
		if i+1 < len(filtered) &&
			filtered[i+1].Race-filtered[i].Race > filtered[i].Race/333 {
			icon = "blue_arrow"
		}
		race.Rows = append(race.Rows, top.DataSetRow{
			Driver: filtered[i].Driver.Name,
			Icon:   icon,
			Value:  util.ConvertLaptime(filtered[i].Race),
		})
	}
	data = append(data, race)

	// // laps
	// laps := top.DataSet{
	// 	Title: "Laps completed",
	// 	Rows:  make([]top.DataSetRow, 0),
	// }
	// // sort by laps
	// sort.Slice(summaries, func(i, j int) bool {
	// 	return summaries[i].LapsCompleted > summaries[j].LapsCompleted
	// })
	// for i := 0; i < topN && i < len(summaries); i++ {
	// 	laps.Rows = append(laps.Rows, top.DataSetRow{
	// 		Driver: summaries[i].Driver.Name,
	// 		Value:  fmt.Sprintf("%d", summaries[i].LapsCompleted),
	// 	})
	// }
	// data = append(data, laps)

	hm := top.New(colorScheme, image, season, raceweek, track, data)
	if err := hm.Draw(headerless); err != nil {
		log.Errorf("could not create weekly top [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, top.Filename(image, seasonID, week))
}

func (h *Handler) weeklyTopSafety(rw http.ResponseWriter, req *http.Request) {
	image := "safety"

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
	if week < 1 || week > 13 {
		week = 1
	}

	// was there a colorScheme given?
	colorScheme := req.URL.Query().Get("colorScheme")

	// was there a topN given?
	topN := 20
	value := req.URL.Query().Get("topN")
	if len(value) > 0 {
		topN, err = strconv.Atoi(value)
		if err != nil {
			log.Errorf("could not convert topN [%s] to int: %v", value, err)
			h.failure(rw, req, err)
			return
		}
	}

	// was there a headerless given?
	headerless := false
	value = req.URL.Query().Get("headerless")
	if len(value) > 0 {
		headerless, err = strconv.ParseBool(value)
		if err != nil {
			log.Errorf("could not convert headerless [%s] to bool: %v", value, err)
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
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(colorScheme, image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}

	// create/update top image
	season, err := h.getSeason(seasonID)
	if err != nil {
		log.Errorf("could not get season: %v", err)
		h.failure(rw, req, err)
		return
	}
	raceweek, track, err := h.getRaceWeek(seasonID, week-1)
	if err != nil {
		log.Debugf("top safety: could not get raceweek for season[%d], week[%d]: %v", seasonID, week-1, err)
		raceweek.RaceWeek = week - 1
		raceweek.LastUpdate = time.Now()
		track.Name = "starting soon..."
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("top safety: could not get raceweek summaries for season[%d], week[%d]: %v", seasonID, week-1, err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// irating-gained
	irating := top.DataSet{
		Title: "Total iRating gained",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by irating-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalIRatingGain > summaries[j].TotalIRatingGain
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		value := fmt.Sprintf("%d", summaries[i].TotalIRatingGain)
		if summaries[i].TotalIRatingGain > 0 {
			value = "+" + value
		}
		irating.Rows = append(irating.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  value,
		})
	}
	data = append(data, irating)

	// sr-gained
	sr := top.DataSet{
		Title: "Total Safety Rating gained",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by sr-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalSafetyRatingGain > summaries[j].TotalSafetyRatingGain
	})
	for i := 0; i < topN && i < len(summaries); i++ {
		value := fmt.Sprintf("%.2f", float64(summaries[i].TotalSafetyRatingGain)/float64(100))
		if summaries[i].TotalSafetyRatingGain > 0 {
			value = "+" + value
		}
		sr.Rows = append(sr.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  value,
		})
	}
	data = append(data, sr)

	// inc/lap
	inc := top.DataSet{
		Title: "Avg. Incidents per Lap (min. 3 races)",
		Icons: "safety",
		Rows:  make([]top.DataSetRow, 0),
	}
	// filter by min. 3 races
	filtered := summaries[:0]
	for _, summary := range summaries {
		if summary.NumberOfRaces >= 3 {
			filtered = append(filtered, summary)
		}
	}
	// sort by inc/lap
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].AverageIncidentsPerLap < filtered[j].AverageIncidentsPerLap
	})
	for i := 0; i < topN && i < len(filtered); i++ {
		inc.Rows = append(inc.Rows, top.DataSetRow{
			Driver: filtered[i].Driver.Name,
			Value:  fmt.Sprintf("%.3f", filtered[i].AverageIncidentsPerLap),
		})
	}
	data = append(data, inc)

	hm := top.New(colorScheme, image, season, raceweek, track, data)
	if err := hm.Draw(headerless); err != nil {
		log.Errorf("could not create weekly top [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, top.Filename(image, seasonID, week))
}
