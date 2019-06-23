package color

import "github.com/fogleman/gg"

type Colorizer interface {
	Border(*gg.Context)
	Background(*gg.Context)
	Transparent(*gg.Context)
	HeaderFG(*gg.Context)
	HeaderLeftBG(*gg.Context)
	HeaderRightBG(*gg.Context)
	TopNHeaderFG(*gg.Context)
	TopNHeaderBG(*gg.Context)
	TopNHeaderOutline(*gg.Context)
	TopNCellDarkerBG(*gg.Context)
	TopNCellLighterBG(*gg.Context)
	TopNCellOutline(*gg.Context)
	TopNCellPosition(*gg.Context)
	TopNCellDriver(*gg.Context)
	TopNCellValue(*gg.Context)
	HeatmapHeaderFG(*gg.Context)
	HeatmapHeaderDarkerBG(*gg.Context)
	HeatmapHeaderLighterBG(*gg.Context)
	HeatmapTimeslotFG(*gg.Context)
	HeatmapTimeslotBG(*gg.Context)
	HeatmapTimeslotMapping(*gg.Context, int, int, int)
	LastUpdate(*gg.Context)
}
