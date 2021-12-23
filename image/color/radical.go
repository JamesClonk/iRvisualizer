package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type radical struct{}

func NewRadicalScheme() Colorizer {
	return &radical{}
}

/*
	Colors:
	dc.SetRGB255(0, 0, 0) // black
	dc.SetRGB255(39, 39, 39) // dark gray 1
	dc.SetRGB255(55, 55, 55) // dark gray 2
	dc.SetRGB255(77, 77, 77) // dark gray 3
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
	dc.SetRGB255(190, 90, 0) // dark orange 1
	dc.SetRGB255(220, 145, 0) // dark yellow 1
	dc.SetRGB255(245, 180, 0) // dark yellow 2
	dc.SetRGB255(250, 190, 0) // dark yellow 3
	dc.SetRGB255(255, 205, 0) // light yellow 1
	dc.SetRGB255(144, 33, 33) // muted dark red
	dc.SetRGB255(199, 0, 0) // red 1
	dc.SetRGB255(222, 0, 0) // red 2
	dc.SetRGB255(155, 0, 0) // dark red 2.5
*/
func (c *radical) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *radical) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *radical) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *radical) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *radical) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(250, 190, 0) // dark yellow 3
}
func (c *radical) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(255, 205, 0) // light yellow 1
}
func (c *radical) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(222, 0, 0) // red 2
}
func (c *radical) TopNHeaderFGDanger(dc *gg.Context) {
	dc.SetRGB255(155, 0, 0) // dark red 2.5
}
func (c *radical) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(255, 205, 0) // light yellow 1
}
func (c *radical) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(77, 77, 77) // dark gray 3
}
func (c *radical) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *radical) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *radical) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *radical) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *radical) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *radical) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(144, 33, 33) // muted dark red
}
func (c *radical) TopNCellValueDanger(dc *gg.Context) {
	dc.SetRGB255(199, 0, 0) // red 1
}
func (c *radical) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *radical) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *radical) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *radical) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *radical) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *radical) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *radical) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(255-image.MapValueIntoRange(15, 0, min, max, value), 220-image.MapValueIntoRange(0, 50, min, max, value), 0, image.MapValueIntoRange(5, 255, min, max, value)) // sof color
}
func (c *radical) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *radical) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
