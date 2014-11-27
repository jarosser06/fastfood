package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jarosser06/fastfood/provider/application"
	"github.com/jarosser06/fastfood/provider/cookbook"
	"github.com/jarosser06/fastfood/provider/database"
	"github.com/mitchellh/mapstructure"
)

type Generator struct {
	MappedArgs map[string]string
}

func GenApp(ckbk cookbook.Cookbook, args map[string]string) {
	app := application.New("app", ckbk)

	mapstructure.Decode(args, &app)

	if err := app.GenDirs(); err != nil {
		fmt.Printf("Error creating dirs: %v\n", err)
	}

	if err := app.GenFiles(); err != nil {
		fmt.Printf("Error creating application: %v\n", err)
	}

	app.Cookbook.AppendDependencies(app.Dependencies())
}

func GenDB(ckbk cookbook.Cookbook, args map[string]string) {
	db := database.New("db", ckbk)

	mapstructure.Decode(args, &db)

	if err := db.GenDirs(); err != nil {
		fmt.Printf("Error creating dirs: %v\n", err)
	}

	if err := db.GenFiles(); err != nil {
		fmt.Printf("Error creating database: %v\n", err)
	}

	db.Cookbook.AppendDependencies(db.Dependencies())
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

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if cookbook.PathIsCookbook(workingDir) {
		ckbk, err := cookbook.NewCookbookFromPath(workingDir)

		if err != nil {
			fmt.Println("Unable to parse cookbook")
			return 1
		}

		// Map the specific gen function to the passed arg
		mappedArgs := MapArgs(args)
		genCommands := map[string]interface{}{
			"app": GenApp,
			"db":  GenDB,
		}

		// Remove the first arg as the command
		genCommand, args := args[0], args[1:len(args)]
		// Attempt to call the command function if it exists
		if com, ok := genCommands[genCommand]; ok {
			com.(func(cookbook.Cookbook, map[string]string))(ckbk, mappedArgs)
		} else {
			fmt.Printf("No generator for %s\n", args[0])
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
