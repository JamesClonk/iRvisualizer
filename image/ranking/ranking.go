package ranking

import (
	"fmt"
	"math"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	scheme "github.com/JamesClonk/iRvisualizer/image/color"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
)

type DataRow struct {
	Driver string
	Value  string
}

type Ranking struct {
	Season       database.Season
	Data         []DataRow
	BorderSize   float64
	FooterHeight float64
	ImageHeight  float64
	ImageWidth   float64
	HeaderHeight float64
	DriverHeight float64
	PaddingSize  float64
	Columns      float64
	ColumnWidth  float64
	Rows         float64
}

func New(season database.Season, data []DataRow) Ranking {
	ranking := Ranking{
		Season:       season,
		Data:         data,
		BorderSize:   float64(2),
		FooterHeight: float64(14),
		ImageWidth:   float64(740),
		HeaderHeight: float64(24),
		DriverHeight: float64(16),
		PaddingSize:  float64(3),
		Columns:      float64(3),
		Rows:         float64(10),
	}
	ranking.ColumnWidth = ranking.ImageWidth / ranking.Columns
	ranking.ImageHeight = float64(ranking.Rows)*ranking.DriverHeight + ranking.HeaderHeight + ranking.PaddingSize*2
	return ranking
}

func IsAvailable(seasonID int) bool {
	return image.IsAvailable("ranking", seasonID, -1)
}

func Filename(seasonID int) string {
	return image.ImageFilename("ranking", seasonID, -1)
}

func (r *Ranking) Filename() string {
	return Filename(r.Season.SeasonID)
}

func (r *Ranking) Draw(colorScheme string, num, ofTotal int) error {
	// raking title
	rankingTitle := fmt.Sprintf("%s - Standings", r.Season.SeasonName)
	if len(r.Season.SeasonName) > 64 {
		rankingTitle = r.Season.SeasonName
	}
	rankingBestOfTitle := fmt.Sprintf("Best %d out of %d week", num, ofTotal)
	if ofTotal > 1 {
		rankingBestOfTitle += "s" // plural
	}

	log.Infof("draw ranking for [%s] - [%s]", rankingTitle, rankingBestOfTitle)

	// colorizer
	var color scheme.Colorizer
	switch colorScheme {
	case "yellow":
		color = scheme.NewYellowScheme()
	case "red":
		color = scheme.NewRedScheme()
	case "simucube":
		color = scheme.NewSimuCubeScheme()
	default:
		color = scheme.NewBlueScheme()
	}

	// create canvas
	dc := gg.NewContext(int(r.ImageWidth), int(r.ImageHeight))

	// background
	color.Background(dc)
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, r.ImageWidth, r.HeaderHeight)
	color.HeaderLeftBG(dc)
	dc.Fill()
	dc.DrawRectangle(r.ImageWidth/1.5, 0, r.ImageWidth/3, r.HeaderHeight)
	color.HeaderRightBG(dc)
	dc.Fill()

	// draw season ranking title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(rankingTitle, r.ImageWidth/3, r.HeaderHeight/2, 0.5, 0.5)
	// draw best-of title
	if err := dc.LoadFontFace("public/fonts/Roboto-BoldItalic.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(rankingBestOfTitle, r.ImageWidth/2+r.ImageWidth/3, r.HeaderHeight/2, 0.5, 0.5)

	// adjust to header height
	yPosColumnStart := r.HeaderHeight + r.PaddingSize
	// draw the columns & rows
	var previousValue string
	xLength := r.ColumnWidth - r.PaddingSize*2
	for d, data := range r.Data {
		if float64(d) >= r.Rows*r.Columns {
			break // abort if too many data rows supplied
		}
		column := math.Floor(float64(d) / r.Rows) // calculate current column based on row index / how many rows per column
		row := float64(d) - (column * r.Rows)     // calculate on which row index of current column
		xPos := r.PaddingSize + column*r.ColumnWidth
		yPos := yPosColumnStart + float64(row)*r.DriverHeight

		// zebra pattern
		dc.DrawRectangle(xPos, yPos, xLength, r.DriverHeight)
		if int(row)%2 == 0 {
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
		if data.Value != previousValue {
			previousValue = data.Value
			// draw trophies
			if d <= 2 {
				// load icon
				iconColor := "gold"
				if d == 1 {
					iconColor = "silver"
				}
				if d == 2 {
					iconColor = "bronze"
				}
				icon, err := gg.LoadPNG(fmt.Sprintf("public/icons/trophy_%s.png", iconColor))
				if err != nil {
					return fmt.Errorf("could not load icon: %v", err)
				}
				dc.DrawImage(icon, int(xPos+r.PaddingSize), int(yPos))
			} else {
				dc.DrawStringAnchored(fmt.Sprintf("%d.", int(row+1)+int(r.Rows*(column))), xPos+r.PaddingSize*2, yPos+r.DriverHeight/2, 0, 0.5)
			}
		}
		// name
		color.TopNCellDriver(dc)
		if err := dc.LoadFontFace("public/fonts/Roboto-Regular.ttf", 11); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(data.Driver, xPos+20+r.PaddingSize*2, yPos+r.DriverHeight/2, 0, 0.5)
		// value
		color.TopNCellValue(dc)
		if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 12); err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		dc.DrawStringAnchored(data.Value, xPos+xLength-r.PaddingSize*2, yPos+r.DriverHeight/2, 1, 0.5)

		// draw outline
		color.TopNCellOutline(dc)
		dc.MoveTo(xPos, yPos)
		dc.LineTo(xPos+xLength, yPos)
		dc.LineTo(xPos+xLength, yPos+r.DriverHeight)
		dc.LineTo(xPos, yPos+r.DriverHeight)
		dc.LineTo(xPos, yPos)
		dc.SetLineWidth(0.5)
		dc.Stroke()
	}

	// add border to image
	bdc := gg.NewContext(int(r.ImageWidth+r.BorderSize*2), int(r.ImageHeight+r.BorderSize*2))
	color.Border(bdc)
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(r.BorderSize), int(r.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(r.FooterHeight))
	color.Transparent(fdc)
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	color.LastUpdate(fdc)
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 10); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	lastUpdate := time.Now().Add(-2 * time.Hour).UTC().Format("2006-01-02 15:04:05 -07 MST")
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", lastUpdate), float64(bdc.Width())-r.FooterHeight/2, float64(bdc.Height())+r.FooterHeight/2, 1, 0.5)

	if err := r.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(r.Filename()) // finally write to file
}
