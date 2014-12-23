package cmd

import (
	"flag"
	"fmt"
	"path"

	"github.com/jarosser06/fastfood"
)

type Creator struct {
	CookbooksPath string
}

func (c *Creator) Help() string {
	return `
Usage:
  fastfood new <flags> [cookbook_name]

  Simple commmand that generates a new cookbook based on a template
  from a templatepack.

  Flags:
    -cookbooks-dir=<directory>  - where to create the new cookbook
    -template-pack=<path>       - path to the template pack to use
`
}

func (c *Creator) Run(args []string) int {
	var cookbookPath, templatePack string
	cmdFlags := flag.NewFlagSet("creator", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Println(c.Help()) }
	cmdFlags.StringVar(&cookbookPath, "cookbooks-dir", "", "directory to store new cookbook")
	cmdFlags.StringVar(&templatePack, "template-pack", DefaultTempPack(), "path to the template pack")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	manifest, err := fastfood.NewManifest(path.Join(templatePack, "manifest.json"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	args = cmdFlags.Args()

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
	return 0
}

func (c *Creator) Synopsis() string {
	return "Creates a new cookbook"
}
