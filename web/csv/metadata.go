package csv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/JamesClonk/iRvisualizer/log"
)

type Metadata struct {
	CSVFilename string `json:"CSVFilename"`
	SeriesID    int
	SeasonID    int
	Season      string
	Year        int
	Quarter     int
	LastUpdated time.Time `json:"LastUpdated"`
}

func MetadataFilename(seriesID, seasonID int) string {
	return fmt.Sprintf("%s.json", Filename(seriesID, seasonID))
}

func GetMetadata(filename string) (meta Metadata) {
	log.Debugf("read metadata of [%s]", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("could not read metadata file [%s]: %v", filename, err)
		return meta
	}
	if err = json.Unmarshal(data, &meta); err != nil {
		log.Errorf("could not unmarshal metadata: %v", string(data))
		log.Errorf("%v", err)
		return meta
	}
	return meta
}

func WriteMetadata(seriesID, seasonID int, season string, year, quarter int) error {
	filename := MetadataFilename(seriesID, seasonID)
	log.Debugf("write metadata to [%s]", filename)

	meta := Metadata{
		CSVFilename: Filename(seriesID, seasonID),
		SeriesID:    seriesID,
		SeasonID:    seasonID,
		Season:      season,
		Year:        year,
		Quarter:     quarter,
		LastUpdated: time.Now().UTC(),
	}

	metaJson, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, metaJson, 0644)
}
