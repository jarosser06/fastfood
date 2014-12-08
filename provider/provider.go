package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood/util"
)

type ProviderType struct {
	Deps        []string           `json:"dependencies"`
	Directories []string           `json:"directories"`
	Files       map[string]string  `json:"files"`
	Opts        map[string]Options `json:"options"`
	Partials    []string           `json:"partials"`
}

type Options struct {
	DefaultValue string `json:"default"`
	Help         string `json:"help"`
}

type Provider struct {
	Cookbook    Cookbook
	DefaultType string                  `json:"default_type"`
	Deps        []string                `json:"dependencies"`
	Opts        map[string]Options      `json:"options"`
	Types       map[string]ProviderType `json:"types"`
}

// Return a new provider, not extremley helpful atm
// Takes a cookbook and path to the provider templates
func NewProvider(ckbk Cookbook) Provider {
	return Provider{
		Cookbook: ckbk,
	}
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
func (p *Provider) GenFiles(typeName string, templatesPath string, opts map[string]string) error {
	templateOpts := struct {
		*Helpers
		Cookbook Cookbook
		Options  map[string]string
	}{
		Cookbook: p.Cookbook,
		Options:  p.MergeOpts(typeName, opts),
	}

	fmt.Printf("%v\n", templateOpts.Options)

	// TODO: Some of this could be cleaned up and added to the provider.Template
	files := p.Types[typeName].Files
	partials := p.Types[typeName].Partials
	for cookbookFile, templateFile := range files {
		cookbookFile = strings.Replace(cookbookFile, "<NAME>", templateOpts.Options["name"], 1)
		if util.FileExist(path.Join(p.Cookbook.Path, cookbookFile)) {
			continue
		}
		var content []string
		b, err := ioutil.ReadFile(path.Join(templatesPath, templateFile))
		if err != nil {
			return errors.New(fmt.Sprintf("provider.GenFiles() reading file returned %v", err))
		}

		content = append(content, string(b))
		for _, partial := range partials {
			b, err := ioutil.ReadFile(path.Join(templatesPath, partial))
			if err != nil {
				return errors.New(fmt.Sprintf("provider.GenFiles() reading file returned %v", err))
			}

			content = append(content, string(b))
		}

		t, err := NewTemplate(cookbookFile, templateOpts, content)

		if err != nil {
			return errors.New(fmt.Sprintf("Error creating template: %v", err))
		}

		t.CleanNewlines()

		if err := t.Flush(path.Join(p.Cookbook.Path, cookbookFile)); err != nil {
			return errors.New(fmt.Sprintf("Error writing file: %v", err))
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

		if !util.FileExist(fullPath) {
			err := os.MkdirAll(path.Join(p.Cookbook.Path, dir), 0755)

			if err != nil {
				return errors.New(fmt.Sprintf("database.GenDirs(): %v", err))
			}
		}
	}

	return nil
}
