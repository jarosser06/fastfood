package fastfood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Name      string              `json:"name,omitempty"`
	Framework string              `json:"framework,omitempty"`
	Stencils  []map[string]string `json:"stencils,omitempty"`
	Target    []string            `json:"target,omitempty"`
}

func NewConfig(path string) (Config, error) {
	var c Config

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("Error reading config %v", err)
	}

	err = json.Unmarshal(file, &c)

	if err != nil {
		return c, fmt.Errorf("Error parsing json file %v", err)
	}

	return c, nil
}
