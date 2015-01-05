package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/framework"
	"github.com/jarosser06/fastfood/framework/chef"
)

type Generator struct {
}

func (g *Generator) Run(args []string) int {
	var force bool
	var templatePack string
	workingDir, _ := os.Getwd()
	cmdFlags := flag.NewFlagSet("gen", flag.ContinueOnError)
	cmdFlags.BoolVar(&force, "force", false, "overwrite existing files")
	cmdFlags.StringVar(&templatePack, "template-pack", DefaultTempPack(), "path to the template pack")
	cmdFlags.Usage = func() { fmt.Println(g.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		fmt.Println(err)
		return 1
	}

	manifest, err := fastfood.NewManifest(path.Join(templatePack, "manifest.json"))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	remainingArgs := cmdFlags.Args()

	// Expect this command to be run inside a cookbook
	if !chef.PathIsCookbook(workingDir) {
		fmt.Println("You must run this command from a cookbook directory")
		return 1
	}

	fopts := fastfood.FrameworkOptions{
		Destination: workingDir,
		Force:       force,
		TemplateDir: templatePack,
	}

	// Initialize the Chef framework
	c := framework.Chef{}
	c.Init(fopts)

	// Remove the first arg as the command
	var setName string
	if len(remainingArgs) > 0 {
		setName, remainingArgs = remainingArgs[0], remainingArgs[1:len(remainingArgs)]
	} else {
		fmt.Println(manifest.Help())
		return 0
	}

	if _, ok := manifest.StencilSets[setName]; ok {
		goto CMDFound
	}

	// If the loop finishes without finding the commnad exit
	fmt.Printf("No provider found for %s\n", setName)
	return 1

	// Command was found continue to execute
CMDFound:

	s, err := fastfood.NewStencilSet(manifest.StencilSets[setName].Manifest)
	if err != nil {
		fmt.Printf("opening stencil set returned %v", err)
		return 1
	}

	if _, err := s.Valid(); err != nil {
		fmt.Printf("invalid stencilset %v", err)
	}

	// No point in setting up a whole flag set here
	if len(remainingArgs) > 0 {
		if remainingArgs[0] == "-h" {
			fmt.Println(s.Help())
			return 0
		}
	}

	mappedArgs := MapArgs(remainingArgs)
	var stencil string
	if val, ok := mappedArgs["stencil"]; ok {
		stencil = val
	} else {
		stencil = s.DefaultStencil
	}

	if _, ok := s.Stencils[stencil]; !ok {
		fmt.Printf("%s is not a valid stencil for stencilset %s\n", stencil, setName)
		return 1
	}

	updatedFiles, err := c.GenerateStencil(stencil, s, mappedArgs)
	if err != nil {
		fmt.Printf("Error generating files %v\n", err)
		return 1
	}

	if len(updatedFiles) > 0 {
		fmt.Println("cookbook has been updated")
	}
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
