package cookbook

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/jarosser06/fastfood/provider"
	"github.com/jarosser06/fastfood/util"
)

type OSTarget struct {
	Distro  string
	Version string
}

type Cookbook struct {
	*provider.Helpers
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

// Given a cookbook path, return a cookbook struct pre-populated
func NewCookbookFromPath(cookbookPath string) (Cookbook, error) {
	cookbook := Cookbook{Year: time.Now().Year()}

	if PathIsCookbook(cookbookPath) {
		f, err := os.Open(path.Join(cookbookPath, "metadata.rb"))

		if err != nil {
			return cookbook, errors.New(
				fmt.Sprintf("cookbook.NewCookbookFromPath: %v", err),
			)
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)

		depFlag := false
		nameFlag := false

		for scanner.Scan() {
			switch {
			case depFlag:
				cookbook.Dependencies = append(
					cookbook.Dependencies,
					strings.Trim(scanner.Text(), "'"),
				)
				depFlag = false
			case nameFlag:
				cookbook.Name = strings.Trim(scanner.Text(), "'")
				nameFlag = false
			case scanner.Text() == "depends":
				depFlag = true
			case scanner.Text() == "name":
				nameFlag = true
			}
		}

		if !(len(cookbook.Name) > 0) {
			return cookbook, errors.New("unable to determine cookbook name")
		}

		return cookbook, nil
	} else {
		return cookbook, errors.New(
			fmt.Sprintf("%s is not a cookbook", cookbookPath),
		)
	}
}

func PathIsCookbook(cookbookPath string) bool {
	_, err := os.Stat(path.Join(cookbookPath, "metadata.rb"))
	if err == nil {
		return true
	} else {
		return false
	}
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

	templateBox, _ := rice.FindBox("../templates/cookbook")
	for _, cookbookFile := range cookbookFiles {
		tmpStr, _ := templateBox.String(cookbookFile)

		t, err := provider.NewTemplate(cookbookFile, c, []string{tmpStr})
		if err != nil {
			return errors.New(fmt.Sprintf("cookbook.GenFiles(): %v", err))
		}

		t.CleanNewlines()
		if err := t.Flush(path.Join(c.Path, cookbookFile)); err != nil {
			return errors.New(fmt.Sprintf("cookbook.GenFiles(): %v", err))
		}
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
		"test/unit/spec",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(path.Join(c.Path, dir), 0755)
		if err != nil {
			return errors.New(fmt.Sprintf("cookbook.Gendirs(): %v", err))
		}
	}

	return nil
}

func (c *Cookbook) AppendDependencies(dependencies []string) {
	var depBuffer []string
	for _, dep := range dependencies {
		exist := false

		for _, existing := range c.Dependencies {
			if existing == dep {
				exist = true
				continue
			}
		}

		if !exist {
			depBuffer = append(depBuffer, fmt.Sprintf("depends '%s'", dep))
		}
	}

	util.AppendFile(
		path.Join(c.Path, "metadata.rb"),
		fmt.Sprintf("\n%s", strings.Join(depBuffer, "\n")),
	)
}
