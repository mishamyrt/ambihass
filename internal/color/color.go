package color

import (
	"image"
	"math"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/mishamyrt/ambihass/internal/hass"
)

const screenshotInterval = 500

func CalcBrightness(c hass.RGBColor, minimum uint32) uint32 {
	return uint32(
		math.Max(
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

func WatchDisplayColors(ch chan<- []hass.RGBColor, display int) {
	bounds := screenshot.GetDisplayBounds(display)
	var img *image.RGBA
	var colors []hass.RGBColor
	for {
		time.Sleep(screenshotInterval * time.Millisecond)
		img, _ = screenshot.CaptureRect(bounds)
		colors = ExtractColors(img)
		ch <- colors
	}
}
