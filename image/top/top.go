package top

import (
	"fmt"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
)

type DataSet struct {
	Title string
	Rows  []DataSetRow
}

type DataSetRow struct {
	Driver string
	Value  string
}

type Top struct {
	Name         string
	Season       database.Season
	Week         database.RaceWeek
	Track        database.Track
	Data         []DataSet
	BorderSize   float64
	FooterHeight float64
	ImageHeight  float64
	ImageWidth   float64
	HeaderHeight float64
	DriverHeight float64
	PaddingSize  float64
	Columns      float64
	ColumnWidth  float64
}

func New(name string, season database.Season, week database.RaceWeek, track database.Track, data []DataSet) Top {
	top := Top{
		Name:         name,
		Season:       season,
		Week:         week,
		Track:        track,
		Data:         data,
		BorderSize:   float64(2),
		FooterHeight: float64(14),
		ImageWidth:   float64(740),
		HeaderHeight: float64(24),
		DriverHeight: float64(18),
		PaddingSize:  float64(3),
		Columns:      float64(len(data)),
	}
	top.ColumnWidth = top.ImageWidth / top.Columns

	maxRows := 0
	for _, d := range data {
		if len(d.Rows) > maxRows {
			maxRows = len(d.Rows)
		}
	}
	top.ImageHeight = float64(maxRows)*top.DriverHeight + top.DriverHeight + top.HeaderHeight + top.PaddingSize*3
	return top
}

func IsAvailable(name string, seasonID, week int) bool {
	return image.IsAvailable("top/"+name, seasonID, week)
}

func Filename(name string, seasonID, week int) string {
	return image.ImageFilename("top/"+name, seasonID, week)
}

func (t *Top) Filename() string {
	return Filename(t.Name, t.Season.SeasonID, t.Week.RaceWeek+1)
}

func (t *Top) Draw() error {
	// top titles, season + track
	topTitle := fmt.Sprintf("%s - Statistics", t.Season.SeasonName)
	topTrackTitle := fmt.Sprintf("Week %d - %s", t.Week.RaceWeek+1, t.Track.Name)
	if t.Week.RaceWeek == -1 { // seasonal avg. top
		topTrackTitle = "Seasonal Average"
	}

	log.Infof("draw top for [%s] - [%s]", topTitle, topTrackTitle)

	// create canvas
	dc := gg.NewContext(int(t.ImageWidth), int(t.ImageHeight))

	// background
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, t.ImageWidth, t.HeaderHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()
	dc.DrawRectangle(t.ImageWidth/2, 0, t.ImageWidth/2, t.HeaderHeight)
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.Fill()

	// draw season name
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(topTitle, t.ImageWidth/4, t.HeaderHeight/2, 0.5, 0.5)
	// draw track title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(topTrackTitle, t.ImageWidth/2+t.ImageWidth/4, t.HeaderHeight/2, 0.5, 0.5)

	// draw the column headers
	yPosColumnHeaderStart := t.HeaderHeight + t.PaddingSize
	xLength := t.ColumnWidth - t.PaddingSize*2
	for column, data := range t.Data {
		xPos := t.PaddingSize + float64(column)*t.ColumnWidth
		yPos := yPosColumnHeaderStart

		// add column header
		dc.DrawRectangle(xPos, yPos, xLength, t.DriverHeight)
		dc.SetRGB255(133, 133, 133) // gray 1
		dc.Fill()

		dc.SetRGB255(255, 255, 255) // white
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(data.Title, xPos+xLength/2, yPos+t.DriverHeight/2, 0.5, 0.5)

		// draw box
		dc.SetRGB255(55, 55, 55) // dark gray 2
		dc.MoveTo(xPos, yPos)
		dc.LineTo(xPos+xLength, yPos)
		dc.LineTo(xPos+xLength, yPos+t.DriverHeight)
		dc.LineTo(xPos, yPos+t.DriverHeight)
		dc.LineTo(xPos, yPos)
		dc.SetLineWidth(1)
		dc.Stroke()
	}

	// draw the columns
	yPosColumnStart := yPosColumnHeaderStart + t.DriverHeight + t.PaddingSize
	for column, data := range t.Data {
		xPos := t.PaddingSize + float64(column)*t.ColumnWidth

		// rows
		var previousValue string
		for row, entry := range data.Rows {
			yPos := yPosColumnStart + float64(row)*t.DriverHeight

			// zebra pattern
			dc.DrawRectangle(xPos, yPos, xLength, t.DriverHeight)
			if row%2 == 0 {
				dc.SetRGB255(225, 225, 225) // light gray 1.5
			} else {
				dc.SetRGB255(241, 241, 241) // light gray 2.5
			}
			dc.Fill()

			dc.SetRGB255(0, 0, 0) // black
			// position
			if err := dc.LoadFontFace("public/fonts/Roboto-Light.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			if entry.Value != previousValue {
				previousValue = entry.Value
				dc.DrawStringAnchored(fmt.Sprintf("%d.", row+1), xPos+t.PaddingSize*2, yPos+t.DriverHeight/2, 0, 0.5)
			}
			// name + value
			if err := dc.LoadFontFace("public/fonts/Roboto-Regular.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			dc.DrawStringAnchored(entry.Driver, xPos+20+t.PaddingSize*2, yPos+t.DriverHeight/2, 0, 0.5)
			dc.SetRGB255(7, 55, 99) // dark blue 3
			if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 12); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			dc.DrawStringAnchored(entry.Value, xPos+xLength-t.PaddingSize*2, yPos+t.DriverHeight/2, 1, 0.5)

			// draw box
			dc.SetRGB255(155, 155, 155) // gray 2
			dc.MoveTo(xPos, yPos)
			dc.LineTo(xPos+xLength, yPos)
			dc.LineTo(xPos+xLength, yPos+t.DriverHeight)
			dc.LineTo(xPos, yPos+t.DriverHeight)
			dc.LineTo(xPos, yPos)
			dc.SetLineWidth(0.5)
			dc.Stroke()
		}
	}

	// add border to image
	bdc := gg.NewContext(int(t.ImageWidth+t.BorderSize*2), int(t.ImageHeight+t.BorderSize*2))
	bdc.SetRGB255(39, 39, 39) // dark gray 1
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(t.BorderSize), int(t.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(t.FooterHeight))
	fdc.SetRGBA255(0, 0, 0, 0) // white
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	fdc.SetRGB255(0, 0, 0) // black
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 10); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	lastUpdate := t.Week.LastUpdate.UTC().Format("2006-01-02 15:04:05 -07 MST")
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", lastUpdate), float64(bdc.Width())-t.FooterHeight/2, float64(bdc.Height())+t.FooterHeight/2, 1, 0.5)

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
	dc.SetRGB255(133, 133, 133) // gray 1
	dc.SetRGB255(155, 155, 155) // gray 2
	dc.SetRGB255(177, 177, 177) // gray 3
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(225, 225, 225) // light gray 1.5
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
