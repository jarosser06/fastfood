package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood"
)

const (
	EnvTempPack      = "FASTFOOD_TEMPLATE_PACK"
	EnvCookbooksPath = "COOKBOOKS"
)

// Templatepack is the path to the templatepack
// CookbookPath is the path to cookbooks
// Manifest is loaded using LoadManifest
type Common struct {
	TemplatePack string
	CookbookPath string
	Force        bool
	Manifest     fastfood.Manifest
}

type ProviderMap map[string]fastfood.Provider

// Load the core manifest so we can provide
// dynamic help options
func (c *Common) LoadManifest() error {
	baseManifest := path.Join(c.TemplatePack, "manifest.json")
	if !fastfood.FileExist(baseManifest) {
		return errors.New(fmt.Sprintf("Error no such file %s\n", baseManifest))
	}

	var err error
	c.Manifest, err = fastfood.NewManifest(baseManifest)
	if err != nil {
		return err
	}

	return nil
}

//Load all providers into memory
func (c *Common) LoadProviders(cookbook fastfood.Cookbook) (ProviderMap, error) {
	providerMap := make(ProviderMap)

	for name, provider := range c.Manifest.Providers {
		p, err := fastfood.NewProviderFromFile(
			cookbook,
			path.Join(c.TemplatePack, provider.Manifest),
		)

		if err != nil {
			return providerMap, errors.New(
				fmt.Sprintf("error loading provider from manifest %v", err),
			)
		}

		providerMap[name] = p
	}

	return providerMap, nil
}

// Translates key:value strings into a map
func MapArgs(args []string) map[string]string {
	var argMap map[string]string
	argMap = make(map[string]string)

	for _, arg := range args {
		if strings.Contains(arg, ":") {
			// Split at the first : in an arg
			splitArg := strings.SplitN(arg, ":", 2)

			argMap[splitArg[0]] = splitArg[1]
		}
	}

	return argMap
}

func DefaultTempPack() string {
	packEnv := os.Getenv(EnvTempPack)
	if packEnv == "" {
		return path.Join(os.Getenv("HOME"), "fastfood")
	} else {
		return packEnv
	}
}
