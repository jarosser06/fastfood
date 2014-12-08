package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood/provider"
	"github.com/jarosser06/fastfood/util"
)

const templatePack = "/home/jim/Projects/fastfood/samples"

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
	workingDir, _ := os.Getwd()
	cmdFlags := flag.NewFlagSet("gen", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Println(g.Help()) }
	//templatesPath := cmdFlags.String("templates-path", "samples", "path to the templates directory")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if provider.PathIsCookbook(workingDir) {
		ckbk, err := provider.NewCookbookFromPath(workingDir)

		if err != nil {
			fmt.Println("Unable to parse cookbook")
			return 1
		}

		cmdManifest := path.Join(templatePack, "manifest.json")
		if !util.FileExist(cmdManifest) {
			fmt.Printf("Error no such file %s\n", cmdManifest)
			return 1
		}

		// Remove the first arg as the command
		genCommand, args := args[0], args[1:len(args)]

		commands := ParseCommandsFromFile(cmdManifest)
		if _, ok := commands[genCommand]; ok {
			goto CMDFound
		}

		// If the loop finishes without finding the commnad exit
		fmt.Printf("No generator found for %s\n", genCommand)
		return 1

		// Command was found continue to execute
	CMDFound:

		p := ParseProviderFromFile(
			ckbk,
			path.Join(templatePack, commands[genCommand].Manifest),
		)

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
Usage: fastfood gen [provider] [options]

  This will generate a recipe and spec file
  based on the provider and options you
  provide that provider.

  Options are passed using using a key:value
  notation so to set the name you would use
  the following:

  name:recipe_name

Generators:

  db     - Creates a database recipe based
           on the type, defaults to MySQL

  app    - Creates an application recipe
           based on the type, defaults to Generic`

	return helpText
}
