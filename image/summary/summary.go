package summary

import (
	"fmt"
	"strconv"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	scheme "github.com/JamesClonk/iRvisualizer/image/color"
	"github.com/JamesClonk/iRvisualizer/log"
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
	Summary database.Summary
	Marked  bool
}

type Summary struct {
	ColorScheme        string
	Team               string
	Name               string
	Season             database.Season
	Week               database.RaceWeek
	Track              database.Track
	Data               []DataSet
	BorderSize         float64
	FooterHeight       float64
	ImageHeight        float64
	ImageWidth         float64
	HeaderHeight       float64
	ColumnHeaderHeight float64
	DriverHeight       float64
	PaddingSize        float64
	Rows               float64
	SummaryColumns     float64
	SummaryColumnWidth float64
	PointsColumnWidth  float64
	DriverColumnWidth  float64
}

func New(colorScheme, team string, season database.Season, week database.RaceWeek, track database.Track, data []DataSet) Summary {
	lap := Summary{
		ColorScheme:        colorScheme,
		Team:               team,
		Name:               "summary",
		Season:             season,
		Week:               week,
		Track:              track,
		Data:               data,
		BorderSize:         float64(2),
		FooterHeight:       float64(14),
		ImageWidth:         float64(756),
		HeaderHeight:       float64(46),
		ColumnHeaderHeight: float64(16),
		DriverHeight:       float64(24),
		PaddingSize:        float64(3),
		Rows:               float64(len(data)),
		SummaryColumns:     float64(10),
		SummaryColumnWidth: float64(50),
		PointsColumnWidth:  float64(48),
	}
	lap.DriverColumnWidth = lap.ImageWidth - (lap.PointsColumnWidth + (lap.SummaryColumnWidth * lap.SummaryColumns))
	lap.ImageHeight = lap.Rows*lap.DriverHeight + lap.ColumnHeaderHeight + lap.HeaderHeight + lap.PaddingSize*3
	return lap
}

func IsAvailable(colorScheme string, seasonID, week int, team string) bool {
	return image.IsAvailable(colorScheme, "summary", seasonID, week, team)
}

func Filename(seasonID, week int, team string) string {
	return image.ImageFilename("summary", seasonID, week, team)
}

func (s *Summary) Filename() string {
	return Filename(s.Season.SeasonID, s.Week.RaceWeek+1, s.Team)
}

func (s *Summary) Draw() error {
	summaryDraws.Inc()

	// summary titles, season + track
	summaryTitle := fmt.Sprintf("%s - Driver summary", s.Season.SeasonName)
	if len(s.Season.SeasonName) > 64 {
		summaryTitle = s.Season.SeasonName
	}
	summaryWeekTitle := fmt.Sprintf("Week %d", s.Week.RaceWeek+1)
	summaryTrackTitle := s.Track.Name

	if s.Week.RaceWeek == -1 { // seasonal summary
		summaryTitle = s.Season.SeasonName
		summaryWeekTitle = ""
		summaryTrackTitle = "Seasonal driver summary"
	}

	log.Infof("draw summary for [%s] - [%s]", summaryTitle, summaryTrackTitle)

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
	dc.DrawStringAnchored(summaryTitle, s.PaddingSize*3, s.HeaderHeight/4, 0, 0.5)
	// draw week title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(summaryWeekTitle, s.ImageWidth/4, s.HeaderHeight/4*3, 0.5, 0.5)
	// draw track title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(summaryTrackTitle, s.ImageWidth/3*2, s.HeaderHeight/4*3, 0.5, 0.5)

	// adjust to header height
	yPosColumnHeaderStart := s.HeaderHeight + s.PaddingSize

	// draw division column header
	xDivisionLength := s.PointsColumnWidth - s.PaddingSize*2
	xPos := s.PaddingSize
	yPos := yPosColumnHeaderStart

	dc.DrawRectangle(xPos, yPos, xDivisionLength, s.ColumnHeaderHeight)
	color.TopNHeaderBG(dc)
	dc.Fill()

	color.TopNHeaderFG(dc)
	if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.DrawStringAnchored("HPts", xPos+xDivisionLength/2, yPos+s.ColumnHeaderHeight/2, 0.5, 0.5)

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

	// draw summary column headers
	titles := []string{"+iR", "+SR", "Races", "Win", "Pole", "Podium", "Top 5", "+Pos", "LLaps", "Inc/L"}
	xColumnLength := s.SummaryColumnWidth - s.PaddingSize
	for column := float64(0); column < s.SummaryColumns; column++ {
		xPos := xDivisionLength + s.PaddingSize*2 + xDriverLength + s.PaddingSize + float64(column)*s.SummaryColumnWidth
		yPos := yPosColumnHeaderStart

		title := titles[int(column)]
		xLength := xColumnLength
		if column == 0 {
			xLength = xLength + s.PaddingSize
		} else {
			xPos = xPos + s.PaddingSize
		}

		dc.DrawRectangle(xPos, yPos, xLength, s.ColumnHeaderHeight)
		color.TopNHeaderBG(dc)
		dc.Fill()

		color.TopNHeaderFG(dc)
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
		color.TopNCellValueDanger(dc)
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(strconv.Itoa(entry.Summary.HighestChampPoints), xPos+xDivisionLength/2, yPos+s.DriverHeight/2, 0.5, 0.5)

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
		dc.DrawStringAnchored(entry.Summary.Driver.Name, xPos, yPos+s.DriverHeight/2, 0, 0.5)

		// draw summary
		xColumnLength := s.SummaryColumnWidth - s.PaddingSize
		for column := float64(0); column < s.SummaryColumns; column++ {
			xPos := xDivisionLength + s.PaddingSize*2 + xDriverLength + s.PaddingSize + float64(column)*s.SummaryColumnWidth
			xLength := xColumnLength
			if column == 0 {
				xLength = xLength + s.PaddingSize
			} else {
				xPos = xPos + s.PaddingSize
			}

			color.TopNCellValue(dc)
			var value string
			switch column {
			case 0:
				value = strconv.Itoa(entry.Summary.TotalIRatingGain)
				if entry.Summary.TotalIRatingGain < 0 {
					color.TopNCellValueDanger(dc)
				}
			case 1:
				value = strconv.Itoa(entry.Summary.TotalSafetyRatingGain)
				if entry.Summary.TotalSafetyRatingGain < 0 {
					color.TopNCellValueDanger(dc)
				}
			case 2:
				value = strconv.Itoa(entry.Summary.NumberOfRaces)
			case 3:
				value = strconv.Itoa(entry.Summary.Wins)
			case 4:
				value = strconv.Itoa(entry.Summary.Poles)
			case 5:
				value = strconv.Itoa(entry.Summary.Podiums)
			case 6:
				value = strconv.Itoa(entry.Summary.Top5)
			case 7:
				value = strconv.Itoa(entry.Summary.TotalPositionsGained)
				if entry.Summary.TotalPositionsGained < 0 {
					color.TopNCellValueDanger(dc)
				}
			case 8:
				value = strconv.Itoa(entry.Summary.LapsLead)
			case 9:
				value = fmt.Sprintf("%0.2f", entry.Summary.AverageIncidentsPerLap)
				if entry.Summary.AverageIncidentsPerLap > 0.5 {
					color.TopNCellValueDanger(dc)
				}
			}

			// marked driver?
			if entry.Marked {
				color.TopNHeaderFGDanger(dc)
			}

			if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			dc.DrawStringAnchored(value, xPos+xLength/2, yPos+s.DriverHeight/2, 0.5, 0.5)
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
