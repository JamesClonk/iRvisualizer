package heatmap

import (
	"fmt"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	scheme "github.com/JamesClonk/iRvisualizer/image/color"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron"
)

var (
	heatmapDraws = promauto.NewCounter(prometheus.CounterOpts{
		Name: "irvisualizer_heatmaps_drawn_total",
		Help: "Total heatmaps drawn by iRvisualizer.",
	})
)

type Heatmap struct {
	ColorScheme    string
	Season         database.Season
	Week           database.RaceWeek
	Track          database.Track
	Results        []database.RaceWeekResult
	BorderSize     float64
	FooterHeight   float64
	ImageHeight    float64
	ImageWidth     float64
	HeaderHeight   float64
	TimeslotHeight float64
	DayWidth       float64
	Days           int
}

func New(colorScheme string, season database.Season, week database.RaceWeek, track database.Track, results []database.RaceWeekResult) Heatmap {
	return Heatmap{
		ColorScheme:    colorScheme,
		Season:         season,
		Week:           week,
		Track:          track,
		Results:        results,
		BorderSize:     float64(3),
		FooterHeight:   float64(18),
		ImageHeight:    float64(480),
		ImageWidth:     float64(1024),
		HeaderHeight:   float64(38),
		TimeslotHeight: float64(50),
		DayWidth:       float64(160),
		Days:           7, // pretty sure that's never gonna change..
	}
}

func IsAvailable(colorScheme string, seasonID, week int) bool {
	return image.IsAvailable(colorScheme, "heatmap", seasonID, week)
}

func Filename(seasonID, week int) string {
	return image.ImageFilename("heatmap", seasonID, week)
}

func (h *Heatmap) Filename() string {
	return Filename(h.Season.SeasonID, h.Week.RaceWeek+1)
}

