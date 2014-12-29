package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/common/fileutil"
)

type Builder struct {
	Common
	config    fastfood.Config
	waitGroup sync.WaitGroup
}

func (b *Builder) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("build", flag.ContinueOnError)
	cmdFlags.StringVar(&b.TemplatePack, "template-pack", DefaultTempPack(), "path to the template pack")
	cmdFlags.StringVar(&b.CookbookPath, "cookbooks-path", DefaultCookbooksPath(), "path to the cookbooks directory")
	cmdFlags.Usage = func() { fmt.Println(b.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := b.LoadManifest(); err != nil {
		fmt.Println(err)
		return 1
	}

	remainingArgs := cmdFlags.Args()
	if len(remainingArgs) == 0 {
		fmt.Println("missing argument configuration file")
		return 1
	}

	configFile := remainingArgs[0]
	if !fileutil.FileExist(configFile) {
		fmt.Printf("file does not exist %s", configFile)
		return 1
	}

	var err error
	b.config, err = fastfood.NewConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	cookbook := b.Cookbook(b.config.Name)

	// Load the providers
	providers, err := b.LoadProviders(cookbook)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// Create the cookbook directories
	if err := cookbook.GenDirs(b.Manifest.Cookbook.Directories); err != nil {
		fmt.Println(err)
		return 1
	}

	// Generate the core cookbook files
	templatePath := path.Join(b.TemplatePack, b.Manifest.Cookbook.TemplatesPath)
	if err := cookbook.GenFiles(b.Manifest.Cookbook.Files, templatePath); err != nil {
		fmt.Println(err)
		return 1
	}

	// Copy the template file to the cookbook if it doesn't exist in the cookbook
	err = fileutil.Copy(configFile, path.Join(cookbook.Path, "fastfood.json"))
	if err != nil {
		// if the file exists then thats fine
		if err.Error() != "file already exists" {
			fmt.Println(err)
			return 1
		}
	}

	// Generate provider files
	for _, provider := range b.config.Providers {
		var providerType string
		providerName := provider["provider"]
		p := providers[providerName]

		// Determine provider type
		if val, ok := provider["type"]; ok {
			providerType = val
		} else {
			if p.DefaultType == "" {
				fmt.Printf("There is no default type for provider %s\n", provider["provider"])
				return 1
			} else {
				providerType = p.DefaultType
			}
		}

		if !p.ValidType(providerType) {
			fmt.Printf("%s is not a valid type for provider %s\n", providerType, providerName)
			return 1
		}

		// Need to append the dependencies so they don't get placed twice
		// TODO: This is a hack and should be handled a bit more elegantly
		deps := p.Dependencies(providerType)
		depsAppended := cookbook.AppendDependencies(deps)
		cookbook.Dependencies = append(cookbook.Dependencies, depsAppended...)
		p.GenDirs(providerType)

		err := p.GenFiles(
			providerType,
			path.Join(b.TemplatePack, providerName),
			false,
			provider,
		)

		if err != nil {
			fmt.Println(err)
			return 1
		}

	}

	fmt.Printf("Cookbook %s updated\n", cookbook.Path)
	return 0
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
		cookbookPath = b.CookbookPath
	}

	return fastfood.NewCookbook(cookbookPath, name)
}

func (b *Builder) Synopsis() string {
	return "Creates a cookbook w/ providers from a config file"
}

func (b *Builder) Help() string {
	return `
Usage: fastfood build [config_file]

  This will create/modify a cookbook and providers.

Flags:
  -template-pack=<path>   - path to the template pack
  -cookbooks-path=<path>  - path to the cookbooks directory
`
}
