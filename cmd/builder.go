package cmd

import (
	"flag"
	"fmt"

	"github.com/jarosser06/fastfood/provider"
)

type Builder struct {
	CookbooksPath string
}

func (b *Builder) Help() string {
	return "TODO: Add help text for new command"
}

func (b *Builder) Run(args []string) int {
	var config, cookbookPath string
	cmdFlags := flag.NewFlagSet("builder", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Println(b.Help()) }
	cmdFlags.StringVar(&cookbookPath, "cookbooks-dir", "", "directory to store new cookbook")
	cmdFlags.StringVar(&config, "config", "", "json config to build from")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	args = cmdFlags.Args()

	if len(config) > 0 {
		fmt.Println("Placeholder for generating from a config")
	} else {
		if len(args) > 0 {
			cookbook := provider.NewCookbook(cookbookPath, args[0])

			//TODO: These can be collapsed into a single function
			if err := cookbook.GenDirs(); err != nil {
				fmt.Println(err)
				return 1
			}

			if err := cookbook.GenFiles(); err != nil {
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
