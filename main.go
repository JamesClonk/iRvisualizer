package main

import "github.com/fogleman/gg"

var (
	imageHeight    = float64(480)
	imageLength    = float64(1024)
	headerHeight   = float64(45)
	timeslotHeight = float64(50)
	dayLength      = float64(160)
)

func main() {
	dc := gg.NewContext(int(imageLength), int(imageHeight))

	// background
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.Clear()

	// header
	dc.DrawRectangle(0, 0, imageLength, headerHeight)
	dc.SetRGB255(7, 55, 99) // dark blue 3
	dc.Fill()

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

	dc.SavePNG("public/test.png")
}

/*
	Colors:
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
