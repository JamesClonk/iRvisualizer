package heatmap

import (
	"github.com/JamesClonk/iRvisualizer/image"
)

func (h *Heatmap) MetadataFilename() string {
	return image.MetadataFilename("heatmap", h.Season.SeasonID, h.Week.RaceWeek+1)
}

func (h *Heatmap) ReadMetadata() (meta image.Metadata) {
	return image.GetMetadata(h.MetadataFilename())
}

func (h *Heatmap) WriteMetadata() error {
	// image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time
	return image.WriteMetadata(h.ColorScheme, "heatmap",
		h.Season.SeasonID, h.Week.RaceWeek+1,
		h.Season.SeasonName, h.Season.Year, h.Season.Quarter,
		h.Track.Name, h.Season.StartDate,
	)
}
