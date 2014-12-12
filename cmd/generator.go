package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood"
)

type Generator struct {
	MappedArgs    map[string]string
	TemplatesPath string
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

func (g *Generator) Run(args []string) int {
	var templatePack string
	workingDir, _ := os.Getwd()
	cmdFlags := flag.NewFlagSet("gen", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Println(g.Help()) }
	cmdFlags.StringVar(&templatePack, "template-pack", DefaultTempPack(), "path to the template pack")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if fastfood.PathIsCookbook(workingDir) {
		ckbk, err := fastfood.NewCookbookFromPath(workingDir)

		if err != nil {
			fmt.Println("Unable to parse cookbook")
			return 1
		}

		cmdManifest := path.Join(templatePack, "manifest.json")
		if !fastfood.FileExist(cmdManifest) {
			fmt.Printf("Error no such file %s\n", cmdManifest)
			return 1
		}

		// Remove the first arg as the command
		genCommand, args := args[0], args[1:len(args)]

		manifest := NewManifest(cmdManifest)
		if _, ok := manifest.Providers[genCommand]; ok {
			goto CMDFound
		}

		// If the loop finishes without finding the commnad exit
		fmt.Printf("No generator found for %s\n", genCommand)
		return 1

		// Command was found continue to execute
	CMDFound:

		p, err := fastfood.NewProviderFromFile(
			ckbk,
			path.Join(templatePack, manifest.Providers[genCommand].Manifest),
		)

		// No point in setting up a whole flag set here
		if len(args) > 0 {
			if args[0] == "-h" {
				fmt.Println(p.Help())
				return 0
			}
		}

		if err != nil {
			fmt.Println("Error loading provider %s: %v", genCommand, err)
			return 1
		}

		mappedArgs := MapArgs(args)
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
			path.Join(templatePack, genCommand),
			mappedArgs,
		)

		if err != nil {
			fmt.Printf("Error generating files %v\n", err)
			return 1
		}

	} else {
		fmt.Println("You must run this command from a cookbook directory")
		return 1
	}
	return 0
}

func (g *Generator) Synopsis() string {
	return "Generates a new recipe for an existing cookbook"
}

// Autogenerate based on commands parsed
func (g *Generator) Help() string {
	helpText := `
Usage: fastfood gen <flags> [provider] [options]

  This will generate a recipe and spec file
  based on the provider and options you
  provide that provider.

  Options are passed using using a key:value
  notation so to set the name you would use
  the following:

  name:recipe_name

  Flags:

    -template-pack=<path>  - is optional
`

	return helpText
}
