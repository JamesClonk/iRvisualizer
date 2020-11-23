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
	Mode        string
	LastUpdated time.Time `json:"LastUpdated"`
}

func MetadataFilename(seriesID int, mode string) string {
	return fmt.Sprintf("%s.json", Filename(seriesID, mode))
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

func WriteMetadata(seriesID int, mode string) error {
	filename := MetadataFilename(seriesID, mode)
	log.Debugf("write metadata to [%s]", filename)

	meta := Metadata{
		CSVFilename: Filename(seriesID, mode),
		SeriesID:    seriesID,
		Mode:        mode,
		LastUpdated: time.Now().UTC(),
	}

	metaJson, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, metaJson, 0644)
}
