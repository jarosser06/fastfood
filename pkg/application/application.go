package application

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/jarosser06/fastfood/pkg/cookbook"
	"github.com/jarosser06/fastfood/pkg/util"
)

const (
	defaultRoot      = "/var/www"
	defaultType      = "generic"
	defaultOwner     = "node['apache']['user']"
	defaultWebserver = "apache"
	defaultRepo      = "github.com/jarosser06/magic"
)

type Application struct {
	Cookbook  cookbook.Cookbook
	Path      string `json:"name,omitempty"`
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
		Cookbook:  ckbk,
		Path:      path.Join(defaultRoot, name),
		Name:      name,
		Owner:     defaultOwner,
		Repo:      defaultRepo,
		Root:      defaultRoot,
		Type:      defaultType,
		Webserver: defaultWebserver,
	}
}

// Using set root also updates the Path variable
func (a *Application) SetRoot(root string) {
	a.Root = root
	a.Path = path.Join(root, a.Name)
}

func (a *Application) GenFiles() error {

	cookbookFiles := map[string]string{
		fmt.Sprintf("%s.rb", a.Name):      "recipes/application.rb",
		fmt.Sprintf("%s_spec.rb", a.Name): "test/unit/spec/application_spec.rb",
	}

	templateBox, _ := rice.FindBox("templates")
	for cookbookFile, templateFile := range cookbookFiles {
		tmpStr, _ := templateBox.String(templateFile)
		t, _ := template.New(templateFile).Parse(tmpStr)

		f, err := os.Create(path.Join(a.Cookbook.Path, cookbookFile))
		defer f.Close()

		if err != nil {
			return errors.New(fmt.Sprintf("application.GenFiles(): %v", err))
		}

		var buffer bytes.Buffer
		t.Execute(&buffer, a)

		cleanStr := util.CollapseNewlines(buffer.String())
		io.WriteString(f, cleanStr)
	}
	return nil
}
