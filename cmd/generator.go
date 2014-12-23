package cmd

import (
	"errors"
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

func (g *Generator) Run(args []string) error {
	var templatePack string
	var force bool
	workingDir, _ := os.Getwd()
	cmdFlags := flag.NewFlagSet("gen", flag.ContinueOnError)
	cmdFlags.BoolVar(&force, "force", false, "overwrite existing files")

	if err := cmdFlags.Parse(args); err != nil {
		return err
	}

	remainingArgs := cmdFlags.Args()
	manifest := g.Manifest

	// Expect this command to be run inside a cookbook
	if !fastfood.PathIsCookbook(workingDir) {
		return errors.New("You must run this command from a cookbook directory")
	}

	// Create a new cookbook
	ckbk, err := fastfood.NewCookbookFromPath(workingDir)
	if err != nil {
		return errors.New("Unable to parse cookbook")
	}

	providers, err := g.LoadProviders(ckbk)
	if err != nil {
		return err
	}

	// Remove the first arg as the command
	var genCommand string
	if len(remainingArgs) > 0 {
		genCommand, remainingArgs = remainingArgs[0], remainingArgs[1:len(remainingArgs)]
	} else {
		fmt.Println(manifest.Help())
		return nil
	}

	if _, ok := providers[genCommand]; ok {
		goto CMDFound
	}

	// If the loop finishes without finding the commnad exit
	return errors.New(fmt.Sprintf("No provider found for %s\n", genCommand))

	// Command was found continue to execute
CMDFound:

	p := providers[genCommand]

	// No point in setting up a whole flag set here
	if len(remainingArgs) > 0 {
		if remainingArgs[0] == "-h" {
			fmt.Println(p.Help())
			return nil
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("Error loading provider %s: %v\n", genCommand, err))
	}

	mappedArgs := MapArgs(remainingArgs)
	var providerType string
	if val, ok := mappedArgs["type"]; ok {
		providerType = val
	} else {
		if p.DefaultType != "" {
			providerType = p.DefaultType
		} else {
			return errors.New("You must pass a type b/c not default type is set")
		}
	}

	// Add the needed dependencies to the metadata
	ckbk.AppendDependencies(p.Dependencies(providerType))
	p.GenDirs(providerType)

	err = p.GenFiles(
		providerType,
		path.Join(templatePack, genCommand),
		force,
		mappedArgs,
	)

	if err != nil {
		return errors.New(fmt.Sprintf("Error generating files %v\n", err))
	}

	return nil
}

func (g *Generator) Description() string {
	return `
      Generates a new recipe for an existing cookbook"

      Options are passed using a key:value notation so
      to set the name you would use the following:

      name:recipe_name
`
}

// Autogenerate based on commands parsed
func (g *Generator) Help() string {
	return "fastfood gen <flags> [provider] [options]"
}
