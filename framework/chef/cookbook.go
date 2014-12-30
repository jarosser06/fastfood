package chef

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jarosser06/fastfood"
	"github.com/jarosser06/fastfood/common/fileutil"
)

type OSTarget struct {
	Distro  string
	Version string
}

type Cookbook struct {
	*fastfood.Helpers
	Berks        []string
	Dependencies []string
	Name         string
	Path         string
	Target       OSTarget
	Year         int
}

// Returns a new empty cookbook
func NewCookbook(cookbookPath string, name string) Cookbook {
	return Cookbook{
		Year: time.Now().Year(),
		Path: cookbookPath,
		Name: name,
	}
}

// Given a cookbook path, return a cookbook struct pre-populated
func NewCookbookFromPath(cookbookPath string) (Cookbook, error) {
	cookbook := Cookbook{Year: time.Now().Year()}

	if PathIsCookbook(cookbookPath) {
		f, err := os.Open(path.Join(cookbookPath, "metadata.rb"))

		if err != nil {
			return cookbook, fmt.Errorf("cookbook.NewCookbookFromPath: %v", err)
		}

		defer f.Close()

		// TODO: This should be cleaned up and stuck elsewhere
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

// Returns a list of dependencies that were written
func (c *Cookbook) AppendDependencies(dependencies []string) []string {
	var depBuffer, newDeps []string

	if len(dependencies) > 0 {
		for _, dep := range dependencies {
			exist := false

			for _, existing := range c.Dependencies {
				if existing == dep {
					exist = true
					continue
				}
			}

			if !exist {
				// Keep track of all new dependencies
				newDeps = append(newDeps, dep)
				depBuffer = append(depBuffer, fmt.Sprintf("depends '%s'", dep))
			}
		}

		// Don't append newlines if all dependencies are up to date
		if len(depBuffer) > 0 {
			fileutil.AppendFile(
				path.Join(c.Path, "metadata.rb"),
				fmt.Sprintf("%s\n", strings.Join(depBuffer, "\n")),
			)
		}
	}

	// Add the new dependencies to the cookbook interface
	c.Dependencies = append(c.Dependencies, newDeps...)

	return newDeps
}
