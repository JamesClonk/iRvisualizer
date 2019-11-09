package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type red struct{}

func NewRedScheme() Colorizer {
	return &red{}
}

/*
	Colors:
	dc.SetRGB255(0, 0, 0) // black
	dc.SetRGB255(39, 39, 39) // dark gray 1
	dc.SetRGB255(55, 55, 55) // dark gray 2
	dc.SetRGB255(255, 255, 255) // white
	dc.SetRGB255(133, 133, 133) // gray 1
	dc.SetRGB255(155, 155, 155) // gray 2
	dc.SetRGB255(177, 177, 177) // gray 3
	dc.SetRGB255(217, 217, 217) // light gray 1
	dc.SetRGB255(225, 225, 225) // light gray 1.5
	dc.SetRGB255(239, 239, 239) // light gray 2
	dc.SetRGB255(241, 241, 241) // light gray 2.5
	dc.SetRGB255(243, 243, 243) // light gray 3
	dc.SetRGB255(85, 0, 0) // dark red 1
	dc.SetRGB255(120, 0, 0) // dark red 2
	dc.SetRGB255(165, 0, 0) // dark red 3
	dc.SetRGB255(240, 50, 50) // light red 1
*/
func (c *red) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *red) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *red) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *red) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *red) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(120, 0, 0) // dark red 2
}
func (c *red) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(165, 0, 0) // dark red 3
}
func (c *red) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *red) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(240, 50, 50) // light red 1
}
func (c *red) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *red) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *red) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *red) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *red) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *red) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *red) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(85, 0, 0) // dark red 1
}
func (c *red) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *red) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *red) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *red) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *red) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *red) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(250-image.MapValueIntoRange(50, 0, min, max, value), 50-image.MapValueIntoRange(0, 45, min, max, value), 0, image.MapValueIntoRange(5, 255, min, max, value)) // sof color
}
func (c *red) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
