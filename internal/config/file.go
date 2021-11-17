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
