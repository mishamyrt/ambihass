package config

import (
	"errors"

	"github.com/mishamyrt/ambihass/internal/hass"
)

type Config struct {
	Token   string             `json:"token"`
	Address string             `json:"address"`
	Display int                `json:"display"`
	Lights  []hass.LightDevice `json:"lights"`
}

func Load(path string) (c Config, err error) {
	if !fileExists(path) {
		return c, errors.New("The config file does not exist.")
	}

	readJSONConfig(path, &c)
	return c, nil
}
