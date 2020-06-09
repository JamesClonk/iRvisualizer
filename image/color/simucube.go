package color

import (
	"github.com/JamesClonk/iRvisualizer/image"
	"github.com/fogleman/gg"
)

type sc struct{}

func NewSimuCubeScheme() Colorizer {
	return &sc{}
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
	dc.SetRGB255(2, 35, 43) // dark green 1
	dc.SetRGB255(27, 51, 56) // dark green 2
	dc.SetRGB255(40, 68, 75) // dark green 3
	dc.SetRGB255(232, 78, 15) // light orange 1
	dc.SetRGB255(200, 40, 10) // dark orange 1
*/
func (c *sc) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *sc) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *sc) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *sc) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(232, 78, 15) // light orange 1
}
func (c *sc) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(2, 35, 43) // dark green 1
}
func (c *sc) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(27, 51, 56) // dark green 2
}
func (c *sc) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *sc) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(40, 68, 75) // dark green 3
}
func (c *sc) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *sc) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *sc) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *sc) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *sc) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *sc) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *sc) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(200, 40, 10) // dark orange 1
}
func (c *sc) HeatmapHeaderFG(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *sc) HeatmapHeaderDarkerBG(dc *gg.Context) {
	dc.SetRGB255(239, 239, 239) // light gray 2
}
func (c *sc) HeatmapHeaderLighterBG(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
func (c *sc) HeatmapTimeslotFG(dc *gg.Context) {
	dc.SetRGB255(2, 35, 43) // dark green 1
}
func (c *sc) HeatmapTimeslotBG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *sc) HeatmapTimeslotZero(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *sc) HeatmapTimeslotMapping(dc *gg.Context, min, max, value int) {
	dc.SetRGBA255(220-image.MapValueIntoRange(50, 0, min, max, value), 100-image.MapValueIntoRange(0, 30, min, max, value), 15, image.MapValueIntoRange(5, 240, min, max, value)) // sof color
}
func (c *sc) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *sc) CreatedBy(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
