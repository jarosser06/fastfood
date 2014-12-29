package stringutil

import "testing"

func TestCapitalizeString(t *testing.T) {
	testString := "rosser"

	if capped := CapitalizeString(testString); capped != "Rosser" {
		t.Errorf("expected CapitalizeString(rosser) to be Rosser not %s", capped)
	}
}
