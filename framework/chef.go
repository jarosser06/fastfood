package framework

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/common/fileutil"
	"github.com/jarosser06/fastfood/common/stringutil"
	"github.com/jarosser06/fastfood/framework/chef"
)

type Chef struct {
	cookbook chef.Cookbook
	options  fastfood.FrameworkOptions
}

// Initializes a framework with options
func (c *Chef) Init(frameworkOptions fastfood.FrameworkOptions) error {
	cpath := path.Join(c.options.Destination, c.options.Name)

	// Create cookbook
	if fileutil.FileExist(cpath) {
		var err error
		c.cookbook, err = chef.NewCookbookFromPath(cpath)
		if err != nil {
			return err
		}
	} else {
		c.cookbook = chef.NewCookbook(cpath, c.options.Name)
	}

	return nil
}

// Generates a basic cookbook structure
func (c *Chef) GenerateBase() ([]string, error) {
	// Generate the cookbook if it doesn't exist
	if !fileutil.FileExist(c.cookbook.Path) {
		err := os.Mkdir(c.cookbook.Path, 0755)
		if err != nil {
			return []string{}, fmt.Errorf(
				"error creating cookbook directory %s: %v",
				c.cookbook.Path,
				err,
			)
		}
	}

	err := c.genDirs(c.options.BaseDirs)
	if err != nil {
		return []string{}, err
	}

	var moddedFile []string
	for _, file := range c.options.BaseFiles {
		cfile := path.Join(c.cookbook.Path, file)
		// If the file already exists then dont overwrite it
		if fileutil.FileExist(cfile) {
			continue
		}

		tfile := path.Join(c.options.TemplateDir, file)
		err := c.genBaseFile(cfile, tfile)
		if err != nil {
			return moddedFile, err
		}

		moddedFile = append(moddedFile, cfile)
	}

	return moddedFile, nil
}

// Generates a stencil
// Takes the stencil name to be generated
func (c *Chef) GenerateStencil(name string, stencilset fastfood.StencilSet, opts map[string]string) ([]string, error) {
	var moddedFiles []string
	stencil := stencilset.Stencils[name]
	opts = stencilset.MergeOpts(name, opts)

	tOpts := struct {
		*fastfood.Helpers
		Cookbook chef.Cookbook
		Options  map[string]string
	}{
		Cookbook: c.cookbook,
		Options:  CapitalizeOptions(opts),
	}

	// Generate the directories
	c.genDirs(append(stencil.Directories))

	// Generate the files
	files := stencil.Files
	pfiles := stencil.Partials

	for cfile, tfile := range files {
		// Find any instance of <NAME> and replace with the name
		cfile = strings.Replace(cfile, "<NAME>", tOpts.Options["Name"], 1)

		var pContent []string
		for _, partial := range pfiles {
			b, err := ioutil.ReadFile(partial)
			if err != nil {
				return moddedFiles, fmt.Errorf("error %v occured while reading file %s", err, partial)
			}

			pContent = append(pContent, string(b))
		}

		err := c.genStencilFile(cfile, tfile, pContent, tOpts)
		if err != nil {
			// don't care if the file exists
			if err.Error() != "file exists" {
				return moddedFiles, err
			}
		}

		moddedFiles = append(moddedFiles, cfile)
	}
	c.cookbook.AppendDependencies(stencilset.Dependencies(name))

	return moddedFiles, nil
}

// Generates a stencil file, meant to be called by Chef interface
func (c *Chef) genStencilFile(cfile string, tfile string, pContent []string, tOpts interface{}) error {
	cfilePath := path.Join(c.cookbook.Path, cfile)

	if fileutil.FileExist(cfilePath) && !c.options.Force {
		return fmt.Errorf("file exists")
	}

	// Read content from
	content := pContent
	b, err := ioutil.ReadFile(tfile)
	if err != nil {
		return fmt.Errorf("error %v occured while reading file %s", err, tfile)
	}
	content = append(content, string(b))

	t, err := fastfood.NewTemplate(cfile, tOpts, content)
	if err != nil {
		return fmt.Errorf("error %v occured while creating template %s", err, tfile)
	}

	t.CleanNewlines()

	if err := t.Flush(cfilePath); err != nil {
		return fmt.Errorf("error %v occured while writing file %s", err, cfile)
	}

	return nil
}

func (c *Chef) genBaseFile(file string, tfile string) error {
	content, err := ioutil.ReadFile(tfile)
	if err != nil {
		return fmt.Errorf("error reading file %s, %v", file, err)
	}

	t, err := fastfood.NewTemplate(file, c.cookbook, []string{string(content)})
	if err != nil {
		return fmt.Errorf("error %v returned while creating new template %s", err, file)
	}

	t.CleanNewlines()
	if err := t.Flush(file); err != nil {
		return fmt.Errorf("error %v while writing file %s", err, file)
	}

	return nil
}

// Generates directories, meant for internal usage to Chef
func (c *Chef) genDirs(dirs []string) error {
	for _, dir := range dirs {
		fpath := path.Join(c.cookbook.Path, dir)

		if !fileutil.FileExist(fpath) {
			err := os.MkdirAll(fpath, 0755)
			if err != nil {
				return fmt.Errorf("error %v occured while creating directory %s", err, fpath)
			}
		}
	}

	return nil
}

func CapitalizeOptions(opts map[string]string) map[string]string {
	cmap := make(map[string]string)

	for k, v := range opts {
		cmap[stringutil.CapitalizeString(k)] = v
	}

	return cmap
}
