/*
 * This handles everything to do with stencilsets and stencils for fastfood
 */
package fastfood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/jarosser06/fastfood/common/fileutil"
)

const stencilAPI = 1

type Option struct {
	DefaultValue string `json:"default"`
	Help         string `json:"help"`
}

type StencilSet struct {
	Name           string                       `json:"id"`
	APIVersion     int                          `json:"api"`
	BerksDeps      map[string]map[string]string `json:"berks_dependencies"`
	DefaultStencil string                       `json:"default_stencil"`
	Deps           []string                     `json:"dependencies"`
	Opts           map[string]Option            `json:"options"`
	Stencils       map[string]struct {
		BerksDeps   map[string]map[string]string `json:"berks_dependencies"`
		Deps        []string                     `json:"dependencies"`
		Directories []string                     `json:"directories"`
		Files       map[string]string            `json:"files"`
		Opts        map[string]Option            `json:"options"`
		Partials    []string                     `json:"partials"`
	} `json:"stencils"`
}

// Return a new stencil set and error
func NewStencilSet(file string) (StencilSet, error) {
	sset := StencilSet{}

	// Verify the file exists before proceeding
	if !fileutil.FileExist(file) {
		return sset, fmt.Errorf("file %s does not exist", file)
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return sset, fmt.Errorf("reading file %s returned error %v", file, err)
	}

	// TODO: replace with override unmarshal that provides better errors
	err = json.Unmarshal(f, &sset)
	if err != nil {
		return sset, fmt.Errorf("unmarshalling stencil set %s return error %v", file, err)
	}

	// Calculate the actual paths to the templates
	sset.calculatePaths(filepath.Dir(file))
	return sset, nil
}

// Calculates full paths for stencil templates and partials
// THIS IS CALLED BY THE NewStencilSet() function
func (s *StencilSet) calculatePaths(basePath string) {
	for name, stencil := range s.Stencils {
		for file, temp := range stencil.Files {
			stencil.Files[file] = path.Join(basePath, temp)
		}

		for i, partial := range stencil.Partials {
			stencil.Partials[i] = path.Join(basePath, partial)
		}

		s.Stencils[name] = stencil
	}

}

// Return true if the type exists in types
func (s *StencilSet) Valid() (bool, error) {
	// Check if stencil version matches api version
	if s.APIVersion != stencilAPI {
		return false, fmt.Errorf(
			"api version mismatch, version %d not compatible with %d",
			s.APIVersion,
			stencilAPI,
		)
	}

	// Verify default stencil is set
	if s.DefaultStencil == "" {
		return false, fmt.Errorf("must have a default stencil defined")
	}

	// Verify default stencil is a valid stencil
	_, ok := s.Stencils[s.DefaultStencil]
	if ok {
		return true, nil
	} else {
		return false, fmt.Errorf("default stencil %s is not a valid stencil")
	}
}

// Returns a string slice of dependencies
func (s *StencilSet) Dependencies(stencil string) []string {
	deps := s.Deps
	deps = append(s.Stencils[stencil].Deps, deps...)

	return deps
}

// Merge all of the options from a given map with the defaults from the
// type and provider
func (s *StencilSet) MergeOpts(stencil string, opts map[string]string) map[string]string {

	// Merge type options first
	// Gives the ability to override provider global options
	for name, val := range s.Stencils[stencil].Opts {
		if _, ok := opts[name]; !ok {
			if v := val.DefaultValue; v != "" {
				opts[name] = val.DefaultValue
			} else {
				opts[name] = ""
			}
		}
	}

	for name, val := range s.Opts {
		if _, ok := opts[name]; !ok {
			if v := val.DefaultValue; v != "" {
				opts[name] = val.DefaultValue
			} else {
				opts[name] = ""
			}
		}
	}

	return opts
}

// Print Provider help
func (s *StencilSet) Help() string {
	var globalOpts, stencils []string
	for name, opt := range s.Opts {
		var help string
		if opt.Help == "" {
			help = "NO HELP FOUND"
		} else {
			help = opt.Help
		}

		globalOpts = append(
			globalOpts,
			fmt.Sprintf("  %-15s - %s", name, help),
		)
	}

	for stencil := range s.Stencils {
		stencils = append(
			stencils,
			fmt.Sprintf("  %s", stencil),
		)
	}
	helpText := fmt.Sprintf(`
Default Stencil: %s

Global Options:

%s

Provider Types:

%s
`, s.DefaultStencil, strings.Join(globalOpts, "\n"), strings.Join(stencils, "\n"))
	return helpText
}
