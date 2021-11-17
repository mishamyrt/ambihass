package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/kbinani/screenshot"
	color "github.com/mishamyrt/ambihass/internal/color"
	"github.com/mishamyrt/ambihass/internal/config"
	"github.com/mishamyrt/ambihass/internal/hass"
	"github.com/mishamyrt/ambihass/internal/lights"
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
	colorChan := make(chan []hass.RGBColor)
	session := hass.NewSession(configuration.Address, configuration.Token)
	controller := lights.Controller{
		Session: session,
		Devices: configuration.Lights,
	}
	fmt.Println(
		"Initiated ambilight for display " +
			fmt.Sprint(configuration.Display) +
			" on " +
			configuration.Address,
	)
	go watchColors(colorChan, configuration.Display)
	for {
		select {
		case colors, _ := <-colorChan:
			controller.SetColor(colors[0])
		}
	}
}

func watchColors(ch chan<- []hass.RGBColor, display int) {
	bounds := screenshot.GetDisplayBounds(display)
	var img *image.RGBA
	var colors []hass.RGBColor
	for {
		time.Sleep(500 * time.Millisecond)
		img, _ = screenshot.CaptureRect(bounds)
		colors = color.ExtractColors(img)
		ch <- colors
	}
}
