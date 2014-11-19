package cookbook

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/jarosser06/fastfood/pkg/util"
)

type OSTarget struct {
	Distro  string
	Version string
}

type Cookbook struct {
	Berks        []string
	Dependencies []string
	Name         string
	Path         string
	Target       OSTarget
	Year         int
}

func NewCookbook(cookbookPath string, name string) Cookbook {
	cookbook := Cookbook{
		Year: time.Now().Year(),
		Path: path.Join(cookbookPath, name),
		Name: name,
	}

	return cookbook
}

func (c *Cookbook) GenFiles() error {

	cookbookFiles := [8]string{
		"Berksfile",
		"CHANGELOG.md",
		"Gemfile",
		"metadata.rb",
		"README.md",
		"recipes/default.rb",
		"test/unit/spec/default_spec.rb",
		"test/unit/spec/spec_helper.rb",
	}

	templateBox, _ := rice.FindBox("templates")
	for _, cookbookFile := range cookbookFiles {
		tmpStr, _ := templateBox.String(cookbookFile)
		t, _ := template.New(cookbookFile).Parse(tmpStr)

		f, err := os.Create(path.Join(c.Path, cookbookFile))
		if err != nil {
			return errors.New(fmt.Sprintf("cookbook.GenFiles(): %v", err))
		}

		// Write template output to a buffer
		var buffer bytes.Buffer
		t.Execute(&buffer, c)

		// Clean up rendered templates and write to file
		cleanStr := util.CollapseNewlines(buffer.String())
		io.WriteString(f, cleanStr)

		f.Close()
	}

	return nil
}

func (c *Cookbook) GenDirs() error {
	err := os.Mkdir(c.Path, 0755)
	if err != nil {
		return errors.New(fmt.Sprintf("Gendirs(): %v", err))
	}

	dirs := [10]string{
		"attributes",
		"files",
		"libraries",
		"providers",
		"recipes",
		"resources",
		"templates",
		"test",
		"test/unit",
		"test/unit/spec",
	}

	for _, dir := range dirs {
		err := os.Mkdir(path.Join(c.Path, dir), 0755)
		if err != nil {
			return errors.New(fmt.Sprintf("cookbook.Gendirs(): %v", err))
		}
	}

	return nil
}

func CWDIsCookbook() bool {
	workingDir, _ := os.Getwd()

	_, err := os.Stat(path.Join(workingDir, "metadata.rb"))
	if err == nil {
		return true
	} else {
		return false
	}
}
