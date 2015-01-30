package stringutil

import "testing"

func TestCapitalizeString(t *testing.T) {
	testString := "rosser"

	if capped := CapitalizeString(testString); capped != "Rosser" {
		t.Errorf("expected CapitalizeString(rosser) to be Rosser not %s", capped)
	}
}

func TestWrap(t *testing.T) {
	testString := "rosser"

	if r := Wrap(testString, "'"); r != "'rosser'" {
		t.Errorf("expected %s to be wrapped in ' not %s", testString, r)
	}
}
