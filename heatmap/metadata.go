package heatmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/JamesClonk/iRvisualizer/log"
)

type Metadata struct {
	ImageFilename string `json:"ImageFilename"`
	Season        string
	Year          int
	Quarter       int
	Week          int
	Track         string
	StartDate     time.Time `json:"StartDate"`
	LastUpdated   time.Time `json:"LastUpdated"`
}

func MetadataFilename(seasonID, week int) string {
	return fmt.Sprintf("%s.json", HeatmapFilename(seasonID, week))
}

func (h *Heatmap) MetadataFilename() string {
	return MetadataFilename(h.Season.SeasonID, h.Week.RaceWeek+1)
}

func GetMetadata(filename string) (meta Metadata) {
	log.Debugf("read metadata of [%s]", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("could not read metadata file [%s]: %v", filename)
		return meta
	}
	if err = json.Unmarshal(data, &meta); err != nil {
		log.Errorf("could not unmarshal metadata: %v", string(data))
		log.Errorf("%v", err)
		return meta
	}
	return meta
}

func (h *Heatmap) ReadMetadata() (meta Metadata) {
	return GetMetadata(h.MetadataFilename())
}

func (h *Heatmap) WriteMetadata() error {
	log.Debugf("write metadata to [%s]", h.MetadataFilename())

	meta := Metadata{
		ImageFilename: h.Filename(),
		Season:        h.Season.SeasonName,
		Year:          h.Season.Year,
		Quarter:       h.Season.Quarter,
		Week:          h.Week.RaceWeek,
		Track:         h.Track.Name,
		StartDate:     h.Season.StartDate,
		LastUpdated:   time.Now().UTC(),
	}

	metaJson, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(h.MetadataFilename(), metaJson, 0644)
}
