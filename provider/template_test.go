package provider

import (
	"strings"
	"testing"
)

func TestCleanNewlines(t *testing.T) {
	temp := Template{
		Content: "include_recipe 'nginx'\n\n\n",
	}

	temp.CleanNewlines()

	if strings.Contains(temp.Content, "\n\n\n") {
		t.Errorf("Expected resulting string not to contain 3 newlines")
	}

	if !strings.Contains(temp.Content, "\n\n") {
		t.Errorf("Expected new string to contain two newlines")
	}
}

func TestCleanNewlines_catches_all(t *testing.T) {
	temp := Template{
		Content: "include_recipe 'nginx'\n\n\ninclude_recipe' 'apache'\n\n\n\n\npackage foo",
	}

	temp.CleanNewlines()
	if strings.Contains(temp.Content, "\n\n\n") {
		t.Errorf("Expected resulting string not to contain 3 newlines")
	}

	if strings.Contains(temp.Content, "\n\n\n\n") {
		t.Errorf("Expected resulting string to not contain 4 newlines")
	}
}
