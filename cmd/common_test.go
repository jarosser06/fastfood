package cmd

import (
	"os"
	"testing"
)

func TestMapArgs(t *testing.T) {
	args := []string{"name:testing"}

	mappedArgs := MapArgs(args)

	if val, ok := mappedArgs["name"]; !ok && val == "testing" {
		t.Errorf("Expected a map with a key of name and value of testing")
	}
}

func TestDefaultCookbookPath(t *testing.T) {
	fakePath := "/my/cookbook/path"
	os.Setenv("COOKBOOKS", fakePath)
	defer os.Unsetenv("COOKBOOKS")

	if res := DefaultCookbookPath(); res != fakePath {
		t.Errorf("Expected %s to be returned", fakePath)
	}
}
