package fastfood

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type OSTarget struct {
	Distro  string
	Version string
}

type Cookbook struct {
	*Helpers
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

		cookbook.Path = cookbookPath
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

func (c *Cookbook) GenFiles(cookbookFiles []string, templatesPath string) error {

	for _, cookbookFile := range cookbookFiles {
		tempStr, err := ioutil.ReadFile(path.Join(templatesPath, cookbookFile))

		if err != nil {
			return errors.New(fmt.Sprintf("cookbook.GenFiles() reading template file: %v", err))
		}

		t, err := NewTemplate(cookbookFile, c, []string{string(tempStr)})
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

func (c *Cookbook) GenDirs(cookbookDirs []string) error {
	if !FileExist(c.Path) {
		err := os.Mkdir(c.Path, 0755)
		if err != nil {
			return errors.New(fmt.Sprintf("cookbook.Gendirs(): %v", err))
		}
	}

	for _, dir := range cookbookDirs {
		if !FileExist(path.Join(c.Path, dir)) {
			err := os.MkdirAll(path.Join(c.Path, dir), 0755)
			if err != nil {
				return errors.New(fmt.Sprintf("cookbook.Gendirs(): %v", err))
			}
		}
	}

	return nil
}

func (c *Cookbook) AppendDependencies(dependencies []string) {
	var depBuffer []string

	if len(dependencies) < 0 {
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

		AppendFile(
			path.Join(c.Path, "metadata.rb"),
			fmt.Sprintf("%s\n", strings.Join(depBuffer, "\n")),
		)
	}
}
