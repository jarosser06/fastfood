package helpers

import "testing"

func TestQString(t *testing.T) {
	s := Template{}

	testMatch := "somestring"
	testNotMatch := "node['cookbook']['attr']"

	if res := s.QString(testMatch); res != "'somestring'" {
		t.Errorf("Expected 'somestring' but recieved %s", res)
	}

	if res := s.QString(testNotMatch); res != testNotMatch {
		t.Errorf("Expected %s but recieved %s", res)
	}
}
