package laptime

import (
	"github.com/JamesClonk/iRvisualizer/image"
)

func (l *Laptime) MetadataFilename() string {
	return image.MetadataFilename("laptimes", l.Season.SeasonID, l.Week.RaceWeek+1)
}

func (l *Laptime) ReadMetadata() (meta image.Metadata) {
	return image.GetMetadata(l.MetadataFilename())
}

func (l *Laptime) WriteMetadata() error {
	// image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time
	return image.WriteMetadata(l.ColorScheme, "laptimes",
		l.Season.SeasonID, l.Week.RaceWeek+1,
		l.Season.SeasonName, l.Season.Year, l.Season.Quarter,
		l.Track.Name, l.Season.StartDate,
	)
}
