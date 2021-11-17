package color_extractor

import (
	"math"

	"github.com/mishamyrt/ambihass/internal/hass"
)

func CalcBrightness(c hass.RGBColor, minimum uint32) uint32 {
	return uint32(
		math.Min(
			math.Max(float64(c[0]), math.Max(float64(c[1]), float64(c[2]))),
			float64(minimum),
		),
	)
}

func IsCloseColors(first, second hass.RGBColor, distance uint32) bool {
	for i := 0; i < 3; i++ {
		if first[i] < second[i]-distance || first[i] > second[i]+distance {
			return false
		}
	}
	return true
}
