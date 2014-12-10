package fastfood

import "testing"

func TestQString(t *testing.T) {
	s := Helpers{}

	testMatch := "somestring"
	testNotMatch := "node['cookbook']['attr']"

	if res := s.QString(testMatch); res != "'somestring'" {
		t.Errorf("Expected 'somestring' but recieved %s", res)
	}

	if res := s.QString(testNotMatch); res != testNotMatch {
		t.Errorf("Expected %s but recieved %s", res)
	}
}

func TestIsNodeAttr(t *testing.T) {
	testMatch := "node['cookbook']['attr']"
	testNotMatch := "something[diff]"

	h := Helpers{}

	if !h.IsNodeAttr(testMatch) {
		t.Errorf("Expected %s to return true", testMatch)
	}

	if h.IsNodeAttr(testNotMatch) {
		t.Errorf("Expected %s to return false", testNotMatch)
	}
}

func TestIsChefVar(t *testing.T) {
	testMatch := "node.chef_environment"
	testNotMatch := "nodesomething"

	h := Helpers{}

	if !h.IsChefVar(testMatch) {
		t.Errorf("Expected %s to return true", testMatch)
	}

	if h.IsChefVar(testNotMatch) {
		t.Errorf("Expected %s to return false", testNotMatch)
	}
}
