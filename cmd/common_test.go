package cmd

import "testing"

func TestMapArgs(t *testing.T) {
	args := []string{"name:testing"}

	mappedArgs := MapArgs(args)

	if val, ok := mappedArgs["name"]; !ok && val == "testing" {
		t.Errorf("Expected a map with a key of name and value of testing")
	}
}
