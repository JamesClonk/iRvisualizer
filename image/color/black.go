package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type black struct{}

func NewBlackScheme() Colorizer {
	return &black{}
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
	dc.SetRGB255(144, 33, 33) // muted dark red
	dc.SetRGB255(199, 0, 0) // red 1
*/
func (c *black) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *black) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *black) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *black) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *black) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *black) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *black) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *black) TopNHeaderFGDanger(dc *gg.Context) {
	dc.SetRGB255(199, 0, 0) // red 1
}
func (c *black) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *black) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *black) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *black) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *black) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *black) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *black) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *black) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *black) TopNCellValueDanger(dc *gg.Context) {
	dc.SetRGB255(144, 33, 33) // muted dark red
}
func (c *black) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *black) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *black) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *black) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *black) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *black) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *black) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(180-image.MapValueIntoRange(0, 150, min, max, value), 180-image.MapValueIntoRange(0, 150, min, max, value), 180-image.MapValueIntoRange(0, 150, min, max, value), image.MapValueIntoRange(5, 255, min, max, value)) // sof color
}
func (c *black) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *black) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
