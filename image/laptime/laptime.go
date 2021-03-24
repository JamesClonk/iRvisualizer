package laptime

import (
	"fmt"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	scheme "github.com/JamesClonk/iRvisualizer/image/color"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	laptimeDraws = promauto.NewCounter(prometheus.CounterOpts{
		Name: "irvisualizer_laptimes_drawn_total",
		Help: "Total laptimes drawn by iRvisualizer.",
	})
)

type DataSet struct {
	Division string
	Driver   string
	Laptime  database.Laptime
}

type Laptime struct {
	ColorScheme         string
	Name                string
	Season              database.Season
	Week                database.RaceWeek
	Track               database.Track
	Data                []DataSet
	BorderSize          float64
	FooterHeight        float64
	ImageHeight         float64
	ImageWidth          float64
	HeaderHeight        float64
	ColumnHeaderHeight  float64
	DriverHeight        float64
	PaddingSize         float64
	Rows                float64
	LaptimeColumns      float64
	LaptimeColumnWidth  float64
	DivisionColumnWidth float64
	DriverColumnWidth   float64
}

func New(colorScheme string, season database.Season, week database.RaceWeek, track database.Track, data []DataSet) Laptime {
	lap := Laptime{
		ColorScheme:         colorScheme,
		Name:                "laptimes",
		Season:              season,
		Week:                week,
		Track:               track,
		Data:                data,
		BorderSize:          float64(2),
		FooterHeight:        float64(14),
		ImageWidth:          float64(528),
		HeaderHeight:        float64(46),
		ColumnHeaderHeight:  float64(16),
		DriverHeight:        float64(24),
		PaddingSize:         float64(3),
		Rows:                float64(len(data)),
		LaptimeColumns:      float64(5),
		LaptimeColumnWidth:  float64(56),
		DivisionColumnWidth: float64(64),
	}
	lap.DriverColumnWidth = lap.ImageWidth - (lap.DivisionColumnWidth + (lap.LaptimeColumnWidth * lap.LaptimeColumns))
	lap.ImageHeight = lap.Rows*lap.DriverHeight + lap.ColumnHeaderHeight + lap.HeaderHeight + lap.PaddingSize*3
	return lap
}

func IsAvailable(colorScheme string, seasonID, week int) bool {
	return image.IsAvailable(colorScheme, "laptimes", seasonID, week)
}

func Filename(seasonID, week int) string {
	return image.ImageFilename("laptimes", seasonID, week)
}

func (l *Laptime) Filename() string {
	return Filename(l.Season.SeasonID, l.Week.RaceWeek+1)
}

