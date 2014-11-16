package util

import (
	"strings"
	"testing"
)

func TestCollapseNewlines(t *testing.T) {
	testString := "include_recipe 'nginx'\n\n\n"

	res := CollapseNewlines(testString)

	if strings.Contains(res, "\n\n\n") {
		t.Errorf("Expected resulting string not to contain 3 newlines")
	}

	if !strings.Contains(res, "\n\n") {
		t.Errorf("Expected new string to contain two newlines")
	}
}
