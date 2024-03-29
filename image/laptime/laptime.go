package laptime

import (
	"fmt"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	scheme "github.com/JamesClonk/iRvisualizer/image/color"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/JamesClonk/iRvisualizer/util"
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
	Marked   bool
}

type Laptime struct {
	ColorScheme         string
	Team                string
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

func New(colorScheme, team string, season database.Season, week database.RaceWeek, track database.Track, data []DataSet) Laptime {
	lap := Laptime{
		ColorScheme:         colorScheme,
		Team:                team,
		Name:                "laptimes",
		Season:              season,
		Week:                week,
		Track:               track,
		Data:                data,
		BorderSize:          float64(2),
		FooterHeight:        float64(14),
		ImageWidth:          float64(756),
		HeaderHeight:        float64(46),
		ColumnHeaderHeight:  float64(16),
		DriverHeight:        float64(24),
		PaddingSize:         float64(3),
		Rows:                float64(len(data)),
		LaptimeColumns:      float64(9),
		LaptimeColumnWidth:  float64(56),
		DivisionColumnWidth: float64(58),
	}
	lap.DriverColumnWidth = lap.ImageWidth - (lap.DivisionColumnWidth + (lap.LaptimeColumnWidth * lap.LaptimeColumns))
	lap.ImageHeight = lap.Rows*lap.DriverHeight + lap.ColumnHeaderHeight + lap.HeaderHeight + lap.PaddingSize*3
	return lap
}

func IsAvailable(colorScheme string, seasonID, week int, team string) bool {
	return image.IsAvailable(colorScheme, "laptimes", seasonID, week, team)
}

func Filename(seasonID, week int, team string) string {
	return image.ImageFilename("laptimes", seasonID, week, team)
}

func (l *Laptime) Filename() string {
	return Filename(l.Season.SeasonID, l.Week.RaceWeek+1, l.Team)
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
	if len(l.ColorScheme) == 0 {
		l.ColorScheme = l.Season.SeriesColorScheme // get series default if needed
	}
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

		title := fmt.Sprintf("%d%%", 100+int(column))
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
		if column == l.LaptimeColumns-1 {
			color.TopNHeaderFGDanger(dc)
		}
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
		xPos := l.PaddingSize
		yPos := yPosRowStart + float64(row)*l.DriverHeight
		xLength := l.ImageWidth - l.PaddingSize*2

		// zebra pattern
		dc.DrawRectangle(xPos, yPos, xLength, l.DriverHeight)
		if row%2 == 0 {
			color.TopNCellDarkerBG(dc)
		} else {
			color.TopNCellLighterBG(dc)
		}
		// marked driver?
		if entry.Marked {
			color.TopNHeaderBG(dc)
		}
		dc.Fill()

		// draw outline
		color.TopNCellOutline(dc)
		dc.MoveTo(xPos, yPos)
		dc.LineTo(xPos+xLength, yPos)
		dc.LineTo(xPos+xLength, yPos+l.DriverHeight)
		dc.LineTo(xPos, yPos+l.DriverHeight)
		dc.LineTo(xPos, yPos)
		dc.SetLineWidth(0.5)
		dc.Stroke()

		// draw division
		color.TopNCellPosition(dc)
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(entry.Division, xPos+xDivisionLength/2, yPos+l.DriverHeight/2, 0.5, 0.5)

		// draw driver
		xPos = xDivisionLength + l.PaddingSize*2
		color.TopNCellDriver(dc)
		// marked driver?
		if entry.Marked {
			color.TopNHeaderFG(dc)
		}
		if err := dc.LoadFontFace("public/fonts/Roboto-Regular.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(entry.Driver, xPos, yPos+l.DriverHeight/2, 0, 0.5)

		// draw laptimes
		xColumnLength := l.LaptimeColumnWidth - l.PaddingSize
		for column := float64(0); column < l.LaptimeColumns; column++ {
			xPos := xDivisionLength + l.PaddingSize*2 + xDriverLength + l.PaddingSize + float64(column)*l.LaptimeColumnWidth

			color.TopNCellValue(dc)
			if column == l.LaptimeColumns-1 {
				color.TopNCellValueDanger(dc)
			}
			if err := dc.LoadFontFace("public/fonts/Roboto-Light.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}

			xLength := xColumnLength
			if column == 0 {
				xLength = xLength + l.PaddingSize
				color.TopNCellValue(dc)
				if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 11); err != nil {
					return fmt.Errorf("could not load font: %v", err)
				}
			} else {
				xPos = xPos + l.PaddingSize
			}

			laptime := util.ConvertLaptime(entry.Laptime)
			if column > 0 {
				percentage := 100 + int(column)
				laptime = util.ConvertLaptime(database.Laptime(float64(entry.Laptime) / float64(100) * float64(percentage)))
			}

			// marked driver?
			if entry.Marked {
				color.TopNHeaderFGDanger(dc)
			}
			dc.DrawStringAnchored(laptime, xPos+xLength/2, yPos+l.DriverHeight/2, 0.5, 0.5)
		}
	}

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
	fdc.DrawStringAnchored("by Fabio Berchtold", l.FooterHeight/2, float64(bdc.Height())+l.FooterHeight/2, 0, 0.5)

	if err := l.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(l.Filename()) // finally write to file
}
