package cmd

import (
	"errors"
	"path"

	"github.com/jarosser06/fastfood"
)

type Creator struct {
	Common
}

func (c *Creator) Help() string {
	return "fastfood new [cookbook_name]"
}

func (c *Creator) Run(args []string) error {
	if len(args) > 0 {
		cookbook := fastfood.NewCookbook(c.cookbookPath, args[0])

		//TODO: These can be collapsed into a single function
		if err := cookbook.GenDirs(c.Manifest.Cookbook.Directories); err != nil {
			return err
		}

		templatePath := path.Join(c.templatePack, c.Manifest.Cookbook.TemplatesPath)
		if err := cookbook.GenFiles(c.Manifest.Cookbook.Files, templatePath); err != nil {
			return err
		}

	} else {
		return errors.New("You must enter the name of the cookbook")
	}
	return nil
}

func (c *Creator) Description() string {
	return "Simple commmand that generates a new cookbook based on a template from a templatepack."
}
