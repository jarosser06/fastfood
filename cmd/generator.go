package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jarosser06/fastfood/pkg/cookbook"
)

type Generator struct {
	MappedArgs map[string]string
}

func GenApp(args map[string]string) {
	fmt.Println("Shit works")
	//TODO: Implement Application Generation code
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

func (g *Generator) Help() string {
	return "TODO: Add help text for gen command"
}

func (g *Generator) Run(args []string) int {
	workingDir, _ := os.Getwd()

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
		}

		// Attempt to call the command function if it exists
		if com, ok := genCommands[args[0]]; ok {
			com.(func(map[string]string))(mappedArgs)
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
