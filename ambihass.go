package main

import (
	"fmt"
	"os"

	"github.com/mishamyrt/ambihass/internal/color"
	"github.com/mishamyrt/ambihass/internal/config"
	"github.com/mishamyrt/ambihass/internal/hass"
	"github.com/mishamyrt/ambihass/internal/lights"
	"github.com/mishamyrt/ambihass/internal/log"
	"github.com/spf13/pflag"
)

const schedulerInterval = 100

var configPath string = os.Getenv("HOME") + ".config/ambihass/config.json"
var debugMode bool

func init() {
	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.BoolVarP(&debugMode, "debug", "d", false, "Debug mode. Prints mode information")
	pflag.Parse()
}

func main() {
	log.DebugMode = debugMode
	log.Debug("Debug mode")
	configuration, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}
	session := hass.NewSession(configuration.Address, configuration.Token)
	err = session.CheckAPI()
	if err != nil {
		panic("Failed to connect to Home Assistant")
	}
	controller := lights.Controller{
		Session: session,
		Devices: configuration.Lights,
	}
	log.Message(
		"Initiated ambilight for display " +
			fmt.Sprint(configuration.Display) +
			" on " +
			configuration.Address,
	)
	colorChan := make(chan []hass.RGBColor)
	go color.WatchDisplayColors(colorChan, configuration.Display)
	controller.Start(schedulerInterval, colorChan)
}
