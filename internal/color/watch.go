package color

import (
	"image"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/mishamyrt/ambihass/internal/hass"
)

func WatchDisplayColors(ch chan<- []hass.RGBColor, display int) {
	bounds := screenshot.GetDisplayBounds(display)
	var img *image.RGBA
	for {
		time.Sleep(screenshotInterval * time.Millisecond)
		img, _ = screenshot.CaptureRect(bounds)
		ch <- ExtractColors(img)
	}
}
