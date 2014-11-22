package application

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/GeertJohan/go.rice"
	"github.com/jarosser06/fastfood/pkg/cookbook"
	"github.com/jarosser06/fastfood/pkg/template"
	"github.com/jarosser06/fastfood/pkg/util"
)

const (
	defaultBranch    = "master"
	defaultRoot      = "/var/www"
	defaultType      = "generic"
	defaultOwner     = "node['apache']['user']"
	defaultWebserver = "apache"
	defaultRepo      = "git@github.com/jarosser06/magic"
)

type Application struct {
	Branch    string `json:"branch,omitempty"`
	Cookbook  cookbook.Cookbook
	Name      string `json:"name,omitempty"`
	Owner     string `json:"owner,omitempty"`
	Repo      string `json:"repo,omitempty"`
	Root      string `json:"docroot,omitempty"`
	Type      string `json:"type,omitempty"`
	Webserver string `json:"webserver,omitempty"`
}

// Return an application with the defaults
func NewApplication(name string, ckbk cookbook.Cookbook) Application {

	return Application{
		Branch:    defaultBranch,
		Cookbook:  ckbk,
		Name:      name,
		Owner:     defaultOwner,
		Repo:      defaultRepo,
		Root:      defaultRoot,
		Type:      defaultType,
		Webserver: defaultWebserver,
	}
}

func (a *Application) Path() string {
	return path.Join(a.Root, a.Name)
}

func (a *Application) GenFiles() error {

	cookbookFiles := map[string]string{
		fmt.Sprintf("recipes/%s.rb", a.Name):             "recipes/application.rb",
		fmt.Sprintf("test/unit/spec/%s_spec.rb", a.Name): "test/unit/spec/application_spec.rb",
	}

	templateBox, _ := rice.FindBox("templates")
	for cookbookFile, templateFile := range cookbookFiles {
		tmpStr, _ := templateBox.String(templateFile)
		partialStr, _ := templateBox.String("partials/site_setup.rb")

		t, err := template.NewTemplate(cookbookFile, a, tmpStr, partialStr)

		if err != nil {
			return errors.New(fmt.Sprintf("Error creating template: %v", err))
		}

		t.CleanNewlines()

		if err := t.Flush(path.Join(a.Cookbook.Path, cookbookFile)); err != nil {
			return errors.New(fmt.Sprintf("Error writing file: %v", err))
		}

	}
	return nil
}

func (a *Application) GenDirs() error {
	dirs := [2]string{
		"recipes",
		"test/unit/spec",
	}

	for _, dir := range dirs {
		fullPath := path.Join(a.Cookbook.Path, dir)

		if !util.FileExist(fullPath) {
			err := os.MkdirAll(path.Join(a.Cookbook.Path, dir), 0755)
			if err != nil {
				return errors.New(fmt.Sprintf("application.GenDirs(): %v", err))
			}
		}
	}

	return nil
}
