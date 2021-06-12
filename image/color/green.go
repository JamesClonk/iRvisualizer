package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type green struct{}

func NewGreenScheme() Colorizer {
	return &green{}
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
	dc.SetRGB255(0, 80, 0) // dark green 1
	dc.SetRGB255(0, 111, 0) // dark green 2
	dc.SetRGB255(22, 133, 22) // dark green 3
	dc.SetRGB255(44, 200, 44) // light green 1
	dc.SetRGB255(144, 33, 33) // muted dark red
	dc.SetRGB255(199, 0, 0) // red 1
*/
func (c *green) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *green) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *green) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *green) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *green) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(0, 111, 0) // dark green 2
}
func (c *green) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(22, 133, 22) // dark green 3
}
func (c *green) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *green) TopNHeaderFGDanger(dc *gg.Context) {
	dc.SetRGB255(199, 0, 0) // red 1
}
func (c *green) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(44, 200, 44) // light green
}
func (c *green) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *green) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *green) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *green) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *green) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *green) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *green) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(0, 80, 0) // dark green 1
}
func (c *green) TopNCellValueDanger(dc *gg.Context) {
	dc.SetRGB255(144, 33, 33) // muted dark red
}
func (c *green) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *green) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *green) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *green) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *green) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *green) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *green) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(0, 180-image.MapValueIntoRange(0, 120, min, max, value), 0, image.MapValueIntoRange(5, 255, min, max, value)) // sof color
}
func (c *green) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *green) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
