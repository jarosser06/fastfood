package provider

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/jarosser06/fastfood/util"
)

type ProviderType struct {
	dependencies []string          `json:"dependencies"`
	files        map[string]string `json:"files"`
	partials     []string          `json:"partials"`
	directories  []string          `json:"directories"`
}

type Options struct {
	DefaultValue string `json:"default"`
	Help         string `json:"help"`
}

type Provider struct {
	*Helpers
	Cookbook     Cookbook
	defaultType  string                  `json:"default_type"`
	dependencies []string                `json:"dependencies"`
	options      map[string]Options      `json:"options"`
	types        map[string]ProviderType `json:"types"`
}

// Return a new provider, not extremley helpful atm
func NewProvider(ckbk Cookbook) Provider {
	return Provider{
		Cookbook: ckbk,
	}
}

// Return true if the type exists in types
func (p *Provider) ValidType(typeName string) bool {
	_, ok := p.types[typeName]
	return ok
}

// Returns a string slice of dependencies
func (p *Provider) Dependencies(typeName string) []string {
	deps := p.dependencies
	deps = append(p.types[typeName].dependencies, deps...)

	return deps
}

// TODO: Rewrite without Rice
func (p *Provider) GenFiles(typeName string) error {
	/*
		if util.FileExist(path.Join(p.Cookbook.Path, recipeFile)) {
			return errors.New(fmt.Sprintf("%s already exists", recipeFile))
		}

		cookbookFiles := map[string]string{
			recipeFile: fmt.Sprintf("recipes/%s_%s.rb", d.Type, d.Role),
			specFile:   fmt.Sprintf("test/unit/spec/%s_%s_spec.rb", d.Type, d.Role),
		}

		templateBox, _ := rice.FindBox("../templates/database")
		for cookbookFile, templateFile := range cookbookFiles {
			tmpStr, _ := templateBox.String(templateFile)

			t, err := provider.NewTemplate(cookbookFile, d, []string{tmpStr})

			if err != nil {
				return errors.New(fmt.Sprintf("Error creating template: %v", err))
			}

			t.CleanNewlines()

			if err := t.Flush(path.Join(d.Cookbook.Path, cookbookFile)); err != nil {
				return errors.New(fmt.Sprintf("Error writing file: %v", err))
			}
		}
	*/
	return nil
}

// Generate any directories needed for the provider
// Always make sure recipes and test/unit/spec are created
// since they are the most common
func (p *Provider) GenDirs(typeName string) error {
	dirs := append(
		p.types[typeName].directories,
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
