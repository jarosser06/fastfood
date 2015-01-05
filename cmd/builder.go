package cmd

import (
	"flag"
	"fmt"
	"path"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/common/fileutil"
	"github.com/jarosser06/fastfood/framework"
)

type Builder struct {
}

func (b *Builder) Run(args []string) int {
	var templatePack, cookbookPath string
	cmdFlags := flag.NewFlagSet("build", flag.ContinueOnError)
	cmdFlags.StringVar(&templatePack, "template-pack", DefaultTempPack(), "path to the template pack")
	cmdFlags.StringVar(&cookbookPath, "cookbooks-path", DefaultCookbooksPath(), "path to the cookbooks directory")
	cmdFlags.Usage = func() { fmt.Println(b.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
		return 1
	}

	manifest, err := fastfood.NewManifest(path.Join(templatePack, "manifest.json"))
	if err != nil {
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

	config, err := fastfood.NewConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fopts := fastfood.FrameworkOptions{
		Destination: cookbookPath,
		BaseFiles:   manifest.Frameworks["chef"].BaseFiles,
		BaseDirs:    manifest.Frameworks["chef"].BaseDirectories,
		Force:       false,
		Name:        config.Name,
		TemplateDir: templatePack,
	}

	c := framework.Chef{}
	if err := c.Init(fopts); err != nil {
		fmt.Printf("chef framework init returned: %v", err)
	}

	updatedFiles, err := c.GenerateBase()
	if err != nil {
		fmt.Println(err)
	}

	// Copy the template file to the cookbook if it doesn't exist in the cookbook
	err = fileutil.Copy(configFile, path.Join(cookbookPath, config.Name, "fastfood.json"))
	if err != nil {
		// if the file exists then thats fine
		if err.Error() != "file already exists" {
			fmt.Println(err)
			return 1
		}
	}

	// Generate Stencils
	for _, setConfig := range config.Stencils {
		// TODO: Validate config
		setName, ok := setConfig["name"]
		if !ok {
			fmt.Printf("missing stencilset name in config file")
			return 1
		}

		sset, err := fastfood.NewStencilSet(manifest.StencilSets[setName].Manifest)
		if err != nil {
			fmt.Printf("opening stencilset %s set returned %v", setName, err)
		}

		if _, err := sset.Valid(); err != nil {
			fmt.Printf("invalid stencilset %s %v", setName, err)
		}

		// Determine provider type
		var stencil string
		if val, ok := setConfig["stencil_set"]; ok {
			stencil = val
		} else {
			stencil = sset.DefaultStencil
		}

		if _, ok := sset.Stencils[stencil]; !ok {
			fmt.Printf("%s is not a valid stencil for stencil set %s\n", stencil, setName)
			return 1
		}

		updatedFiles, err = c.GenerateStencil(stencil, sset, setConfig)
		if err != nil {
			fmt.Println(err)
			return 1
		}

	}

	if len(updatedFiles) > 0 {
		fmt.Println("cookbook has been updated")
	}
	return 0
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
