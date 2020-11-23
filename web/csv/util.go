package csv

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
)

func IsAvailable(seriesID int, mode string) bool {
	// check if file already exists
	csvFilename := Filename(seriesID, mode)
	metaFilename := MetadataFilename(seriesID, mode)
	if util.FileExists(metaFilename) && util.FileExists(csvFilename) {
		metadata := GetMetadata(metaFilename)
		// if it's newer than 24 hours
		if time.Since(metadata.LastUpdated) < time.Hour*24 {
			return true
		}
	}
	return false
}

func Filename(seriesID int, mode string) string {
	if len(mode) == 0 || mode == "weekly" {
		return fmt.Sprintf("public/csv/weekly_%d.csv", seriesID)
	}
	return fmt.Sprintf("public/csv/seasons_%d.csv", seriesID)
}

func Write(seriesID int, mode string, data []byte) error {
	filename := Filename(seriesID, mode)
	log.Debugf("write csv to [%s]", filename)
	return ioutil.WriteFile(filename, data, 0644)
}
