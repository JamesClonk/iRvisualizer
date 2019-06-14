package web

import (
	"sort"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/log"
)

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
