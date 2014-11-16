package cookbook

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
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

func NewCookbook(path string) Cookbook {
	cookbook := Cookbook{
		Year: time.Now().Year(),
		Path: path,
	}

	return cookbook
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

	_, err := os.Stat(path.Join(workingDir, "metatdata.rb"))
	if err == nil {
		return true
	} else {
		return false
	}
}
