package cmd

import (
	"flag"
	"fmt"
	"path"

	"github.com/jarosser06/fastfood"
)

type Creator struct {
	Common
}

func (c *Creator) Help() string {
	return `
Usage: fastfood new <flags> [cookbook_name]

Flags:
  -template-pack=<path>  - path to the template pack
  -cookbooks-path=<path> - path to the cookbooks directory
`
}

func (c *Creator) Run(args []string) int {
	var cookbooksPath string
	cmdFlags := flag.NewFlagSet("new", flag.ContinueOnError)
	cmdFlags.StringVar(&c.TemplatePack, "template-pack", DefaultTempPack(), "path to the template pack")
	cmdFlags.StringVar(&cookbooksPath, "cookbooks-path", DefaultCookbooksPath(), "base cookbooks directory")
	cmdFlags.Usage = func() { fmt.Println(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := c.LoadManifest(); err != nil {
		fmt.Println(err)
		return 1
	}

	remainingArgs := cmdFlags.Args()
	if len(remainingArgs) > 0 {
		cookbook := fastfood.NewCookbook(cookbooksPath, remainingArgs[0])

		//TODO: These can be collapsed into a single function
		if err := cookbook.GenDirs(c.Manifest.Cookbook.Directories); err != nil {
			fmt.Println(err)
			return 1
		}

		templatePath := path.Join(c.TemplatePack, c.Manifest.Cookbook.TemplatesPath)
		if err := cookbook.GenFiles(c.Manifest.Cookbook.Files, templatePath); err != nil {
			fmt.Println(err)
			return 1
		}

	} else {
		fmt.Println("You must enter the name of the cookbook")
		return 1
	}
	return 0
}

func (c *Creator) Synopsis() string {
	return "Creates a new empty cookbook"
}
