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
	TopNHeaderFGDanger(*gg.Context)
	TopNHeaderBG(*gg.Context)
	TopNHeaderOutline(*gg.Context)
	TopNCellDarkerBG(*gg.Context)
	TopNCellLighterBG(*gg.Context)
	TopNCellOutline(*gg.Context)
	TopNCellPosition(*gg.Context)
	TopNCellDriver(*gg.Context)
	TopNCellValue(*gg.Context)
	TopNCellValueDanger(*gg.Context)
	HeatmapHeaderFG(*gg.Context)
	HeatmapHeaderDarkerBG(*gg.Context)
	HeatmapHeaderLighterBG(*gg.Context)
	HeatmapTimeslotFG(*gg.Context)
	HeatmapTimeslotBG(*gg.Context)
	HeatmapTimeslotZero(*gg.Context)
	HeatmapTimeslotMapping(*gg.Context, int, int, int)
	LastUpdate(*gg.Context)
	CreatedBy(*gg.Context)
}

func Get(scheme string) Colorizer {
	var c Colorizer
	switch scheme {
	case "blue":
		c = NewBlueScheme()
	case "green":
		c = NewGreenScheme()
	case "yellow":
		c = NewYellowScheme()
	case "red":
		c = NewRedScheme()
	case "black":
		c = NewBlackScheme()
	case "simucube":
		c = NewSimuCubeScheme()
	case "apex":
		c = NewApexScheme()
	case "indypro":
		c = NewPMScheme()
	default:
		c = NewPMScheme()
	}
	return c
}
