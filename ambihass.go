package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/kbinani/screenshot"
	color "github.com/mishamyrt/ambihass/internal/color_extractor"
	"github.com/mishamyrt/ambihass/internal/config"
	"github.com/mishamyrt/ambihass/internal/hass"
	"github.com/spf13/pflag"
)

const deadZone = 5

var configPath string = os.Getenv("HOME") + ".config/ambihass/config.json"
var debugMode bool

func init() {
	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.BoolVarP(&debugMode, "debug", "d", false, "Debug mode. Prints mode information")
	pflag.Parse()
}

func debug(a ...interface{}) {
	if debugMode {
		fmt.Println(a...)
	}
}

func main() {
	configuration, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}
	bounds := screenshot.GetDisplayBounds(1)
	session := hass.NewSession(configuration.Address, configuration.Token)
	var img *image.RGBA
	var colors []hass.RGBColor
	var last hass.RGBColor

	fmt.Println(
		"Initiated ambilight for display " +
			fmt.Sprint(configuration.Display) +
			" on " +
			configuration.Address,
	)

	for {
		time.Sleep(500 * time.Millisecond)
		img, err = screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		colors = color.ExtractColors(img)
		debug("Colors:", colors)
		if color.IsCloseColors(colors[0], last, deadZone) {
			debug("Skip, too close colors")
			continue
		}
		debug("Update")
		session.TurnOn(hass.LightService{
			Entity:     configuration.Lights[0],
			Color:      colors[0],
			Brightness: color.CalcBrightness(colors[0], 140),
		})
		last = colors[0]
	}
}
