package heatmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/JamesClonk/iRcollector/database"
)

type Metadata struct {
	ImageFilename string
	Season        string
	Year          int
	Quarter       int
	Track         string
	StartDate     time.Time
	LastUpdated   time.Time
}

func MetadataFilename(season database.Season, week database.RaceWeek) {
	return fmt.Sprintf("public/heatmaps/season_%d_week_%d.png.json", season.SeasonID, week.RaceWeek+1)
}

func ReadMetadata(season database.Season, week database.RaceWeek) (meta Metadata) {
	filename := MetadataFilename(season, week)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return meta
	}
	if err = json.Unmarshal(data, meta); err != nil {
		return meta
	}
	return meta
}

func WriteMetadata(season database.Season, week database.RaceWeek, track database.Track) error {
	meta := Metadata{
		ImageFilename: Filename(season, week),
		Season:        season.SeasonName,
		Year:          season.Year,
		Quarter:       season.Quarter,
		Track:         track.Name,
		StartDate:     season.StartDate,
		LastUpdated:   time.Now(),
	}

	metaJson, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(MetadataFilename(season, week), metaJson, 0644)
}
