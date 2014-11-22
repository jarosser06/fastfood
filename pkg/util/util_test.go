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

func TestIsNodeAttr(t *testing.T) {
	testMatch := "node['cookbook']['attr']"
	testNotMatch := "something[diff]"

	if !IsNodeAttr(testMatch) {
		t.Errorf("Expected %s to return true", testMatch)
	}

	if IsNodeAttr(testNotMatch) {
		t.Errorf("Expected %s to return false", testNotMatch)
	}
}
