package heatmap

import (
	"fmt"
	"time"

	"github.com/JamesClonk/iRcollector/database"
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
	"github.com/robfig/cron"
)

type Heatmap struct {
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

func New(season database.Season, week database.RaceWeek, track database.Track, results []database.RaceWeekResult) Heatmap {
	return Heatmap{
		Season:         season,
		Week:           week,
		Track:          track,
		Results:        results,
		BorderSize:     float64(3),
		FooterHeight:   float64(18),
		ImageHeight:    float64(480),
		ImageWidth:     float64(1024),
		HeaderHeight:   float64(45),
		TimeslotHeight: float64(50),
		DayWidth:       float64(160),
		Days:           7, // pretty sure that's never gonna change..
	}
}

func IsAvailable(seasonID, week int) bool {
	return image.IsAvailable("heatmap", seasonID, week)
}

func Filename(seasonID, week int) string {
	return image.ImageFilename("heatmap", seasonID, week)
}

func (h *Heatmap) Filename() string {
	return Filename(h.Season.SeasonID, h.Week.RaceWeek+1)
}

func (h *Heatmap) Draw(minSOF, maxSOF int, drawEmptySlots bool) error {
	// heatmap titles, season + track
	heatmapTitle := fmt.Sprintf("%s - Week %d", h.Season.SeasonName, h.Week.RaceWeek+1)
	heatmap2ndTitle := h.Track.Name
	if h.Week.RaceWeek == -1 { // seasonal avg. map
		heatmapTitle = fmt.Sprintf("%s", h.Season.SeasonName)
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
	weekStart := database.WeekStart(h.Season.StartDate.UTC().AddDate(0, 0, (h.Week.RaceWeek+1)*h.Days))
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

	// create canvas
	dc := gg.NewContext(int(h.ImageWidth), int(h.ImageHeight))

	// background
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, h.ImageWidth, h.HeaderHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()
	dc.DrawRectangle(h.ImageWidth/2+h.DayWidth/2, 0, h.ImageWidth/2, h.HeaderHeight)
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.Fill()

	// draw season name
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 19); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(heatmapTitle, h.DayWidth/7, h.HeaderHeight/2, 0, 0.5)
	// draw track config
	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 19); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored(heatmap2ndTitle, h.ImageWidth-h.DayWidth/7, h.HeaderHeight/2, 1, 0.5)

	// timeslots
	dc.DrawRectangle(0, h.HeaderHeight, h.DayWidth, h.TimeslotHeight)
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.Fill()
	if err := dc.LoadFontFace("public/fonts/roboto-mono_thin.ttf", 14); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	dc.SetRGB255(0, 0, 0) // black
	dc.DrawStringAnchored("UTC / GMT+0", h.DayWidth/2, h.HeaderHeight+h.TimeslotHeight/2, 0.5, 0.5)
	if err := dc.LoadFontFace("public/fonts/roboto-mono_medium.ttf", 16); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	timeslotWidth := ((h.ImageWidth - h.DayWidth) / float64(len(timeslots))) - 1
	for slot := 0; slot < len(timeslots); slot++ {
		dc.DrawRectangle((float64(slot)*(timeslotWidth+1))+(h.DayWidth+1), h.HeaderHeight, timeslotWidth, h.TimeslotHeight)
		if slot%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
		// draw timeslot starting time
		dc.SetRGB255(0, 0, 0) // black
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
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
		// draw weekday name
		dc.SetRGB255(0, 0, 0) // black
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
			dc.SetRGB255(255, 255, 255) // white
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
							dc.SetRGBA255(0, 0, 240-image.MapValueIntoRange(0, 120, minSOF, maxSOF, sof), image.MapValueIntoRange(10, 200, minSOF, maxSOF, sof)) // sof color
							dc.Fill()
						}
					}

					dc.SetRGB255(39, 39, 39) // dark gray 1
					dc.SetLineWidth(1)
					dc.DrawLine(slotX+eventWidth/3, slotY+eventHeight/2, slotX+eventWidth/1.5, slotY+eventHeight/2)
					dc.Stroke()

					dc.SetRGB255(0, 0, 0) // black
					if err := dc.LoadFontFace("public/fonts/roboto-mono_regular.ttf", 15); err != nil {
						return fmt.Errorf("could not load font: %v", err)
					}
					dc.DrawStringAnchored(fmt.Sprintf("%d", sof), slotX+eventWidth/2, slotY+eventHeight/3-1, 0.5, 0.5)
					if err := dc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 13); err != nil {
						return fmt.Errorf("could not load font: %v", err)
					}
					dc.DrawStringAnchored(fmt.Sprintf("%d", result.SizeOfField), slotX+eventWidth/2, slotY+eventHeight/1.5+1, 0.5, 0.5)
				}
			}
		}
	}

	// add border to image
	bdc := gg.NewContext(int(h.ImageWidth+h.BorderSize*2), int(h.ImageHeight+h.BorderSize*2))
	bdc.SetRGB255(39, 39, 39) // dark gray 1
	bdc.Clear()
	bdc.DrawImage(dc.Image(), int(h.BorderSize), int(h.BorderSize))

	// add footer to image
	fdc := gg.NewContext(bdc.Width(), bdc.Height()+int(h.FooterHeight))
	fdc.SetRGBA255(0, 0, 0, 0) // white
	fdc.Clear()
	fdc.DrawImage(bdc.Image(), 0, 0)
	// add last-update text
	fdc.SetRGB255(0, 0, 0) // black
	if err := fdc.LoadFontFace("public/fonts/roboto-mono_light.ttf", 12); err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	fdc.DrawStringAnchored(fmt.Sprintf("Last Update: %s", h.Week.LastUpdate), float64(bdc.Width())-h.FooterHeight/2, float64(bdc.Height())+h.FooterHeight/2, 1, 0.5)

	if err := h.WriteMetadata(); err != nil {
		return err
	}
	return fdc.SavePNG(h.Filename()) // finally write to file
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
