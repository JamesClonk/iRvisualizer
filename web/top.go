package web

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/JamesClonk/iRvisualizer/image/top"
	"github.com/JamesClonk/iRvisualizer/log"
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
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
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
		log.Errorf("could not get raceweek: %v", err)
		h.failure(rw, req, err)
		return
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek summaries: %v", err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// champ points
	champ := top.DataSet{
		Title: "Highest Championship Points",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by champ points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].HighestChampPoints > summaries[j].HighestChampPoints
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
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
	for i := 0; i < 25 && i < len(summaries); i++ {
		club.Rows = append(club.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].TotalClubPoints),
		})
	}
	data = append(data, club)

	// podiums
	podiums := top.DataSet{
		Title: "Podium positions",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by podiums
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Podiums > summaries[j].Podiums
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		podiums.Rows = append(podiums.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].Podiums),
		})
	}
	data = append(data, podiums)

	hm := top.New(image, season, raceweek, track, data)
	if err := hm.Draw(); err != nil {
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
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
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
		log.Errorf("could not get raceweek: %v", err)
		h.failure(rw, req, err)
		return
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek summaries: %v", err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// irating-gained
	irating := top.DataSet{
		Title: "TODO: TOP5-HYPE",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by irating-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].HighestIRatingGain > summaries[j].HighestIRatingGain
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		irating.Rows = append(irating.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].HighestIRatingGain),
		})
	}
	data = append(data, irating)

	// positions-gained
	positions := top.DataSet{
		Title: "Positions gained / Hard Charger",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by positions-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalPositionsGained > summaries[j].TotalPositionsGained
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		value := fmt.Sprintf("%d", summaries[i].TotalPositionsGained)
		if summaries[i].TotalPositionsGained > 0 {
			value = "+" + value
		} else {
			value = "-" + value
		}
		positions.Rows = append(positions.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  value,
		})
	}
	data = append(data, positions)

	// sr-gained
	sr := top.DataSet{
		Title: "TODO: NOF RACES",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by sr-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalSafetyRatingGain > summaries[j].TotalSafetyRatingGain
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		sr.Rows = append(sr.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].TotalSafetyRatingGain),
		})
	}
	data = append(data, sr)

	hm := top.New(image, season, raceweek, track, data)
	if err := hm.Draw(); err != nil {
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
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
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
		log.Errorf("could not get raceweek: %v", err)
		h.failure(rw, req, err)
		return
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek summaries: %v", err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// champ points
	champ := top.DataSet{
		Title: "TODO: RACE LAPS",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by champ points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].HighestChampPoints > summaries[j].HighestChampPoints
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		champ.Rows = append(champ.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].HighestChampPoints),
		})
	}
	data = append(data, champ)

	// club points
	club := top.DataSet{
		Title: "TODO: TT LAPS",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by club points
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalClubPoints > summaries[j].TotalClubPoints
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		club.Rows = append(club.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].TotalClubPoints),
		})
	}
	data = append(data, club)

	// laps
	laps := top.DataSet{
		Title: "Laps completed",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by laps
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].LapsCompleted > summaries[j].LapsCompleted
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		laps.Rows = append(laps.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  fmt.Sprintf("%d", summaries[i].LapsCompleted),
		})
	}
	data = append(data, laps)

	hm := top.New(image, season, raceweek, track, data)
	if err := hm.Draw(); err != nil {
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
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
		http.ServeFile(rw, req, top.Filename(image, seasonID, week))
		return
	}
	// lock global mutex
	topMutex.Lock()
	defer topMutex.Unlock()
	// doublecheck, to make sure it wasn't updated by now by another goroutine that held the lock before
	if !forceOverwrite && top.IsAvailable(image, seasonID, week) {
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
		log.Errorf("could not get raceweek: %v", err)
		h.failure(rw, req, err)
		return
	}
	summaries, err := h.getRaceWeekSummaries(seasonID, week-1)
	if err != nil {
		log.Errorf("could not get raceweek summaries: %v", err)
		h.failure(rw, req, err)
		return
	}

	data := make([]top.DataSet, 0)
	// irating-gained
	irating := top.DataSet{
		Title: "Highest iRating gained",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by irating-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].HighestIRatingGain > summaries[j].HighestIRatingGain
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		value := fmt.Sprintf("%d", summaries[i].HighestIRatingGain)
		if summaries[i].HighestIRatingGain > 0 {
			value = "+" + value
		} else {
			value = "-" + value
		}
		irating.Rows = append(irating.Rows, top.DataSetRow{
			Driver: summaries[i].Driver.Name,
			Value:  value,
		})
	}
	data = append(data, irating)

	// sr-gained
	sr := top.DataSet{
		Title: "Total Safety Rating",
		Rows:  make([]top.DataSetRow, 0),
	}
	// sort by sr-gained
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalSafetyRatingGain > summaries[j].TotalSafetyRatingGain
	})
	for i := 0; i < 25 && i < len(summaries); i++ {
		value := fmt.Sprintf("%.2f", float64(summaries[i].TotalSafetyRatingGain)/float64(100))
		if summaries[i].TotalSafetyRatingGain > 0 {
			value = "+" + value
		} else {
			value = "-" + value
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
	for i := 0; i < 25 && i < len(filtered); i++ {
		inc.Rows = append(inc.Rows, top.DataSetRow{
			Driver: filtered[i].Driver.Name,
			Value:  fmt.Sprintf("%.3f", filtered[i].AverageIncidentsPerLap),
		})
	}
	data = append(data, inc)

	hm := top.New(image, season, raceweek, track, data)
	if err := hm.Draw(); err != nil {
		log.Errorf("could not create weekly top [%s]: %v", image, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated image
	http.ServeFile(rw, req, top.Filename(image, seasonID, week))
}
