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
	bounds := screenshot.GetDisplayBounds(configuration.Display)
	session := hass.NewSession(configuration.Address, configuration.Token)
	var img *image.RGBA
	var colors []hass.RGBColor
	var current hass.RGBColor
	var light config.Light

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
		if color.IsCloseColors(colors[0], current, deadZone) {
			debug("Skip, too close colors")
			continue
		}
		current = colors[0]
		for _, light = range configuration.Lights {
			debug("Update " + light.ID)
			session.TurnOn(hass.LightService{
				Entity:     light.ID,
				Color:      current,
				Brightness: color.CalcBrightness(current, light.MinBrightness),
			})
		}
	}
}
