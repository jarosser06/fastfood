package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Applications []map[string]string `json:"applications,omitempty"`
	CookbookPath string              `json:"cookbook_path,omitempty"`
	Databases    []map[string]string `json:"database,omitempty"`
	Name         string              `json:"name,omitempty"`
	Queues       []map[string]string `json:"queues,omitempty"`
	Target       []string            `json:"target,omitempty"`
}

func Parse(path string) (Config, error) {
	var c Config

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return c, errors.New(fmt.Sprintf("Error reading config %v", err))
	}

	err = json.Unmarshal(file, &c)

	if err != nil {
		return c, errors.New(fmt.Sprintf("Error parsing json file %v", err))
	}

	return c, nil
}
