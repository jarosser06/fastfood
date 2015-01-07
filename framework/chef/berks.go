package chef

import (
	"io"
	"os"
)

type Berks struct {
	Git      string `json:"git"`
	Path     string `json:"path"`
	Revision string `json:"revision"`
}

// Read a berkshelf file and return a berks struct
func BerksFromFile(f string) (Berks, error) {
	var b Berks

	r, err := os.Open(f)
	if err != nil {
		return b, err
	}
	defer r.Close()

	b.Parse(r)

	return b, nil
}

// Parse the Berksfile
func (b *Berks) Parse(r io.Reader) {
}