func (h *Heatmap) Draw(minSOF, maxSOF int, drawEmptySlots bool) error {
	heatmapDraws.Inc()

	// heatmap titles, season + track
	heatmapTitle := fmt.Sprintf("%s - Week %d", h.Season.SeasonName, h.Week.RaceWeek+1)
	heatmap2ndTitle := h.Track.Name
	if h.Week.RaceWeek == -1 { // seasonal avg. map
		heatmapTitle = h.Season.SeasonName
		heatmap2ndTitle = "Seasonal Average"
	}
	// if len(track.Config) > 0 {
	// 	heatmap2ndTitle = fmt.Sprintf("%s - %s", h.Track.Name, h.Track.Config)
	// }

	log.Infof("draw heatmap for [%s] - [%s]", heatmapTitle, heatmap2ndTitle)

	// figure out timeslots schedule
	p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := p.Parse(h.Season.Timeslots)
	if err != nil {
		return fmt.Errorf("could not parse timeslot [%s] to crontab format: %v", h.Season.Timeslots, err)
	}
	// start -1 minute to previous day, to make sure schedule.Next will catch a midnight start (00:00)
	weekStart := database.WeekStart(h.Season.StartDate.UTC().AddDate(0, 0, (h.Week.RaceWeek)*h.Days))
	start := weekStart.Add(-1 * time.Minute)
	timeslots := make([]time.Time, 0)
	next := schedule.Next(start)                             // get first timeslot
	for next.Before(schedule.Next(start.AddDate(0, 0, 1))) { // collect all timeslots of 1 day
		timeslots = append(timeslots, next)
		next = schedule.Next(next)
	}

	// figure out dynamic SOF
	if minSOF == 0 {
		minSOF = 1000
	}
	if maxSOF == 0 {
		maxSOF = minSOF * 2
		for _, result := range h.Results {
			if result.StrengthOfField > maxSOF {
				maxSOF = result.StrengthOfField
			}
		}
	}

	// colorizer
	color := scheme.Get(h.ColorScheme)

	// create canvas
	dc := gg.NewContext(int(h.ImageWidth), int(h.ImageHeight))

	// background
	color.Background(dc)
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, h.ImageWidth, h.HeaderHeight)
	color.HeaderLeftBG(dc)
	dc.Fill()
	dc.DrawRectangle(h.ImageWidth/2+h.DayWidth/2, 0, h.ImageWidth/2, h.HeaderHeight)
	color.HeaderRightBG(dc)
	dc.Fill()

	// draw season name
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 19); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(heatmapTitle, h.DayWidth/7, h.HeaderHeight/2, 0, 0.5)
	// draw track config
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 19); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeaderFG(dc)
	dc.DrawStringAnchored(heatmap2ndTitle, h.ImageWidth-h.DayWidth/7, h.HeaderHeight/2, 1, 0.5)

	// timeslots
	dc.DrawRectangle(0, h.HeaderHeight, h.DayWidth, h.TimeslotHeight)
	color.HeatmapHeaderDarkerBG(dc)
	dc.Fill()
	if err := dc.LoadFontFace("public/fonts/roboto-mono_thin.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	color.HeatmapHeaderFG(dc)
	dc.DrawStringAnchored("UTC / GMT+0", h.DayWidth/2, h.HeaderHeight+h.TimeslotHeight/2, 0.5, 0.5)
	if err := dc.LoadFontFace("public/fonts/roboto-mono_medium.ttf", 16); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	timeslotWidth := ((h.ImageWidth - h.DayWidth) / float64(len(timeslots))) - 1
	for slot := 0; slot < len(timeslots); slot++ {
		dc.DrawRectangle((float64(slot)*(timeslotWidth+1))+(h.DayWidth+1), h.HeaderHeight, timeslotWidth, h.TimeslotHeight)
		if slot%2 == 0 {
			color.HeatmapHeaderLighterBG(dc)
		} else {
			color.HeatmapHeaderDarkerBG(dc)
		}
		dc.Fill()
		// draw timeslot starting time
		color.HeatmapHeaderFG(dc)
		dc.DrawStringAnchored(
			timeslots[slot].Format("15:04"),
			(float64(slot)*(timeslotWidth+1))+(h.DayWidth+1)+(timeslotWidth/2),
			h.HeaderHeight+h.TimeslotHeight/2,
			0.5, 0.5)
	}

	// weekdays
	if err := dc.LoadFontFace("public/fonts/RobotoCondensed-Regular.ttf", 20); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dayHeight := ((h.ImageHeight - h.HeaderHeight - h.TimeslotHeight) / float64(h.Days)) - 1
	for day := 0; day < h.Days; day++ {
		dc.DrawRectangle(0, (float64(day)*(dayHeight+1))+(h.HeaderHeight+h.TimeslotHeight+1), h.DayWidth, dayHeight)
		if day%2 == 0 {
			color.HeatmapHeaderLighterBG(dc)
		} else {
			color.HeatmapHeaderDarkerBG(dc)
		}
		dc.Fill()
		// draw weekday name
		color.HeatmapHeaderFG(dc)
		dc.DrawStringAnchored(
			weekStart.AddDate(0, 0, day).Weekday().String(),
			h.DayWidth/2,
			(float64(day)*(dayHeight+1))+(h.HeaderHeight+h.TimeslotHeight+1)+dayHeight/2,
			0.5, 0.5)
	}

	// events
	eventHeight := ((h.ImageHeight - h.HeaderHeight - h.TimeslotHeight) / float64(h.Days)) - 1
	eventWidth := ((h.ImageWidth - h.DayWidth) / float64(len(timeslots))) - 1
	for day := 0; day < h.Days; day++ {
		for slot := 0; slot < len(timeslots); slot++ {
			slotX := (float64(slot) * (eventWidth + 1)) + (h.DayWidth + 1)
			slotY := (float64(day) * (eventHeight + 1)) + (h.HeaderHeight + h.TimeslotHeight + 1)

			dc.DrawRectangle(slotX, slotY, eventWidth, eventHeight)
			color.HeatmapTimeslotBG(dc)
			dc.Fill()

			// draw event values
			timeslot := weekStart.AddDate(0, 0, day).Add(time.Hour * time.Duration(timeslots[slot].Hour())).Add(time.Minute * time.Duration(timeslots[slot].Minute()))
			result := image.GetResult(timeslot, h.Results)

			// only draw empty slots if enabled
			if result.Official || drawEmptySlots {
				// only draw event if a session actually happened already
				if timeslot.Before(time.Now().Add(time.Hour * -2)) {
					sof := 0
					if result.Official {
						sof = result.StrengthOfField
						if result.StrengthOfField >= minSOF {
							// draw background color
							dc.DrawRectangle(slotX, slotY, eventWidth, eventHeight)
							color.HeatmapTimeslotMapping(dc, minSOF, maxSOF, sof) // sof color
							dc.Fill()
						}
					}

					color.HeatmapTimeslotFG(dc)
					dc.SetLineWidth(1)
					dc.DrawLine(slotX+eventWidth/3, slotY+eventHeight/2, slotX+eventWidth/1.5, slotY+eventHeight/2)
					dc.Stroke()

					if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 15); err != nil {
						return fmt.Errorf("could not load font: %v", err)
					}
					textWithBorder(dc, color, fmt.Sprintf("%d", sof), slotX+eventWidth/2, slotY+eventHeight/3-1)

					if err := dc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 13); err != nil {
						return fmt.Errorf("could not load font: %v", err)
					}
					textWithBorder(dc, color, fmt.Sprintf("%d", result.SizeOfField), slotX+eventWidth/2, slotY+eventHeight/1.5+1)
				}
			}
		}
	}

	// add border to image
	bdc := gg.NewContext(int(h.ImageWidth+h.BorderSize*2), int(h.ImageHeight+h.BorderSize*2))
	color.Border(bdc)
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(h.BorderSize), int(h.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(h.FooterHeight))
	color.Transparent(fdc)
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	color.LastUpdate(fdc)
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	lastUpdate := h.Week.LastUpdate.UTC().Format("2006-01-02 15:04:05 -07 MST")
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", lastUpdate), float64(bdc.Width())-h.FooterHeight/2, float64(bdc.Height())+h.FooterHeight/2, 1, 0.5)

	color.CreatedBy(fdc)
	if err := fdc.LoadFontFace("public/fonts/Roboto-Light.ttf", 10); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	fdc.DrawStringAnchored("by Fabio Berchtold", h.FooterHeight/2, float64(bdc.Height())+h.FooterHeight/2, 0, 0.5)

	if err := h.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(h.Filename()) // finally write to file
}

func textWithBorder(dc *gg.Context, color scheme.Colorizer, text string, X, Y float64) {
	if text != "0" {
		color.Border(dc)
		n := 1
		for dy := -n; dy <= n; dy++ {
			for dx := -n; dx <= n; dx++ {
				x := X + float64(dx)
				y := Y + float64(dy)
				dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
			}
		}
		color.HeatmapTimeslotBG(dc)
	} else {
		color.HeatmapTimeslotZero(dc)
	}
	dc.DrawStringAnchored(text, X, Y, 0.5, 0.5)
}
