package top

import (
	"github.com/JamesClonk/iRvisualizer/image"
)

func (t *Top) MetadataFilename() string {
	return image.MetadataFilename("top/"+t.Name, t.Season.SeasonID, t.Week.RaceWeek+1, t.Team)
}

func (t *Top) ReadMetadata() (meta image.Metadata) {
	return image.GetMetadata(t.MetadataFilename())
}

func (t *Top) WriteMetadata() error {
	// image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time
	return image.WriteMetadata(t.ColorScheme, "top/"+t.Name,
		t.Season.SeasonID, t.Week.RaceWeek+1,
		t.Season.SeasonName, t.Season.Year, t.Season.Quarter,
		t.Track.Name, t.Team, t.Season.StartDate,
	)
}
