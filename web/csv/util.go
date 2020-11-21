package csv

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
)

func IsAvailable(seriesID, seasonID int) bool {
	// check if file already exists
	csvFilename := Filename(seriesID, seasonID)
	metaFilename := MetadataFilename(seriesID, seasonID)
	if util.FileExists(metaFilename) && util.FileExists(csvFilename) {
		metadata := GetMetadata(metaFilename)
		// if it's newer than 24 hours
		if time.Since(metadata.LastUpdated) < time.Hour*24 {
			return true
		}
	}
	return false
}

func Filename(seriesID, seasonID int) string {
	if seasonID > 0 {
		return fmt.Sprintf("public/csv/season_%d.csv", seasonID)
	}
	return fmt.Sprintf("public/csv/series_%d.csv", seriesID)
}

func Write(seriesID, seasonID int, data []byte) error {
	filename := Filename(seriesID, seasonID)
	log.Debugf("write csv to [%s]", filename)
	return ioutil.WriteFile(filename, data, 0644)
}
