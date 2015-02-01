package chef

import (
	"os"
	"path"
	"testing"

	"github.com/jarosser06/fastfood/common/fileutil"
)

var testCookbook = "/tmp/testcookbook"

func TestMain(m *testing.M) {
	if _, err := os.Stat(testCookbook); err == nil {
		os.RemoveAll(testCookbook)
	}

	os.Mkdir(testCookbook, 0755)
	fileutil.Copy("../../tests/testcookbook/metadata.rb", path.Join(testCookbook, "metadata.rb"))
	status := m.Run()
	os.RemoveAll(testCookbook)

	os.Exit(status)
}

func TestNewCookbook(t *testing.T) {
	c := NewCookbook("/tmp", "testcookbook")

	if c.Path != path.Join("/tmp", "testcookbook") {
		t.Errorf("unexpected path %s from new cookbook")
	}
}

func TestPathIsCookbook(t *testing.T) {
	if !PathIsCookbook(testCookbook) {
		t.Errorf("%s should be a valid cookbook", testCookbook)
	}
}

func TestNewCookbookFromPath(t *testing.T) {
	cookbook, err := NewCookbookFromPath(testCookbook)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if cookbook.Name != "testcookbook" {
		t.Errorf("expected cookbook name to be testcookbook")
	}

	if cookbook.Path != testCookbook {
		t.Errorf("expected cookbook path to be %s not %s", testCookbook, cookbook.Path)
	}

	if len(cookbook.Dependencies) != 1 {
		t.Errorf("cookbook should have one dependencies")
	}
}

func TestAppendDependencies(t *testing.T) {
	c, _ := NewCookbookFromPath(testCookbook)
	deps := map[string]CookbookDependency{
		"couchdb": {Name: "couchdb", Options: []string{">= 2.5.3"}},
		"mongodb": {Name: "mongodb"},
	}
	appended := c.AppendDependencies(deps)

	if len(appended) != 2 {
		t.Errorf("expected 2 dependencies to be added")
	}

	if len(c.Dependencies) != 3 {
		t.Errorf("expected 3 dependencies to be added to cookbook struct")
	}

	appended = c.AppendDependencies(deps)
	if len(appended) != 0 {
		t.Errorf("expected AppendDependencies to have 0 elements")
	}

	if len(c.Dependencies) != 3 {
		t.Errorf("expected cookbook dependencies to have 3 elements")
	}
}
