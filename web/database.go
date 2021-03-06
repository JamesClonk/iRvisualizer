package web

import (
	"sort"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/log"
)

func (h *Handler) getSeries() ([]database.Series, error) {
	log.Infof("collect active series")
	return h.DB.GetActiveSeries()
}

func (h *Handler) getSeasons(seriesID int) ([]database.Season, error) {
	log.Infof("collect seasons by series ID [%d]", seriesID)
	return h.DB.GetSeasonsBySeriesID(seriesID)
}

func (h *Handler) getSeason(seasonID int) (database.Season, error) {
	log.Infof("collect season [%d]", seasonID)
	return h.DB.GetSeasonByID(seasonID)
}

func (h *Handler) getSeasonMetrics(seriesID int) ([]database.SeasonMetrics, error) {
	log.Infof("collect season metrics for series [%d]", seriesID)
	return h.DB.GetSeasonMetricsBySeriesID(seriesID)
}

func (h *Handler) getRaceWeekMetrics(seasonID, week int) (database.RaceWeekMetrics, error) {
	log.Infof("collect raceweek metrics for season [%d], week [%d]", seasonID, week)
	return h.DB.GetRaceWeekMetricsBySeasonIDAndWeek(seasonID, week)
}

func (h *Handler) getRaceWeek(seasonID, week int) (database.RaceWeek, database.Track, error) {
	log.Infof("collect raceweek for season [%d], week [%d]", seasonID, week)

	raceweek, err := h.DB.GetRaceWeekBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return database.RaceWeek{}, database.Track{}, err
	}
	track, err := h.DB.GetTrackByID(raceweek.TrackID)
	if err != nil {
		return raceweek, database.Track{}, err
	}
	return raceweek, track, nil
}

func (h *Handler) getRaceWeekResults(seasonID, week int) ([]database.RaceWeekResult, error) {
	log.Infof("collect raceweek results for season [%d], week [%d]", seasonID, week)

	results, err := h.DB.GetRaceWeekResultsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return nil, err
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.Before(results[j].StartTime)
	})
	return results, nil
}

func (h *Handler) getRaceResults(seasonID, week int) ([]database.RaceResult, error) {
	log.Infof("collect race results for season [%d], week [%d]", seasonID, week)

	results, err := h.DB.GetRaceResultsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return nil, err
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].SessionStartTime < results[j].SessionStartTime
	})
	return results, nil
}

func (h *Handler) getRaceWeekSummaries(seasonID, week int) ([]database.Summary, error) {
	log.Infof("collect raceweek summaries for season [%d], week [%d]", seasonID, week)

	summaries, err := h.DB.GetDriverSummariesBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}

func (h *Handler) getRaceWeekTimeRankings(seasonID, week int) ([]database.TimeRanking, error) {
	log.Infof("collect raceweek timerankings for season [%d], week [%d]", seasonID, week)

	timeRankings, err := h.DB.GetTimeRankingsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return nil, err
	}
	return timeRankings, nil
}

func (h *Handler) getChampPoints(seasonID, week int) ([]database.Points, error) {
	log.Infof("collect championship points for season [%d], week [%d]", seasonID, week)

	points, err := h.DB.GetPointsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (h *Handler) getChampPointsForOvals(seasonID, week int) ([]database.Points, error) {
	log.Infof("collect championship points for oval tracks for season [%d], week [%d]", seasonID, week)

	points, err := h.DB.GetPointsBySeasonIDAndWeekAndTrackCategory(seasonID, week, "Oval")
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (h *Handler) getTTStandings(seasonID, week int) ([]database.TimeTrialResult, error) {
	log.Infof("collect time trial results for season [%d], week [%d]", seasonID, week)

	results, err := h.DB.GetTimeTrialResultsBySeasonIDAndWeek(seasonID, week)
	if err != nil {
		return nil, err
	}
	return results, nil
}
