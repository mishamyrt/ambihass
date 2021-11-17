package config

import (
	"errors"
)

type Light struct {
	ID            string `json:"id"`
	MinBrightness uint32
}

// "id": "light.screen_back_middle",
//             "minBrightness": 60

type Config struct {
	Token   string  `json:"token"`
	Address string  `json:"address"`
	Display int     `json:"display"`
	Lights  []Light `json:"lights"`
}

func Load(path string) (c Config, err error) {
	if !fileExists(path) {
		return c, errors.New("The config file does not exist.")
	}

	readJSONConfig(path, &c)
	return c, nil
}
