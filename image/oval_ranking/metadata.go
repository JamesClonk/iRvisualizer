package oval_ranking

import (
	"github.com/JamesClonk/iRvisualizer/image"
)

func (r *Ranking) MetadataFilename() string {
	return image.MetadataFilename("oval_ranking", r.Season.SeasonID, -1)
}

func (r *Ranking) ReadMetadata() (meta image.Metadata) {
	return image.GetMetadata(r.MetadataFilename())
}

func (r *Ranking) WriteMetadata() error {
	// image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time
	return image.WriteMetadata(r.ColorScheme, "oval_ranking",
		r.Season.SeasonID, -1,
		r.Season.SeasonName, r.Season.Year, r.Season.Quarter,
		"oval_ranking", r.Season.StartDate,
	)
}
