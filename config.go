package fastfood

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Config struct {
	CookbookPath string              `json:"cookbook_path,omitempty"`
	Name         string              `json:"name,omitempty"`
	Providers    []map[string]string `json:"providers,omitempty"`
	Target       []string            `json:"target,omitempty"`
}

func NewConfig(path string) (Config, error) {
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
