package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood/provider/cookbook"
	"github.com/jarosser06/fastfood/util"
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
	workingDir, _ := os.Getwd()
	cmdFlags := flag.NewFlagSet("gen", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Println(g.Help()) }
	//templatesPath := cmdFlags.String("templates-path", "samples", "path to the templates directory")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if cookbook.PathIsCookbook(workingDir) {
		/*
			ckbk, err := cookbook.NewCookbookFromPath(workingDir)

			if err != nil {
				fmt.Println("Unable to parse cookbook")
				return 1
			}
		*/

		cmdManifest := path.Join("/home/jim/Projects/fastfood/samples", "manifest.json")
		if !util.FileExist(cmdManifest) {
			fmt.Printf("Error no such file %s\n", cmdManifest)
			return 1
		}

		// Remove the first arg as the command
		genCommand, args := args[0], args[1:len(args)]

		commands := ParseCommandsFromFile(cmdManifest)
		for _, command := range commands {
			if command.Name == genCommand {
				goto CMDFound
			}
		}

		// If the loop finishes without finding the commnad exit
		fmt.Printf("No generator found for %s\n", genCommand)
		return 1

		// Command was found continue to execute
	CMDFound:

		fmt.Println(MapArgs(args))

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
