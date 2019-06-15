package top20

import (
	"github.com/JamesClonk/iRvisualizer/image"
)

func (t *Top20) MetadataFilename() string {
	return image.MetadataFilename("top20/"+t.Name, t.Season.SeasonID, t.Week.RaceWeek+1)
}

func (t *Top20) ReadMetadata() (meta image.Metadata) {
	return image.GetMetadata(t.MetadataFilename())
}

func (t *Top20) WriteMetadata() error {
	// image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time
	return image.WriteMetadata("top20/"+t.Name,
		t.Season.SeasonID, t.Week.RaceWeek+1,
		t.Season.SeasonName, t.Season.Year, t.Season.Quarter,
		t.Track.Name, t.Season.StartDate,
	)
}
