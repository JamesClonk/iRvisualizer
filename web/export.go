package web

import (
	"bytes"
	"fmt"
	"math"
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

	_, _ = rw.Write([]byte("SERIES_ID;SERIES_NAME;CURRENT_SEASON;CURRENT_SEASON_ID;CURRENT_WEEK\n"))
	for _, series := range series {
		_, _ = rw.Write([]byte(fmt.Sprintf("%d;%s;%s;%d;%d\n", series.SeriesID, series.SeriesNameShort, series.CurrentSeason, series.CurrentSeasonID, series.CurrentWeek)))
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
	if !forceOverwrite && csv.IsAvailable(seriesID, "weekly") {
		http.ServeFile(rw, req, csv.Filename(seriesID, "weekly"))
		return
	}

	var data bytes.Buffer
	_, _ = data.WriteString("ID;SEASON;WEEK;TRACK;TYPE;LAPS;TIME_OF_DAY;OFFICIAL_RACES;AVG_CAUTIONS;AVG_LAPTIME;FASTEST_LAPTIME;AVG_SOF;HIGHEST_SOF;LOWEST_SOF;NUM_OF_SPLITS;AVG_DRIVERS_PER_SPLIT;UNIQUE_DRIVERS;TOTAL_DRIVERS;AVG_RACES_PER_UNIQUE_DRIVER;STDEV_RACES_PER_DRIVER;STDEV_AVG_RACES_PER_WEEK\n")

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
		racesPerDriver := make([]float64, 0)
		lines := make([]string, 0)

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
			driverIDs := make(map[int]int)
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
							if _, ok := driverIDs[race.Driver.DriverID]; ok {
								driverIDs[race.Driver.DriverID] = driverIDs[race.Driver.DriverID] + 1
							} else {
								driverIDs[race.Driver.DriverID] = 1
							}
						}
					}
				}
			}
			numOfSplits = len(splitSubSessionIDs)
			uniqueDrivers = len(driverIDs)

			// stdev of races per driver this week
			racesPerDriver = append(racesPerDriver, float64(totalDrivers)/float64(uniqueDrivers))
			var stdev float64
			for _, races := range driverIDs {
				stdev += math.Pow(float64(races)-racesPerDriver[week-1], 2)
			}
			stdev = math.Sqrt(stdev / float64(uniqueDrivers))

			// buffer week data, to also add week-stdev below
			lines = append(lines, fmt.Sprintf("%dS%dW%d;%dS%d;%d;%s;%s;%d;%s;%d;%d;%s;%s;%d;%d;%d;%d;%d;%d;%d;%.2f;%.2f",
				season.Year, season.Quarter, week, season.Year, season.Quarter, week, track.Name, track.Category,
				weekMetrics.Laps, weekMetrics.TimeOfDay.Format("2006-01-02 15:04"), officialRaces,
				weekMetrics.AvgCautions, weekMetrics.AvgLaptime, weekMetrics.FastestLaptime,
				weekMetrics.AvgSOF, weekMetrics.MaxSOF, weekMetrics.MinSOF, numOfSplits, weekMetrics.AvgSize,
				uniqueDrivers, totalDrivers, racesPerDriver[week-1], stdev,
			))
		}

		for week, line := range lines {
			// stdev of races per driver compared to all weeks
			var stdev float64
			for _, weeks := range racesPerDriver {
				stdev += math.Pow(weeks-racesPerDriver[week], 2)
			}
			stdev = math.Sqrt(stdev / float64(len(racesPerDriver)))

			// print metrics
			_, _ = data.WriteString(fmt.Sprintf("%s;%.2f", line, stdev))
			_, _ = data.WriteString("\n")
		}
	}

	if err := csv.Write(seriesID, "weekly", data.Bytes()); err != nil {
		log.Errorf("could not write csv file for seriesID [%d]: %v", seriesID, err)
		h.failure(rw, req, err)
		return
	}
	if err := csv.WriteMetadata(seriesID, "weekly"); err != nil {
		log.Errorf("could not write metadata file for seriesID [%d]: %v", seriesID, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated csv file
	http.ServeFile(rw, req, csv.Filename(seriesID, "weekly"))
}

func (h *Handler) seriesSeasonExport(rw http.ResponseWriter, req *http.Request) {
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
	if !forceOverwrite && csv.IsAvailable(seriesID, "season") {
		http.ServeFile(rw, req, csv.Filename(seriesID, "season"))
		return
	}

	var data bytes.Buffer
	_, _ = data.WriteString("SEASON;TIMESLOTS;WEEKS;OFFICIAL_RACES;AVG_DRIVERS_PER_SESSION;AVG_SOF;TOTAL_DRIVERS;UNIQUE_DRIVERS;UNIQUE_ROAD_DRIVERS;UNIQUE_ROAD_ONLY_DRIVERS;UNIQUE_COMMITTED_ROAD_ONLY_DRIVERS;UNIQUE_OVAL_DRIVERS;UNIQUE_OVAL_ONLY_DRIVERS;UNIQUE_COMMITTED_OVAL_ONLY_DRIVERS;UNIQUE_BOTH_DRIVERS;UNIQUE_EIGHT_WEEKS_DRIVERS;UNIQUE_FULL_SEASON_DRIVERS\n")

	// get all season metrics
	metrics, err := h.getSeasonMetrics(seriesID)
	if err != nil {
		log.Errorf("could not get season metrics: %v", err)
		h.failure(rw, req, err)
		return
	}

	// print metrics
	for _, season := range metrics {
		_, _ = data.WriteString(fmt.Sprintf("%dS%d;%s;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d;%d",
			season.Year, season.Quarter, season.Timeslots, season.Weeks, season.Sessions,
			season.AvgSize, season.AvgSOF, season.Drivers,
			season.UniqueDrivers, season.UniqueRoadDrivers, season.UniqueRoadDrivers-season.UniqueBothDrivers, season.UniqueCommittedRoadOnlyDrivers,
			season.UniqueOvalDrivers, season.UniqueOvalDrivers-season.UniqueBothDrivers, season.UniqueCommittedOvalOnlyDrivers,
			season.UniqueBothDrivers, season.UniqueEightWeeksDrivers, season.UniqueFullSeasonDrivers,
		))
		_, _ = data.WriteString("\n")
	}

	if err := csv.Write(seriesID, "season", data.Bytes()); err != nil {
		log.Errorf("could not write season metrics csv file for seriesID [%d]: %v", seriesID, err)
		h.failure(rw, req, err)
		return
	}
	if err := csv.WriteMetadata(seriesID, "season"); err != nil {
		log.Errorf("could not write season metrics metadata file for seriesID [%d]: %v", seriesID, err)
		h.failure(rw, req, err)
		return
	}

	// serve new/updated csv file
	http.ServeFile(rw, req, csv.Filename(seriesID, "season"))
}
