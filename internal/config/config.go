package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func readJSONConfig(filePath string, config *Config) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), &config)
	return err
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

type Config struct {
	Token   string   `json:"token"`
	Address string   `json:"address"`
	Display uint64   `json:"display"`
	Lights  []string `json:"lights"`
}

func Load(path string) (c Config, err error) {
	if !fileExists(path) {
		return c, errors.New("The config file does not exist.")
	}

	readJSONConfig(path, &c)
	return c, nil
}
