package main

import (
	"github.com/JamesClonk/iRvisualizer/log"
	"github.com/fogleman/gg"
)

const (
	imageHeight    = float64(480)
	imageLength    = float64(1024)
	headerHeight   = float64(45)
	timeslotHeight = float64(50)
	dayLength      = float64(160)
)

func main() {
	dc := gg.NewContext(int(imageLength), int(imageHeight))

	// background
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, imageLength, headerHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()
	dc.DrawRectangle(imageLength/2+dayLength/2, 0, imageLength/2, headerHeight)
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.Fill()

	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 20); err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored("Pro Mazda Championship - 2019 Season 2 - Week 10", dayLength/4, headerHeight/2, 0, 0.5)

	if err := dc.LoadFontFace("public/fonts/Roboto-Italic.ttf", 20); err != nil {
		log.Fatalf("could not load font: %v", err)
	}
	dc.SetRGB255(255, 255, 255) // white
	dc.DrawStringAnchored("Spa-Francorchamps - Classic Pits", imageLength-dayLength/4, headerHeight/2, 1, 0.5)

	// timeslots
	dc.DrawRectangle(0, headerHeight, dayLength, timeslotHeight)
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.Fill()
	slots := 12
	timeslotLength := ((imageLength - dayLength) / float64(slots)) - 1
	for slot := 0; slot < slots; slot++ {
		dc.DrawRectangle((float64(slot)*(timeslotLength+1))+(dayLength+1), headerHeight, timeslotLength, timeslotHeight)
		if slot%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
	}

	// weekdays
	days := 7
	dayHeight := ((imageHeight - headerHeight - timeslotHeight) / float64(days)) - 1
	for day := 0; day < days; day++ {
		dc.DrawRectangle(0, (float64(day)*(dayHeight+1))+(headerHeight+timeslotHeight+1), dayLength, dayHeight)
		if day%2 == 0 {
			dc.SetRGB255(243, 243, 243) // light gray 3
		} else {
			dc.SetRGB255(239, 239, 239) // light gray 2
		}
		dc.Fill()
	}

	// empty events
	eventDays := 7
	eventSlots := 12
	eventHeight := ((imageHeight - headerHeight - timeslotHeight) / float64(eventDays)) - 1
	eventLength := ((imageLength - dayLength) / float64(eventSlots)) - 1
	for day := 0; day < eventDays; day++ {
		for slot := 0; slot < eventSlots; slot++ {
			dc.DrawRectangle(
				(float64(slot)*(eventLength+1))+(dayLength+1),
				(float64(day)*(eventHeight+1))+(headerHeight+timeslotHeight+1),
				eventLength, eventHeight)
			dc.SetRGB255(255, 255, 255) // white
			dc.Fill()
		}
	}

	dc.SavePNG("public/test.png")
}

/*
	Colors:
	dc.SetRGB255(255, 255, 255) // white
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
