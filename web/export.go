package web

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/web/csv"
	"github.com/gorilla/mux"
)

func (h *Handler) series(rw http.ResponseWriter, req *http.Request) {
	// get all active series
	series, err := h.getSeries()
	if err != nil {
		log.Errorf("could not get series: %v", err)
		h.failure(rw, req, err)
		return
	}

	_, _ = rw.Write([]byte("SERIES_ID;SERIES_NAME\n"))
	for _, series := range series {
		_, _ = rw.Write([]byte(fmt.Sprintf("%d;%s\n", series.SeriesID, series.SeriesName)))
	}
}

func (h *Handler) seriesWeeklyExport(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	seriesID, err := strconv.Atoi(vars["seriesID"])
	if err != nil {
		log.Errorf("could not convert seriesID [%s] to int: %v", vars["seriesID"], err)
		h.failure(rw, req, err)
		return
	}
	if seriesID < 1 || seriesID > 99 {
		seriesID = 2
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

	// do we need to update the cached csv file?
	// check if file already exists and is up-to-date, serve it immediately if yes
	if !forceOverwrite && csv.IsAvailable(seriesID, 0) {
		http.ServeFile(rw, req, csv.Filename(seriesID, 0))
		return
	}

	var data bytes.Buffer
	_, _ = data.WriteString("ID;SEASON;WEEK;TRACK;TYPE;LAPS;TIME_OF_DAY;OFFICIAL_RACES;AVG_CAUTIONS;AVG_LAPTIME;FASTEST_LAPTIME;AVG_SOF;HIGHEST_SOF;LOWEST_SOF;NUM_OF_SPLITS;AVG_DRIVERS_PER_SPLIT;UNIQUE_DRIVERS;TOTAL_DRIVERS\n")

	// get all seasons
	seasons, err := h.getSeasons(seriesID)
	if err != nil {
		log.Errorf("could not get seasons: %v", err)
		h.failure(rw, req, err)
		return
	}
	// sort seasons ascending
	sort.Slice(seasons, func(i, j int) bool {
		return seasons[i].StartDate.Before(seasons[j].StartDate)
	})

	// get all 12 weeks for all seasons
	for _, season := range seasons {
		for week := 1; week <= 13; week++ {
			_, track, err := h.getRaceWeek(season.SeasonID, week-1)
			if err != nil {
				log.Debugf("data export: could not get raceweek/track for season[%d], week[%d]: %v", season.SeasonID, week, err)
				continue
			}
			weekResults, err := h.getRaceWeekResults(season.SeasonID, week-1)
			if err != nil {
				log.Errorf("data export: could not get raceweek results for season[%d], week[%d]: %v", season.SeasonID, week-1, err)
				continue
			}
			raceResults, err := h.getRaceResults(season.SeasonID, week-1)
			if err != nil {
				log.Errorf("data export: could not get race results for season[%d], week[%d]: %v", season.SeasonID, week-1, err)
				continue
			}
			weekMetrics, err := h.getRaceWeekMetrics(season.SeasonID, week-1)
			if err != nil {
				log.Errorf("data export: could not get raceweek metrics for season[%d], week[%d]: %v", season.SeasonID, week-1, err)
				continue
			}

			var numOfSplits, uniqueDrivers, totalDrivers, officialRaces int
			splitSubSessionIDs := make(map[int]bool)
			driverIDs := make(map[int]bool)
			for _, result := range weekResults {
				if result.Official {
					officialRaces++
					totalDrivers += result.SizeOfField

					// check if there was a split session
					for _, r2 := range weekResults {
						if r2.SessionID == result.SessionID && r2.SubsessionID != result.SubsessionID {
							splitSubSessionIDs[result.SubsessionID] = true
						}
					}

					// get driver stats
					for _, race := range raceResults {
						if race.SubsessionID == result.SubsessionID {
							driverIDs[race.Driver.DriverID] = true
						}
					}
				}
			}
			numOfSplits = len(splitSubSessionIDs)
			uniqueDrivers = len(driverIDs)

			_, _ = data.WriteString(fmt.Sprintf("%dS%dW%d;%dS%d;%d;%s;%s;%d;%s;%d;%d;%s;%s;%d;%d;%d;%d;%d;%d;%d",
				season.Year, season.Quarter, week, season.Year, season.Quarter, week, track.Name, track.Category,
				weekMetrics.Laps, weekMetrics.TimeOfDay.Format("2006-01-02 15:04"), officialRaces,
				weekMetrics.AvgCautions, weekMetrics.AvgLaptime, weekMetrics.FastestLaptime,
				weekMetrics.AvgSOF, weekMetrics.MaxSOF, weekMetrics.MinSOF, numOfSplits, weekMetrics.AvgSize,
				uniqueDrivers, totalDrivers,
			))
			_, _ = data.WriteString("\n")
		}
	}

	if err := csv.Write(seriesID, 0, data.Bytes()); err != nil {
		log.Errorf("could not write csv file for seriesID [%d]: %v", seriesID, err)
		h.failure(rw, req, err)
		return
	}
	if err := csv.WriteMetadata(seriesID, 0, "", 0, 0); err != nil {
		log.Errorf("could not write metadata file for seriesID [%d]: %v", seriesID, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated csv file
	http.ServeFile(rw, req, csv.Filename(seriesID, 0))
}