func (l *Laptime) Draw() error {
	laptimeDraws.Inc()

	// laptime titles, season + track
	lapTitle := fmt.Sprintf("%s - Fastest Laptimes", l.Season.SeasonName)
	if len(l.Season.SeasonName) > 64 {
		lapTitle = l.Season.SeasonName
	}
	lapWeekTitle := fmt.Sprintf("Week %d", l.Week.RaceWeek+1)
	lapTrackTitle := l.Track.Name

	log.Infof("draw laptimes for [%s] - [%s]", lapTitle, lapTrackTitle)

	// colorizer
	color := scheme.Get(l.ColorScheme)

	// create canvas
	dc := gg.NewContext(int(l.ImageWidth), int(l.ImageHeight))

	// background
	color.Background(dc)
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, l.ImageWidth, l.HeaderHeight/2)
	color.HeaderLeftBG(dc)
	dc.Fill()
	dc.DrawRectangle(0, l.HeaderHeight/2, l.ImageWidth, l.HeaderHeight/2)
	color.HeaderRightBG(dc)
	dc.Fill()

	// draw season title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(lapTitle, l.PaddingSize*3, l.HeaderHeight/4, 0, 0.5)
	// draw week title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(lapWeekTitle, l.ImageWidth/4, l.HeaderHeight/4*3, 0.5, 0.5)
	// draw track title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(lapTrackTitle, l.ImageWidth/3*2, l.HeaderHeight/4*3, 0.5, 0.5)

	// adjust to header height
	yPosColumnHeaderStart := l.HeaderHeight + l.PaddingSize

	// draw division column header
	xDivisionLength := l.DivisionColumnWidth - l.PaddingSize*2
	xPos := l.PaddingSize
	yPos := yPosColumnHeaderStart

	dc.DrawRectangle(xPos, yPos, xDivisionLength, l.ColumnHeaderHeight)
	color.TopNHeaderBG(dc)
	dc.Fill()

	color.TopNHeaderFG(dc)
	if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.DrawStringAnchored("Division", xPos+xDivisionLength/2, yPos+l.ColumnHeaderHeight/2, 0.5, 0.5)

	// draw outline
	color.TopNHeaderOutline(dc)
	dc.MoveTo(xPos, yPos)
	dc.LineTo(xPos+xDivisionLength, yPos)
	dc.LineTo(xPos+xDivisionLength, yPos+l.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos+l.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos)
	dc.SetLineWidth(1)
	dc.Stroke()

	// draw driver column header
	xDriverLength := l.DriverColumnWidth - l.PaddingSize*2
	xPos = xDivisionLength + l.PaddingSize*2

	dc.DrawRectangle(xPos, yPos, xDriverLength, l.ColumnHeaderHeight)
	color.TopNHeaderBG(dc)
	dc.Fill()

	color.TopNHeaderFG(dc)
	if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.DrawStringAnchored("Driver", xPos+xDriverLength/2, yPos+l.ColumnHeaderHeight/2, 0.5, 0.5)

	// draw outline
	color.TopNHeaderOutline(dc)
	dc.MoveTo(xPos, yPos)
	dc.LineTo(xPos+xDriverLength, yPos)
	dc.LineTo(xPos+xDriverLength, yPos+l.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos+l.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos)
	dc.SetLineWidth(1)
	dc.Stroke()

	// draw laptime column headers
	xColumnLength := l.LaptimeColumnWidth - l.PaddingSize
	for column := float64(0); column < l.LaptimeColumns; column++ {
		xPos := xDivisionLength + l.PaddingSize*2 + xDriverLength + l.PaddingSize + float64(column)*l.LaptimeColumnWidth
		yPos := yPosColumnHeaderStart

		title := fmt.Sprintf("%d%%", 100+int(column)+int(column-1))
		xLength := xColumnLength
		if column == 0 {
			xLength = xLength + l.PaddingSize
			title = "100%"
		} else {
			xPos = xPos + l.PaddingSize
		}

		dc.DrawRectangle(xPos, yPos, xLength, l.ColumnHeaderHeight)
		color.TopNHeaderBG(dc)
		dc.Fill()

		color.TopNHeaderFG(dc)
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(title, xPos+xLength/2, yPos+l.ColumnHeaderHeight/2, 0.5, 0.5)

		// draw outline
		color.TopNHeaderOutline(dc)
		dc.MoveTo(xPos, yPos)
		dc.LineTo(xPos+xLength, yPos)
		dc.LineTo(xPos+xLength, yPos+l.ColumnHeaderHeight)
		dc.LineTo(xPos, yPos+l.ColumnHeaderHeight)
		dc.LineTo(xPos, yPos)
		dc.SetLineWidth(1)
		dc.Stroke()
	}

	// draw rows
	yPosRowStart := yPosColumnHeaderStart + l.ColumnHeaderHeight + l.PaddingSize
	for row, entry := range l.Data {
		// draw division
		xPos := l.PaddingSize
		yPos := yPosRowStart + float64(row)*l.DriverHeight

		// zebra pattern
		dc.DrawRectangle(xPos, yPos, l.ImageWidth-l.PaddingSize*2, l.DriverHeight)
		if row%2 == 0 {
			color.TopNCellDarkerBG(dc)
		} else {
			color.TopNCellLighterBG(dc)
		}
		dc.Fill()
		log.Debugf("%#v", entry)
	}

	//------------------------------------------------------------------------------------------------------------------

	// // draw the columns
	// yPosColumnStart := yPosColumnHeaderStart + l.DriverHeight + l.PaddingSize
	// for column, data := range l.Data {
	// 	xPos := l.PaddingSize + float64(column)*l.ColumnWidth

	// 	// rows
	// 	var previousValue string
	// 	for row, entry := range data.Rows {
	// 		yPos := yPosColumnStart + float64(row)*l.DriverHeight

	// 		// zebra pattern
	// 		dc.DrawRectangle(xPos, yPos, xLength, l.DriverHeight)
	// 		if row%2 == 0 {
	// 			color.TopNCellDarkerBG(dc)
	// 		} else {
	// 			color.TopNCellLighterBG(dc)
	// 		}
	// 		dc.Fill()

	// 		// position
	// 		color.TopNCellPosition(dc)
	// 		if err := dc.LoadFontFace("public/fonts/Roboto-Light.ttf", 11); err != nil {
	// 			return fmt.Errorf("could not load font: %v", err)
	// 		}
	// 		if entry.Value != previousValue {
	// 			previousValue = entry.Value

	// 			// draw icons if specified
	// 			if len(data.Icons) > 0 && row <= 2 {
	// 				// load icon
	// 				iconColor := "gold"
	// 				if row == 1 {
	// 					iconColor = "silver"
	// 				}
	// 				if row == 2 {
	// 					iconColor = "bronze"
	// 				}
	// 				icon, err := gg.LoadPNG(fmt.Sprintf("public/icons/%s_%s.png", data.Icons, iconColor))
	// 				if err != nil {
	// 					return fmt.Errorf("could not load icon: %v", err)
	// 				}
	// 				dc.DrawImage(icon, int(xPos+l.PaddingSize), int(yPos))
	// 			} else {
	// 				dc.DrawStringAnchored(fmt.Sprintf("%d.", row+1), xPos+l.PaddingSize*2, yPos+l.DriverHeight/2, 0, 0.5)
	// 			}
	// 		}
	// 		// name
	// 		color.TopNCellDriver(dc)
	// 		if err := dc.LoadFontFace("public/fonts/Roboto-Regular.ttf", 11); err != nil {
	// 			return fmt.Errorf("could not load font: %v", err)
	// 		}
	// 		dc.DrawStringAnchored(entry.Driver, xPos+20+l.PaddingSize*2, yPos+l.DriverHeight/2, 0, 0.5)
	// 		// value
	// 		color.TopNCellValue(dc)
	// 		if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 12); err != nil {
	// 			return fmt.Errorf("could not load font: %v", err)
	// 		}
	// 		dc.DrawStringAnchored(entry.Value, xPos+xLength-l.PaddingSize*2, yPos+l.DriverHeight/2, 1, 0.5)
	// 		// draw an icon if specified
	// 		if len(entry.Icon) > 0 {
	// 			icon, err := gg.LoadPNG(fmt.Sprintf("public/icons/%s.png", entry.Icon))
	// 			if err != nil {
	// 				return fmt.Errorf("could not load icon: %v", err)
	// 			}
	// 			dc.DrawImageAnchored(icon, int(xPos+xLength-l.PaddingSize*2)-entry.IconPosition, int(yPos), 1, 0)
	// 		}

	// 		// draw outline
	// 		color.TopNCellOutline(dc)
	// 		dc.MoveTo(xPos, yPos)
	// 		dc.LineTo(xPos+xLength, yPos)
	// 		dc.LineTo(xPos+xLength, yPos+l.DriverHeight)
	// 		dc.LineTo(xPos, yPos+l.DriverHeight)
	// 		dc.LineTo(xPos, yPos)
	// 		dc.SetLineWidth(0.5)
	// 		dc.Stroke()
	// 	}
	// }

	// add border to image
	bdc := gg.NewContext(int(l.ImageWidth+l.BorderSize*2), int(l.ImageHeight+l.BorderSize*2))
	color.Border(bdc)
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(l.BorderSize), int(l.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(l.FooterHeight))
	color.Transparent(fdc)
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	color.LastUpdate(fdc)
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 10); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	lastUpdate := l.Week.LastUpdate.UTC().Format("2006-01-02 15:04:05 -07 MST")
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", lastUpdate), float64(bdc.Width())-l.FooterHeight/2, float64(bdc.Height())+l.FooterHeight/2, 1, 0.5)

	color.CreatedBy(fdc)
	if err := fdc.LoadFontFace("public/fonts/Roboto-Light.ttf", 9); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	fdc.DrawStringAnchored("created by Fabio Berchtold", l.FooterHeight/2, float64(bdc.Height())+l.FooterHeight/2, 0, 0.5)

	if err := l.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(l.Filename()) // finally write to file
}
