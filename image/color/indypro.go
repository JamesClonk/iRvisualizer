package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type pm18 struct{}

func NewPMScheme() Colorizer {
	return &pm18{}
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
	dc.SetRGB255(144, 33, 33) // muted dark red
	dc.SetRGB255(155, 0, 0) // dark red 2.5
*/
func (c *pm18) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *pm18) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *pm18) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *pm18) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *pm18) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(7, 55, 99) // dark blue 3
}
func (c *pm18) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(11, 83, 148) // dark blue 2
}
func (c *pm18) TopNHeaderFGDanger(dc *gg.Context) {
	dc.SetRGB255(155, 0, 0) // dark red 2.5
}
func (c *pm18) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *pm18) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *pm18) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *pm18) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *pm18) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *pm18) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *pm18) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *pm18) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *pm18) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(7, 55, 99) // dark blue 3
}
func (c *pm18) TopNCellValueDanger(dc *gg.Context) {
	dc.SetRGB255(144, 33, 33) // muted dark red
}
func (c *pm18) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *pm18) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *pm18) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *pm18) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *pm18) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *pm18) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *pm18) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(50-image.MapValueIntoRange(0, 45, min, max, value), 150-image.MapValueIntoRange(0, 120, min, max, value), 255-image.MapValueIntoRange(0, 160, min, max, value), image.MapValueIntoRange(10, 225, min, max, value)) // sof color
}
func (c *pm18) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *pm18) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
