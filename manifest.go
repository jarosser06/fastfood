package fastfood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

const templatePackAPI = 1

type Manifest struct {
	API         int `json:"api"`
	Path        string
	StencilSets map[string]struct {
		Name          string
		Manifest      string
		Help          string `json:"help"`
		templatesPath string
	} `json:"stencil_sets"`

	Frameworks map[string]struct {
		Name            string
		BaseFiles       []string `json:"files"`
		BaseDirectories []string `json:"directories"`
	}
}

func NewManifest(mpath string) (Manifest, error) {

	var m Manifest

	f, err := ioutil.ReadFile(mpath)
	if err != nil {
		return m, fmt.Errorf("reading manifest %s: %v", mpath, err)
	}

	err = json.Unmarshal(f, &m)
	if err != nil {
		return m, fmt.Errorf("parsing manifest %s: %v", mpath, err)
	}

	// Set helpful values for stencilsets
	for n, _ := range m.StencilSets {
		tmp := m.StencilSets[n]
		tmp.Name = n
		tmp.Manifest = path.Join(
			filepath.Dir(mpath),
			"stencils",
			n,
			"manifest.json",
		)
		m.StencilSets[n] = tmp
	}

	// Set name for frameworks
	for n, _ := range m.Frameworks {
		tmp := m.Frameworks[n]
		tmp.Name = n
		m.Frameworks[n] = tmp
	}

	return m, nil
}

func (m *Manifest) Help() string {
	var shelp []string

	for n, s := range m.StencilSets {
		var help string
		if s.Help == "" {
			help = "NO HELP FOUND"
		} else {
			help = s.Help
		}

		shelp = append(
			shelp,
			fmt.Sprintf("  %-15s - %s", n, help),
		)
	}

	return fmt.Sprintf(`
Available Stencil Sets:

%s
`, strings.Join(shelp, "\n"))
}

func (m *Manifest) Valid() (bool, error) {
	if m.API != templatePackAPI {
		return false, fmt.Errorf("template pack api version %d not compatible with %d", m.API, templatePackAPI)
	}

	return true, nil
}
