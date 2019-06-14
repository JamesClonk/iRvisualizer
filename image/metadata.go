package image

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

func MetadataFilename(image string, seasonID, week int) string {
	return fmt.Sprintf("%s.json", ImageFilename(image, seasonID, week))
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

func WriteMetadata(image string, seasonID, week int, season string, year, quarter int, track string, startDate time.Time) error {
	filename := MetadataFilename(image, seasonID, week)
	log.Debugf("write metadata to [%s]", filename)

	meta := Metadata{
		ImageFilename: ImageFilename(image, seasonID, week),
		Season:        season,
		Year:          year,
		Quarter:       quarter,
		Week:          week,
		Track:         track,
		StartDate:     startDate,
		LastUpdated:   time.Now().UTC(),
	}

	metaJson, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, metaJson, 0644)
}
