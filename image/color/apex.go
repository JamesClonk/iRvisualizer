package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type apex struct{}

func NewApexScheme() Colorizer {
	return &apex{}
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
	dc.SetRGB255(2, 35, 43) // dark green 1
	dc.SetRGB255(33, 144, 55) // dark green 2
	dc.SetRGB255(27, 204, 110) // lime green
*/
func (c *apex) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *apex) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *apex) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *apex) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *apex) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(27, 204, 110) // lime green
}
func (c *apex) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *apex) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *apex) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *apex) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *apex) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(33, 144, 55) // dark green 2
}
func (c *apex) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *apex) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *apex) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(2, 35, 43) // dark green 1
}
func (c *apex) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *apex) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *apex) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(33-image.MapValueIntoRange(22, 0, min, max, value), 240-image.MapValueIntoRange(0, 150, min, max, value), 77-image.MapValueIntoRange(55, 0, min, max, value), image.MapValueIntoRange(5, 255, min, max, value)) // sof color
}
func (c *apex) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *apex) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
