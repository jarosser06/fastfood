package fastfood

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood/common/fileutil"
	"github.com/jarosser06/fastfood/common/stringutil"
)

type Option struct {
	DefaultValue string `json:"default"`
	Help         string `json:"help"`
}

type Provider struct {
	Cookbook    Cookbook
	DefaultType string            `json:"default_type"`
	Deps        []string          `json:"dependencies"`
	Opts        map[string]Option `json:"options"`
	Types       map[string]struct {
		Deps        []string          `json:"dependencies"`
		Directories []string          `json:"directories"`
		Files       map[string]string `json:"files"`
		Opts        map[string]Option `json:"options"`
		Partials    []string          `json:"partials"`
	}
}

// Return a new provider, not extremley helpful atm
// Takes a cookbook and path to the provider templates
func NewProvider(ckbk Cookbook) Provider {
	return Provider{
		Cookbook: ckbk,
	}
}

//TODO: Return proper errors instead fo panicing
func NewProviderFromFile(ckbk Cookbook, file string) (Provider, error) {
	provider := NewProvider(ckbk)

	f, err := ioutil.ReadFile(file)
	// Probably shouldn't Panic, might scare people
	if err != nil {
		return provider, fmt.Errorf("Failed to read file %s: %v", file, err)
	}

	err = json.Unmarshal(f, &provider)
	if err != nil {
		return provider, fmt.Errorf("Failed to unmarshal provider json: %v", err)
	}

	return provider, nil
}

// Return true if the type exists in types
func (p *Provider) ValidType(typeName string) bool {
	_, ok := p.Types[typeName]
	return ok
}

// Returns a string slice of dependencies
func (p *Provider) Dependencies(typeName string) []string {
	deps := p.Deps
	deps = append(p.Types[typeName].Deps, deps...)

	return deps
}

// Merge all of the options from a given map with the defaults from the
// type and provider
func (p *Provider) MergeOpts(typeName string, opts map[string]string) map[string]string {

	// Merge type options first
	// Gives the ability to override provider global options
	for optName, optVal := range p.Types[typeName].Opts {
		if _, ok := opts[optName]; !ok {
			if val := optVal.DefaultValue; val != "" {
				opts[optName] = optVal.DefaultValue
			} else {
				opts[optName] = ""
			}
		}
	}

	for optName, optVal := range p.Opts {
		if _, ok := opts[optName]; !ok {
			if val := optVal.DefaultValue; val != "" {
				opts[optName] = optVal.DefaultValue
			} else {
				opts[optName] = ""
			}
		}
	}

	return opts
}

// Creates the expected struct for all templates and renders each template one by one
func (p *Provider) GenFiles(typeName string, templatesPath string, forceWrite bool, opts map[string]string) error {
	mergedOpts := p.MergeOpts(typeName, opts)
	cappedMap := make(map[string]string)
	for key, val := range mergedOpts {
		cappedMap[stringutil.CapitalizeString(key)] = val
	}

	templateOpts := struct {
		*Helpers
		Cookbook Cookbook
		Options  map[string]string
	}{
		Cookbook: p.Cookbook,
		Options:  cappedMap,
	}

	// TODO: Some of this could be cleaned up and added to the provider.Template
	files := p.Types[typeName].Files
	partials := p.Types[typeName].Partials
	for cookbookFile, templateFile := range files {
		cookbookFile = strings.Replace(cookbookFile, "<NAME>", templateOpts.Options["Name"], 1)
		if fileutil.FileExist(path.Join(p.Cookbook.Path, cookbookFile)) && !forceWrite {
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

		t, err := NewTemplate(cookbookFile, templateOpts, content)

		if err != nil {
			return fmt.Errorf("Error creating template: %v", err)
		}

		t.CleanNewlines()

		if err := t.Flush(path.Join(p.Cookbook.Path, cookbookFile)); err != nil {
			return fmt.Errorf("Error writing file: %v", err)
		}
	}
	return nil
}

// Generate any directories needed for the provider
// Always make sure recipes and test/unit/spec are created
// since they are the most common
func (p *Provider) GenDirs(typeName string) error {
	dirs := append(
		p.Types[typeName].Directories,
		"recipes",
		"test/unit/spec",
	)

	for _, dir := range dirs {
		fullPath := path.Join(p.Cookbook.Path, dir)

		if !fileutil.FileExist(fullPath) {
			err := os.MkdirAll(path.Join(p.Cookbook.Path, dir), 0755)

			if err != nil {
				return fmt.Errorf("database.GenDirs(): %v", err)
			}
		}
	}

	return nil
}

// Print Provider help
func (p *Provider) Help() string {
	var globalOpts, providerTypes []string
	for name, opt := range p.Opts {
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

	for providerType := range p.Types {
		providerTypes = append(
			providerTypes,
			fmt.Sprintf("  %s", providerType),
		)
	}
	helpText := fmt.Sprintf(`
Default Type: %s

Global Options:

%s

Provider Types:

%s
`, p.DefaultType, strings.Join(globalOpts, "\n"), strings.Join(providerTypes, "\n"))
	return helpText
}
