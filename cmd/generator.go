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
		// Call NewCookbook within this path
	} else {
		fmt.Println("You must run this command from a cookbook directory")
		return 1
	}
	return 0
}

func (g *Generator) Synopsis() string {
	return "Generates a new recipe for an existing cookbook"
}
