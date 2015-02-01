package chef

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jarosser06/fastfood/common/fileutil"
)

type BerksFile struct {
	Cookbooks map[string]BerksCookbook
}

type BerksCookbook struct {
	Branch   string `json:"branch"`
	Name     string `json:"name"`
	Git      string `json:"git"`
	Path     string `json:"path"`
	Ref      string `json:"ref"`
	Revision string `json:"revision"`
	Tag      string `json:"tag"`
}

// Read a berkshelf file and return a berks struct
func BerksFromFile(f string) (BerksFile, error) {
	b := BerksFile{Cookbooks: make(map[string]BerksCookbook)}

	r, err := os.Open(f)
	if err != nil {
		return b, err
	}
	defer r.Close()

	b.Parse(r)

	return b, nil
}

func (c *BerksCookbook) String() string {
	s := fmt.Sprintf("cookbook \"%s\"", c.Name)
	if c.Git != "" {
		s = fmt.Sprintf("%s, git: \"%s\"", s, c.Git)

		switch {
		case c.Branch != "":
			s = fmt.Sprintf("%s, branch: \"%s\"", s, c.Branch)
		case c.Ref != "":
			s = fmt.Sprintf("%s, ref: \"%s\"", s, c.Ref)
		case c.Revision != "":
			s = fmt.Sprintf("%s, revision: \"%s\"", s, c.Revision)
		case c.Tag != "":
			s = fmt.Sprintf("%s, tag: \"%s\"", s, c.Tag)
		}

	} else if c.Path != "" {
		s = fmt.Sprintf("%s, path: \"%s\"", s, c.Path)
	}

	return s
}

// Parse the Berksfile
func (b *BerksFile) Parse(r io.Reader) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	for s.Scan() {
		switch s.Text() {
		case "cookbook":
			s.Scan()

			cName := strings.Trim(s.Text(), "',\"")
			c := BerksCookbook{Name: cName}
			b.Cookbooks[cName] = c
		}
	}
}

// Append Dependencies to a Berksfile
func (b *BerksFile) Append(f string, deps []BerksCookbook) []string {
	var added []string
	var buffer []string

	if len(deps) == 0 {
		return added
	}

	// Catches issue with Cookbooks not being created
	if b.Cookbooks == nil {
		b.Cookbooks = make(map[string]BerksCookbook)
	}

	for _, d := range deps {
		if _, ok := b.Cookbooks[d.Name]; !ok {
			b.Cookbooks[d.Name] = d
			added = append(added, d.Name)
			buffer = append(buffer, d.String())
		}
	}

	if len(added) > 0 {
		fileutil.AppendFile(
			f,
			fmt.Sprintf("%s\n", strings.Join(buffer, "\n")),
		)
	}

	return added
}
