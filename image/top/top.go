package top

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
	topNDraws = promauto.NewCounter(prometheus.CounterOpts{
		Name: "irvisualizer_topn_drawn_total",
		Help: "Total topN drawn by iRvisualizer.",
	})
)

type DataSet struct {
	Title string
	Icons string
	Rows  []DataSetRow
}

type DataSetRow struct {
	Driver       string
	Value        string
	Icon         string
	IconPosition int
}

type Top struct {
	ColorScheme  string
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

func New(colorScheme, name string, season database.Season, week database.RaceWeek, track database.Track, data []DataSet) Top {
	top := Top{
		ColorScheme:  colorScheme,
		Name:         name,
		Season:       season,
		Week:         week,
		Track:        track,
		Data:         data,
		BorderSize:   float64(2),
		FooterHeight: float64(14),
		ImageWidth:   float64(740),
		HeaderHeight: float64(24),
		DriverHeight: float64(16),
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

func IsAvailable(colorScheme string, name string, seasonID, week int) bool {
	return image.IsAvailable(colorScheme, "top/"+name, seasonID, week)
}

func Filename(name string, seasonID, week int) string {
	return image.ImageFilename("top/"+name, seasonID, week)
}

func (t *Top) Filename() string {
	return Filename(t.Name, t.Season.SeasonID, t.Week.RaceWeek+1)
}

func (t *Top) Draw(headerless bool) error {
	topNDraws.Inc()

	// top titles, season + track
	topTitle := fmt.Sprintf("%s - Statistics", t.Season.SeasonName)
	if len(t.Season.SeasonName) > 38 {
		topTitle = t.Season.SeasonName
	}
	topTrackTitle := fmt.Sprintf("Week %d - %s", t.Week.RaceWeek+1, t.Track.Name)
	if t.Week.RaceWeek == -1 { // seasonal avg. top
		topTrackTitle = "Seasonal Average"
	}

	log.Infof("draw top for [%s] - [%s]", topTitle, topTrackTitle)

	// colorizer
	color := scheme.Get(t.ColorScheme)

	// strip header?
	if headerless {
		t.ImageHeight = t.ImageHeight - t.HeaderHeight
	}

	// create canvas
	dc := gg.NewContext(int(t.ImageWidth), int(t.ImageHeight))

	// background
	color.Background(dc)
	dc.Clear()

	yPosColumnHeaderStart := t.PaddingSize
	if !headerless {
		// header
		dc.DrawRectangle(0, 0, t.ImageWidth, t.HeaderHeight)
		color.HeaderLeftBG(dc)
		dc.Fill()
		dc.DrawRectangle(t.ImageWidth/2, 0, t.ImageWidth/2, t.HeaderHeight)
		color.HeaderRightBG(dc)
		dc.Fill()

		// draw season name
		if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		color.HeaderFG(dc)
		dc.DrawStringAnchored(topTitle, t.ImageWidth/4, t.HeaderHeight/2, 0.5, 0.5)
		// draw track title
		if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		color.HeaderFG(dc)
		dc.DrawStringAnchored(topTrackTitle, t.ImageWidth/2+t.ImageWidth/4, t.HeaderHeight/2, 0.5, 0.5)

		// adjust to header height
		yPosColumnHeaderStart = t.HeaderHeight + t.PaddingSize
	}

	// draw the column headers
	xLength := t.ColumnWidth - t.PaddingSize*2
	for column, data := range t.Data {
		xPos := t.PaddingSize + float64(column)*t.ColumnWidth
		yPos := yPosColumnHeaderStart

		// add column header
		dc.DrawRectangle(xPos, yPos, xLength, t.DriverHeight)
		color.TopNHeaderBG(dc)
		dc.Fill()

		color.TopNHeaderFG(dc)
		if err := dc.LoadFontFace("public/fonts/Roboto-Medium.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(data.Title, xPos+xLength/2, yPos+t.DriverHeight/2, 0.5, 0.5)

		// draw outline
		color.TopNHeaderOutline(dc)
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
				color.TopNCellDarkerBG(dc)
			} else {
				color.TopNCellLighterBG(dc)
			}
			dc.Fill()

			// position
			color.TopNCellPosition(dc)
			if err := dc.LoadFontFace("public/fonts/Roboto-Light.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			if entry.Value != previousValue {
				previousValue = entry.Value

				// draw icons if specified
				if len(data.Icons) > 0 && row <= 2 {
					// load icon
					iconColor := "gold"
					if row == 1 {
						iconColor = "silver"
					}
					if row == 2 {
						iconColor = "bronze"
					}
					icon, err := gg.LoadPNG(fmt.Sprintf("public/icons/%s_%s.png", data.Icons, iconColor))
					if err != nil {
						return fmt.Errorf("could not load icon: %v", err)
					}
					dc.DrawImage(icon, int(xPos+t.PaddingSize), int(yPos))
				} else {
					dc.DrawStringAnchored(fmt.Sprintf("%d.", row+1), xPos+t.PaddingSize*2, yPos+t.DriverHeight/2, 0, 0.5)
				}
			}
			// name
			color.TopNCellDriver(dc)
			if err := dc.LoadFontFace("public/fonts/Roboto-Regular.ttf", 11); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			dc.DrawStringAnchored(entry.Driver, xPos+20+t.PaddingSize*2, yPos+t.DriverHeight/2, 0, 0.5)
			// value
			color.TopNCellValue(dc)
			if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 12); err != nil {
				return fmt.Errorf("could not load font: %v", err)
			}
			dc.DrawStringAnchored(entry.Value, xPos+xLength-t.PaddingSize*2, yPos+t.DriverHeight/2, 1, 0.5)
			// draw an icon if specified
			if len(entry.Icon) > 0 {
				icon, err := gg.LoadPNG(fmt.Sprintf("public/icons/%s.png", entry.Icon))
				if err != nil {
					return fmt.Errorf("could not load icon: %v", err)
				}
				dc.DrawImageAnchored(icon, int(xPos+xLength-t.PaddingSize*2)-entry.IconPosition, int(yPos), 1, 0)
			}

			// draw outline
			color.TopNCellOutline(dc)
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
	color.Border(bdc)
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(t.BorderSize), int(t.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(t.FooterHeight))
	color.Transparent(fdc)
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	color.LastUpdate(fdc)
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 10); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	lastUpdate := t.Week.LastUpdate.UTC().Format("2006-01-02 15:04:05 -07 MST")
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", lastUpdate), float64(bdc.Width())-t.FooterHeight/2, float64(bdc.Height())+t.FooterHeight/2, 1, 0.5)

	color.CreatedBy(fdc)
	if err := fdc.LoadFontFace("public/fonts/Roboto-Light.ttf", 9); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	fdc.DrawStringAnchored("created by Fabio Berchtold", t.FooterHeight/2, float64(bdc.Height())+t.FooterHeight/2, 0, 0.5)

	if err := t.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(t.Filename()) // finally write to file
}
