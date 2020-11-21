package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/gorilla/mux"
)

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

	_, _ = rw.Write([]byte("SEASON;WEEK;TRACK;TYPE;LAPS;OFFICIAL_RACES;AVG_LAPTIME;FASTEST_LAPTIME;AVG_SOF;HIGHEST_SOF;NUM_OF_SPLITS;AVG_DRIVERS_PER_SPLIT;UNIQUE_DRIVERS;TOTAL_DRIVERS\n"))

	// get all seasons
	seasons, err := h.getSeasons(seriesID)
	if err != nil {
		log.Errorf("could not get seasons: %v", err)
		h.failure(rw, req, err)
		return
	}

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

			var numOfSplits, averageDrivers, uniqueDrivers, totalDrivers, averageSOF, highestSOF, officialRaces int
			splitSubSessionIDs := make(map[int]bool)
			driverIDs := make(map[int]bool)
			for _, result := range weekResults {
				if result.Official {
					officialRaces++
					averageDrivers += result.SizeOfField
					totalDrivers += result.SizeOfField
					averageSOF += result.StrengthOfField
					if result.StrengthOfField > highestSOF {
						highestSOF = result.StrengthOfField
					}

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
			averageDrivers = int(averageDrivers / officialRaces)
			uniqueDrivers = len(driverIDs)
			averageSOF = int(averageSOF / officialRaces)

			_, _ = rw.Write([]byte(fmt.Sprintf("%dS%d;%d;%s;%s;%d;%d;%s;%s;%d;%d;%d;%d;%d;%d",
				season.Year, season.Quarter, week, track.Name, track.Category,
				0, officialRaces, "avg-laptime", "fastest-laptime",
				averageSOF, highestSOF, numOfSplits, averageDrivers,
				uniqueDrivers, totalDrivers,
			)))
			_, _ = rw.Write([]byte("\n"))
		}
	}
}
