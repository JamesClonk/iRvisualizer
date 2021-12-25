package summary

import (
	"github.com/JamesClonk/iRvisualizer/image"
)

func (s *Summary) MetadataFilename() string {
	return image.MetadataFilename("summary", s.Season.SeasonID, s.Week.RaceWeek+1, s.Team)
}

func (s *Summary) ReadMetadata() (meta image.Metadata) {
	return image.GetMetadata(s.MetadataFilename())
}

func (s *Summary) WriteMetadata() error {
	// image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time
	return image.WriteMetadata(s.ColorScheme, "summary",
		s.Season.SeasonID, s.Week.RaceWeek+1,
		s.Season.SeasonName, s.Season.Year, s.Season.Quarter,
		s.Track.Name, s.Team, s.Season.StartDate,
	)
}
