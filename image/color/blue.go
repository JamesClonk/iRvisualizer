package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type blue struct{}

func NewBlueScheme() Colorizer {
	return &blue{}
}

/*
	Colors:
	dc.SetRGB255(0, 0, 0) // black
	dc.SetRGB255(39, 39, 39) // dark gray 1
	dc.SetRGB255(55, 55, 55) // dark gray 2
	dc.SetRGB255(255, 255, 255) // white
	dc.SetRGB255(133, 133, 133) // gray 1
	dc.SetRGB255(155, 155, 155) // gray 2
	dc.SetRGB255(166, 166, 166) // gray 2.5
	dc.SetRGB255(177, 177, 177) // gray 3
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(225, 225, 225) // light gray 1.5
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
func (c *blue) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *blue) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *blue) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *blue) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *blue) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(7, 55, 99) // dark blue 3
}
func (c *blue) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(11, 83, 148) // dark blue 2
}
func (c *blue) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *blue) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *blue) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *blue) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *blue) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *blue) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *blue) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *blue) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *blue) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(7, 55, 99) // dark blue 3
}
func (c *blue) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *blue) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *blue) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *blue) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *blue) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *blue) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *blue) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(0, 0, 240-image.MapValueIntoRange(0, 120, min, max, value), image.MapValueIntoRange(10, 225, min, max, value)) // sof color
}
func (c *blue) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *blue) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(166, 166, 166) // gray 2.5
}
