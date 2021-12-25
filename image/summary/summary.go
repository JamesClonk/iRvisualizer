package summary

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
	summaryDraws = promauto.NewCounter(prometheus.CounterOpts{
		Name: "irvisualizer_summaries_drawn_total",
		Help: "Total driver summaries drawn by iRvisualizer.",
	})
)

type DataSet struct {
	Division string
	Driver   string
	Laptime  database.Laptime
	Marked   bool
}

type Summary struct {
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

func New(colorScheme, team string, season database.Season, week database.RaceWeek, track database.Track, data []DataSet) Summary {
	lap := Summary{
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

func (s *Summary) Filename() string {
	return Filename(s.Season.SeasonID, s.Week.RaceWeek+1, s.Team)
}

func (s *Summary) Draw() error {
	summaryDraws.Inc()

	// laptime titles, season + track
	lapTitle := fmt.Sprintf("%s - Fastest Laptimes", s.Season.SeasonName)
	if len(s.Season.SeasonName) > 64 {
		lapTitle = s.Season.SeasonName
	}
	lapWeekTitle := fmt.Sprintf("Week %d", s.Week.RaceWeek+1)
	lapTrackTitle := s.Track.Name

	log.Infof("draw laptimes for [%s] - [%s]", lapTitle, lapTrackTitle)

	// colorizer
	if len(s.ColorScheme) == 0 {
		s.ColorScheme = s.Season.SeriesColorScheme // get series default if needed
	}
	color := scheme.Get(s.ColorScheme)

	// create canvas
	dc := gg.NewContext(int(s.ImageWidth), int(s.ImageHeight))

	// background
	color.Background(dc)
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, s.ImageWidth, s.HeaderHeight/2)
	color.HeaderLeftBG(dc)
	dc.Fill()
	dc.DrawRectangle(0, s.HeaderHeight/2, s.ImageWidth, s.HeaderHeight/2)
	color.HeaderRightBG(dc)
	dc.Fill()

	// draw season title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(lapTitle, s.PaddingSize*3, s.HeaderHeight/4, 0, 0.5)
	// draw week title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(lapWeekTitle, s.ImageWidth/4, s.HeaderHeight/4*3, 0.5, 0.5)
	// draw track title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(lapTrackTitle, s.ImageWidth/3*2, s.HeaderHeight/4*3, 0.5, 0.5)

	// adjust to header height
	yPosColumnHeaderStart := s.HeaderHeight + s.PaddingSize

	// draw division column header
	xDivisionLength := s.DivisionColumnWidth - s.PaddingSize*2
	xPos := s.PaddingSize
	yPos := yPosColumnHeaderStart

	dc.DrawRectangle(xPos, yPos, xDivisionLength, s.ColumnHeaderHeight)
	color.TopNHeaderBG(dc)
	dc.Fill()

	color.TopNHeaderFG(dc)
	if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.DrawStringAnchored("Division", xPos+xDivisionLength/2, yPos+s.ColumnHeaderHeight/2, 0.5, 0.5)

	// draw outline
	color.TopNHeaderOutline(dc)
	dc.MoveTo(xPos, yPos)
	dc.LineTo(xPos+xDivisionLength, yPos)
	dc.LineTo(xPos+xDivisionLength, yPos+s.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos+s.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos)
	dc.SetLineWidth(1)
	dc.Stroke()

	// draw driver column header
	xDriverLength := s.DriverColumnWidth - s.PaddingSize*2
	xPos = xDivisionLength + s.PaddingSize*2

	dc.DrawRectangle(xPos, yPos, xDriverLength, s.ColumnHeaderHeight)
	color.TopNHeaderBG(dc)
	dc.Fill()

	color.TopNHeaderFG(dc)
	if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.DrawStringAnchored("Driver", xPos+xDriverLength/2, yPos+s.ColumnHeaderHeight/2, 0.5, 0.5)

	// draw outline
	color.TopNHeaderOutline(dc)
	dc.MoveTo(xPos, yPos)
	dc.LineTo(xPos+xDriverLength, yPos)
	dc.LineTo(xPos+xDriverLength, yPos+s.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos+s.ColumnHeaderHeight)
	dc.LineTo(xPos, yPos)
	dc.SetLineWidth(1)
	dc.Stroke()

	// draw laptime column headers
	xColumnLength := s.LaptimeColumnWidth - s.PaddingSize
	for column := float64(0); column < s.LaptimeColumns; column++ {
		xPos := xDivisionLength + s.PaddingSize*2 + xDriverLength + s.PaddingSize + float64(column)*s.LaptimeColumnWidth
		yPos := yPosColumnHeaderStart

		title := fmt.Sprintf("%d%%", 100+int(column))
		xLength := xColumnLength
		if column == 0 {
			xLength = xLength + s.PaddingSize
			title = "100%"
		} else {
			xPos = xPos + s.PaddingSize
		}

		dc.DrawRectangle(xPos, yPos, xLength, s.ColumnHeaderHeight)
		color.TopNHeaderBG(dc)
		dc.Fill()

		color.TopNHeaderFG(dc)
		if column == s.LaptimeColumns-1 {
			color.TopNHeaderFGDanger(dc)
		}
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(title, xPos+xLength/2, yPos+s.ColumnHeaderHeight/2, 0.5, 0.5)

		// draw outline
		color.TopNHeaderOutline(dc)
		dc.MoveTo(xPos, yPos)
		dc.LineTo(xPos+xLength, yPos)
		dc.LineTo(xPos+xLength, yPos+s.ColumnHeaderHeight)
		dc.LineTo(xPos, yPos+s.ColumnHeaderHeight)
		dc.LineTo(xPos, yPos)
		dc.SetLineWidth(1)
		dc.Stroke()
	}

	// draw rows
	yPosRowStart := yPosColumnHeaderStart + s.ColumnHeaderHeight + s.PaddingSize
	for row, entry := range s.Data {
		xPos := s.PaddingSize
		yPos := yPosRowStart + float64(row)*s.DriverHeight
		xLength := s.ImageWidth - s.PaddingSize*2

		// zebra pattern
		dc.DrawRectangle(xPos, yPos, xLength, s.DriverHeight)
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
		dc.LineTo(xPos+xLength, yPos+s.DriverHeight)
		dc.LineTo(xPos, yPos+s.DriverHeight)
		dc.LineTo(xPos, yPos)
		dc.SetLineWidth(0.5)
		dc.Stroke()

		// draw division
		color.TopNCellPosition(dc)
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(entry.Division, xPos+xDivisionLength/2, yPos+s.DriverHeight/2, 0.5, 0.5)

		// draw driver
		xPos = xDivisionLength + s.PaddingSize*2
		color.TopNCellDriver(dc)
		// marked driver?
		if entry.Marked {
			color.TopNHeaderFG(dc)
		}
		if err := dc.LoadFontFace("public/fonts/Roboto-Regular.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(entry.Driver, xPos, yPos+s.DriverHeight/2, 0, 0.5)

		// draw laptimes
		xColumnLength := s.LaptimeColumnWidth - s.PaddingSize
		for column := float64(0); column < s.LaptimeColumns; column++ {
			xPos := xDivisionLength + s.PaddingSize*2 + xDriverLength + s.PaddingSize + float64(column)*s.LaptimeColumnWidth

			color.TopNCellValue(dc)
			if column == s.LaptimeColumns-1 {
				color.TopNCellValueDanger(dc)
			}
			if err := dc.LoadFontFace("public/fonts/Roboto-Light.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}

			xLength := xColumnLength
			if column == 0 {
				xLength = xLength + s.PaddingSize
				color.TopNCellValue(dc)
				if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 11); err != nil {
					return fmt.Errorf("could not load font: %v", err)
				}
			} else {
				xPos = xPos + s.PaddingSize
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
			dc.DrawStringAnchored(laptime, xPos+xLength/2, yPos+s.DriverHeight/2, 0.5, 0.5)
		}
	}

	// add border to image
	bdc := gg.NewContext(int(s.ImageWidth+s.BorderSize*2), int(s.ImageHeight+s.BorderSize*2))
	color.Border(bdc)
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(s.BorderSize), int(s.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(s.FooterHeight))
	color.Transparent(fdc)
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	color.LastUpdate(fdc)
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 10); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	lastUpdate := s.Week.LastUpdate.UTC().Format("2006-01-02 15:04:05 -07 MST")
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", lastUpdate), float64(bdc.Width())-s.FooterHeight/2, float64(bdc.Height())+s.FooterHeight/2, 1, 0.5)

	color.CreatedBy(fdc)
	if err := fdc.LoadFontFace("public/fonts/Roboto-Light.ttf", 9); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	fdc.DrawStringAnchored("by Fabio Berchtold", s.FooterHeight/2, float64(bdc.Height())+s.FooterHeight/2, 0, 0.5)

	if err := s.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(s.Filename()) // finally write to file
}
