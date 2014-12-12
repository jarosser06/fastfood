package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	tempPackEnvVar    = "FASTFOOD_TEMPLATE_PACK"
	cookbookTemplates = "cookbook"
)

type Manifest struct {
	Providers map[string]struct {
		Name          string `json:"name"`
		Manifest      string `json:"manifest"`
		Help          string `json:"help"`
		templatesPath string
	}

	Cookbook struct {
		Directories   []string `json:"directories"`
		Files         []string `json:"files"`
		TemplatesPath string   `json:"templates_path"`
	}
}

func NewManifest(path string) Manifest {

	var manifest Manifest

	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file %s", path))
	}

	err = json.Unmarshal(f, &manifest)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse json: %v", err))
	}

	if manifest.Cookbook.TemplatesPath == "" {
		manifest.Cookbook.TemplatesPath = cookbookTemplates
	}

	return manifest

}

func DefaultTempPack() string {
	packEnv := os.Getenv(tempPackEnvVar)
	if packEnv == "" {
		return path.Join(os.Getenv("HOME"), "fastfood")
	} else {
		return packEnv
	}
}
