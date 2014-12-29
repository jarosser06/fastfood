package chef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/common/fileutil"
	"github.com/jarosser06/fastfood/common/stringutil"
)

type Option struct {
	DefaultValue string `json:"default"`
	Help         string `json:"help"`
}

type StencilSet struct {
	Cookbook       Cookbook
	DefaultStencil string            `json:"default_stencil"`
	Deps           []string          `json:"dependencies"`
	Opts           map[string]Option `json:"options"`
	Stencils       map[string]struct {
		Deps        []string          `json:"dependencies"`
		Directories []string          `json:"directories"`
		Files       map[string]string `json:"files"`
		Opts        map[string]Option `json:"options"`
		Partials    []string          `json:"partials"`
	} `json:"stencils"`
}

// Return a new provider, not extremley helpful atm
// Takes a cookbook and path to the provider templates
func NewStencilSet(ckbk Cookbook) StencilSet {
	return StencilSet{
		Cookbook: ckbk,
	}
}

func NewStencilSetFromFile(ckbk Cookbook, file string) (StencilSet, error) {
	sset := NewStencilSet(ckbk)

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return sset, fmt.Errorf("Failed to read file %s: %v", file, err)
	}

	err = json.Unmarshal(f, &sset)
	if err != nil {
		return sset, fmt.Errorf("Failed to unmarshal provider json: %v", err)
	}

	return sset, nil
}

// Return true if the type exists in types
func (s *StencilSet) ValidType(stencil string) bool {
	_, ok := s.Stencils[stencil]
	return ok
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

// Creates the expected struct for all templates and renders each template one by one
func (s *StencilSet) GenFiles(stencil string, templatesPath string, forceWrite bool, opts map[string]string) error {
	mergedOpts := s.MergeOpts(stencil, opts)
	cmap := make(map[string]string)
	for key, val := range mergedOpts {
		cmap[stringutil.CapitalizeString(key)] = val
	}

	templateOpts := struct {
		*fastfood.Helpers
		Cookbook Cookbook
		Options  map[string]string
	}{
		Cookbook: s.Cookbook,
		Options:  cmap,
	}

	// TODO: Some of this could be cleaned up and added to the provider.Template
	files := s.Stencils[stencil].Files
	partials := s.Stencils[stencil].Partials
	for cookbookFile, templateFile := range files {
		cookbookFile = strings.Replace(cookbookFile, "<NAME>", templateOpts.Options["Name"], 1)
		if fileutil.FileExist(path.Join(s.Cookbook.Path, cookbookFile)) && !forceWrite {
			continue
		}
		var content []string
		b, err := ioutil.ReadFile(path.Join(templatesPath, templateFile))
		if err != nil {
			return fmt.Errorf("provider.GenFiles() reading file returned %v", err)
		}

		content = append(content, string(b))
		for _, partial := range partials {
			b, err := ioutil.ReadFile(path.Join(templatesPath, partial))
			if err != nil {
				return fmt.Errorf("provider.GenFiles() reading file returned %v", err)
			}

			content = append(content, string(b))
		}

		t, err := fastfood.NewTemplate(cookbookFile, templateOpts, content)

		if err != nil {
			return fmt.Errorf("Error creating template: %v", err)
		}

		t.CleanNewlines()

		if err := t.Flush(path.Join(s.Cookbook.Path, cookbookFile)); err != nil {
			return fmt.Errorf("Error writing file: %v", err)
		}
	}
	return nil
}

// Generate any directories needed for the provider
// Always make sure recipes and test/unit/spec are created
// since they are the most common
func (s *StencilSet) GenDirs(stencil string) error {
	dirs := append(
		s.Stencils[stencil].Directories,
		"recipes",
		"test/unit/spec",
	)

	for _, dir := range dirs {
		fullPath := path.Join(s.Cookbook.Path, dir)

		if !fileutil.FileExist(fullPath) {
			err := os.MkdirAll(path.Join(s.Cookbook.Path, dir), 0755)

			if err != nil {
				return fmt.Errorf("database.GenDirs(): %v", err)
			}
		}
	}

	return nil
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
Default Type: %s

Global Options:

%s

Provider Types:

%s
`, s.DefaultStencil, strings.Join(globalOpts, "\n"), strings.Join(stencils, "\n"))
	return helpText
}
