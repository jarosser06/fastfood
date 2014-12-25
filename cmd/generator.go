package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/jarosser06/fastfood"
)

type Generator struct {
	Common
	MappedArgs    map[string]string
	TemplatesPath string
}

func (g *Generator) Run(args []string) int {
	var force bool
	workingDir, _ := os.Getwd()
	cmdFlags := flag.NewFlagSet("gen", flag.ContinueOnError)
	cmdFlags.BoolVar(&force, "force", false, "overwrite existing files")
	cmdFlags.StringVar(&g.TemplatePack, "template-pack", DefaultTempPack(), "path to the template pack")
	cmdFlags.Usage = func() { fmt.Println(g.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
		return 1
	}

	if err := g.LoadManifest(); err != nil {
		fmt.Println(err)
		return 1
	}

	remainingArgs := cmdFlags.Args()
	manifest := g.Manifest

	// Expect this command to be run inside a cookbook
	if !fastfood.PathIsCookbook(workingDir) {
		fmt.Println("You must run this command from a cookbook directory")
		return 1
	}

	// Create a new cookbook
	ckbk, err := fastfood.NewCookbookFromPath(workingDir)
	if err != nil {
		fmt.Println("Unable to parse cookbook")
		return 1
	}

	providers, err := g.LoadProviders(ckbk)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// Remove the first arg as the command
	var genCommand string
	if len(remainingArgs) > 0 {
		genCommand, remainingArgs = remainingArgs[0], remainingArgs[1:len(remainingArgs)]
	} else {
		fmt.Println(manifest.Help())
		return 0
	}

	if _, ok := providers[genCommand]; ok {
		goto CMDFound
	}

	// If the loop finishes without finding the commnad exit
	fmt.Printf("No provider found for %s\n", genCommand)
	return 1

	// Command was found continue to execute
CMDFound:

	p := providers[genCommand]

	// No point in setting up a whole flag set here
	if len(remainingArgs) > 0 {
		if remainingArgs[0] == "-h" {
			fmt.Println(p.Help())
			return 0
		}
	}

	if err != nil {
		fmt.Printf("Error loading provider %s: %v\n", genCommand, err)
		return 1
	}

	mappedArgs := MapArgs(remainingArgs)
	var providerType string
	if val, ok := mappedArgs["type"]; ok {
		providerType = val
	} else {
		if p.DefaultType != "" {
			providerType = p.DefaultType
		} else {
			fmt.Println("You must pass a type b/c not default type is set")
			return 1
		}
	}

	// Add the needed dependencies to the metadata
	ckbk.AppendDependencies(p.Dependencies(providerType))
	p.GenDirs(providerType)

	err = p.GenFiles(
		providerType,
		path.Join(g.TemplatePack, genCommand),
		force,
		mappedArgs,
	)

	if err != nil {
		fmt.Printf("Error generating files %v\n", err)
		return 1
	}

	fmt.Printf("Cookbook %s updated\n", ckbk.Path)
	return 0
}

func (g *Generator) Synopsis() string {
	return "Generates a new recipe for an existing cookbook"
}

// Autogenerate based on commands parsed
func (g *Generator) Help() string {
	return `
Usage: fastfood gen <flags> [provider] [options]
  This will generate a recipe and spec file
  based on the provider and options you
  provide that provider.
  Options are passed using using a key:value
  notation so to set the name you would use
  the following:
  name:recipe_name

Flags:
  -template-pack=<path> - path to the template pack
  -force                - overwrite any existing files
`
}
