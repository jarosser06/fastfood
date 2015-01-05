package cmd

import (
	"flag"
	"fmt"
	"path"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/framework"
)

type Creator struct {
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
	var cookbooksPath, templatePack string
	cmdFlags := flag.NewFlagSet("new", flag.ContinueOnError)
	cmdFlags.StringVar(&templatePack, "template-pack", DefaultTempPack(), "path to the template pack")
	cmdFlags.StringVar(&cookbooksPath, "cookbooks-path", DefaultCookbooksPath(), "base cookbooks directory")
	cmdFlags.Usage = func() { fmt.Println(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
		return 1
	}

	remainingArgs := cmdFlags.Args()
	if len(remainingArgs) > 0 {
		name := remainingArgs[0]

		manifest, err := fastfood.NewManifest(path.Join(templatePack, "manifest.json"))
		if err != nil {
			fmt.Println(err)
			return 1
		}

		fopts := fastfood.FrameworkOptions{
			Destination: cookbooksPath,
			BaseFiles:   manifest.Frameworks["chef"].BaseFiles,
			BaseDirs:    manifest.Frameworks["chef"].BaseDirectories,
			Name:        name,
			TemplateDir: templatePack,
		}

		c := framework.Chef{}
		c.Init(fopts)
		_, err = c.GenerateBase()
		if err != nil {
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
