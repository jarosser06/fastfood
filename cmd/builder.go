package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/jarosser06/fastfood"
)

type Builder struct {
	Common
	config    fastfood.Config
	waitGroup sync.WaitGroup
}

func (b *Builder) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("missing argument configuration file")
	}

	configFile := args[0]
	if !fastfood.FileExist(configFile) {
		return errors.New(fmt.Sprintf("file does not exist %s", configFile))
	}

	var err error
	b.config, err = fastfood.NewConfig(configFile)
	if err != nil {
		return err
	}

	cookbook := b.Cookbook(b.config.Name)

	// Load the providers
	providers, err := b.LoadProviders(cookbook)
	if err != nil {
		return err
	}

	if err := cookbook.GenDirs(b.Manifest.Cookbook.Directories); err != nil {
		return err
	}

	templatePath := path.Join(b.templatePack, b.Manifest.Cookbook.TemplatesPath)
	if err := cookbook.GenFiles(b.Manifest.Cookbook.Files, templatePath); err != nil {
		return err
	}

	for _, provider := range b.config.Providers {
		var providerType string
		providerName := provider["provider"]
		p := providers[providerName]

		// Determine provider type
		if val, ok := provider["type"]; ok {
			providerType = val
		} else {
			if p.DefaultType == "" {
				return errors.New(fmt.Sprintf("There is no default type for provider %s", provider["provider"]))
			} else {
				providerType = p.DefaultType
			}
		}

		// Need to append the dependencies so they don't get placed twice
		// TODO: This is a hack and should be handled a bit more elegantly
		deps := p.Dependencies(providerType)
		cookbook.AppendDependencies(deps)
		cookbook.Dependencies = append(cookbook.Dependencies, deps...)
		p.GenDirs(providerType)

		err := p.GenFiles(
			providerType,
			path.Join(b.templatePack, providerName),
			false,
			provider,
		)

		if err != nil {
			return err
		}

	}

	return nil
}

// If this command is run from inside a cookbook we are going
// to assume we want to modify this cookbook
func (b *Builder) Cookbook(name string) fastfood.Cookbook {
	workingDir, _ := os.Getwd()

	if fastfood.PathIsCookbook(workingDir) {
		cookbook, _ := fastfood.NewCookbookFromPath(workingDir)

		if cookbook.Name == name {
			return cookbook
		}
	}

	var cookbookPath string
	if b.config.CookbookPath != "" {
		cookbookPath = b.config.CookbookPath
	} else {
		cookbookPath = b.cookbookPath
	}

	return fastfood.NewCookbook(cookbookPath, name)
}

func (b *Builder) Description() string {
	return "Creates a cookbook w/ providers from a config file"
}

func (b *Builder) Help() string {
	return "fastfood build [config_file]"
}
