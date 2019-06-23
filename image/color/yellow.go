package color

package color

import "github.com/fogleman/gg"

type yellow struct{}

func NewYellowScheme() Colorizer {
	return &yellow{}
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
	dc.SetRGB255(61, 133, 198) // dark blue 1
	dc.SetRGB255(11, 83, 148) // dark blue 2
	dc.SetRGB255(7, 55, 99) // dark blue 3
*/
func (c *yellow) Border(dc *gg.Context) {
	dc.SetRGB255(39, 39, 39) // dark gray 1
}
func (c *yellow) Background(dc *gg.Context) {
	dc.SetRGB255(243, 243, 243) // light gray 3
}
-
func (c *yellow) Transparent(dc *gg.Context) {
	dc.SetRGBA255(0, 0, 0, 0) // transparent
}
func (c *yellow) HeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *yellow) HeaderLeftBG(dc *gg.Context) {
	dc.SetRGB255(7, 55, 99) // dark blue 3
}
func (c *yellow) HeaderRightBG(dc *gg.Context) {
	dc.SetRGB255(11, 83, 148) // dark blue 2
}
func (c *yellow) TopNHeaderFG(dc *gg.Context) {
	dc.SetRGB255(255, 255, 255) // white
}
func (c *yellow) TopNHeaderBG(dc *gg.Context) {
	dc.SetRGB255(133, 133, 133) // gray 1
}
func (c *yellow) TopNHeaderOutline(dc *gg.Context) {
	dc.SetRGB255(55, 55, 55) // dark gray 2
}
func (c *yellow) TopNCellDarkerBG(dc *gg.Context) {
	dc.SetRGB255(225, 225, 225) // light gray 1.5
}
func (c *yellow) TopNCellLighterBG(dc *gg.Context) {
	dc.SetRGB255(241, 241, 241) // light gray 2.5
}
func (c *yellow) TopNCellOutline(dc *gg.Context) {
	dc.SetRGB255(155, 155, 155) // gray 2
}
func (c *yellow) TopNCellPosition(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *yellow) TopNCellDriver(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
func (c *yellow) TopNCellValue(dc *gg.Context) {
	dc.SetRGB255(7, 55, 99) // dark blue 3
}
func (c *yellow) LastUpdate(dc *gg.Context) {
	dc.SetRGB255(0, 0, 0) // black
}
