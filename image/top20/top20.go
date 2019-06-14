package top20

import (
	"fmt"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
)

type Top20 struct {
	Season       database.Season
	Week         database.RaceWeek
	Track        database.Track
	Results      []database.RaceWeekResult
	BorderSize   float64
	ImageHeight  float64
	ImageWidth   float64
	HeaderHeight float64
	DriverHeight float64
	ColumnWidth  float64
}

func New(season database.Season, week database.RaceWeek, track database.Track, results []database.RaceWeekResult) Top20 {
	return Top20{
		Season:       season,
		Week:         week,
		Track:        track,
		Results:      results,
		BorderSize:   float64(2),
		ImageHeight:  float64(480),
		ImageWidth:   float64(1024),
		HeaderHeight: float64(30),
		DriverHeight: float64(26),
		ColumnWidth:  float64(300),
	}
}

func IsAvailable(seasonID, week int) bool {
	return image.IsAvailable("top20", seasonID, week)
}

func Filename(seasonID, week int) string {
	return image.ImageFilename("top20", seasonID, week)
}

func (t *Top20) Filename() string {
	return Filename(t.Season.SeasonID, t.Week.RaceWeek+1)
}

func (t *Top20) Draw() error {
	// top20 titles, season + track
	top20Title := fmt.Sprintf("%s - Statistics", t.Season.SeasonName)
	top20TrackTitle := fmt.Sprintf("Week %d - %s", t.Week.RaceWeek+1, t.Track.Name)
	if t.Week.RaceWeek == -1 { // seasonal avg. top20
		top20TrackTitle = "Seasonal Average"
	}

	log.Infof("draw top20 for [%s] - [%s]", top20Title, top20TrackTitle)

	// create canvas
	dc := gg.NewContext(int(t.ImageWidth), int(t.ImageHeight))

	// background
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, t.ImageWidth, t.HeaderHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()
	dc.DrawRectangle(t.ImageWidth/2, 0, t.ImageWidth/2, t.HeaderHeight)
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.Fill()

	// draw season name
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(top20Title, t.ImageWidth/4, t.HeaderHeight/2, 0.5, 0.5)
	// draw track title
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(top20TrackTitle, t.ImageWidth/2+t.ImageWidth/4, t.HeaderHeight/2, 0.5, 0.5)

	// add border to image
	fdc := gg.NewContext(int(t.ImageWidth+t.BorderSize*2), int(t.ImageHeight+t.BorderSize*2))
	fdc.SetRGB255(39, 39, 39) // dark gray 1
	fdc.Clear()
	fdc.DrawImage(dc.Image(), int(t.BorderSize), int(t.BorderSize))

	if err := t.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(t.Filename()) // finally write to file
}

/*
	Colors:
	dc.SetRGB255(0, 0, 0) // black
	dc.SetRGB255(39, 39, 39) // dark gray 1
	dc.SetRGB255(55, 55, 55) // dark gray 2
	dc.SetRGB255(255, 255, 255) // white
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
