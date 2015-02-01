package chef

import (
	"os"
	"testing"

	"github.com/jarosser06/fastfood/common/fileutil"
)

const berksTestFile = "../../tests/Berksfile"

func TestBerksFromFile(t *testing.T) {
	b, err := BerksFromFile(berksTestFile)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if _, ok := b.Cookbooks["apt"]; !ok {
		t.Errorf("expected apt to be a parsed cookbook")
	}
}

func TestBerksAppend(t *testing.T) {
	tmpFile := "/tmp/Berksfile"
	fileutil.Copy(berksTestFile, tmpFile)

	defer os.Remove(tmpFile)

	b, _ := BerksFromFile(tmpFile)

	newDeps := []BerksCookbook{
		BerksCookbook{Name: "apt"},
		BerksCookbook{Name: "apache", Git: "https://github.com/viverae-cookbooks/apache2"},
	}

	m := b.Append(tmpFile, newDeps)

	if len(m) != 1 {
		t.Errorf("expected 1 dependency to be added got %d", len(m))
	}
}

func TestBerksAppend_withEmptySlice(t *testing.T) {
	tmpFile := "/tmp/Berksfile"
	fileutil.Copy(berksTestFile, tmpFile)

	defer os.Remove(tmpFile)

	b, _ := BerksFromFile(tmpFile)

	var newDeps []BerksCookbook

	b.Append(tmpFile, newDeps)
}

func TestBerksAppend_wontDup(t *testing.T) {
	tmpFile := "/tmp/Berksfile"
	fileutil.Copy(berksTestFile, tmpFile)

	defer os.Remove(tmpFile)

	b, _ := BerksFromFile(tmpFile)

	newDeps := []BerksCookbook{
		BerksCookbook{Name: "apt"},
		BerksCookbook{Name: "apache", Git: "https://github.com/viverae-cookbooks/apache2"},
		BerksCookbook{Name: "couchdb", Git: "https://github.com/the-galley/chef-couchdb", Ref: "88a9afa3b7e29ce987d7f4f43df5d2070491ed05"},
	}

	// This run returns 1
	b.Append(tmpFile, newDeps)

	// This run should return 0
	m := b.Append(tmpFile, newDeps)
	if len(m) > 0 {
		t.Errorf("expected no changes to be made but got %d", len(m))
	}
}
