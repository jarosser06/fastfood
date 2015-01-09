package chef

import (
	"bufio"
	"io"
	"os"
)

type BerksFile struct {
	Cookbooks []BerksCookbook
}

type BerksCookbook struct {
	Name     string `json:"name"`
	Git      string `json:"git"`
	Path     string `json:"path"`
	Revision string `json:"revision"`
}

// Read a berkshelf file and return a berks struct
func BerksFromFile(f string) (BerksFile, error) {
	var b BerksFile

	r, err := os.Open(f)
	if err != nil {
		return b, err
	}
	defer r.Close()

	b.Parse(r)

	return b, nil
}

// Parse the Berksfile
func (b *BerksFile) Parse(r io.Reader) {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)

	for s.Scan() {
		switch s.Text() {
		case "cookbook":
			s.Scan()
			c := BerksCookbook{Name: s.Text()}
			b.Cookbooks = append(b.Cookbooks, c)
		}
	}
}

// Append Dependencies to a Berksfile
func (b *BerksFile) Append(deps []BerksCookbook) []string {

	return []string{}
}
