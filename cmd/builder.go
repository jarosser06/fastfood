package cmd

import (
	"flag"
	"fmt"
	"path"

	"github.com/jarosser06/fastfood"
)

type Builder struct {
	CookbooksPath string
}

func (b *Builder) Help() string {
	return "TODO: Add help text for new command"
}

func (b *Builder) Run(args []string) int {
	var config, cookbookPath, templatePack string
	cmdFlags := flag.NewFlagSet("builder", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Println(b.Help()) }
	cmdFlags.StringVar(&cookbookPath, "cookbooks-dir", "", "directory to store new cookbook")
	cmdFlags.StringVar(&config, "config", DefaultCookbookPath(), "json config to build from NOT IMPLEMENTED")
	cmdFlags.StringVar(&templatePack, "template-pack", DefaultTempPack(), "path to the template pack")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	manifest, err := fastfood.NewManifest(path.Join(templatePack, "manifest.json"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	args = cmdFlags.Args()

	if len(config) > 0 {
		fmt.Println("Placeholder for generating from a config")
	} else {
		if len(args) > 0 {
			cookbook := fastfood.NewCookbook(cookbookPath, args[0])

			//TODO: These can be collapsed into a single function
			if err := cookbook.GenDirs(manifest.Cookbook.Directories); err != nil {
				fmt.Println(err)
				return 1
			}

			templatePath := path.Join(templatePack, manifest.Cookbook.TemplatesPath)
			if err := cookbook.GenFiles(manifest.Cookbook.Files, templatePath); err != nil {
				fmt.Println(err)
				return 1
			}

		} else {
			fmt.Println("You must enter the name of the cookbook")
			return 1
		}
	}
	return 0
}

func (b *Builder) Synopsis() string {
	return "Builds a new cookbook"
}
