package maputil

import "testing"

func TestMerge(t *testing.T) {
	hMap := map[string]string{
		"override": "higherValue",
		"hkey":     "hvalue",
	}

	lMap := map[string]string{
		"override": "lowerValue",
		"lkey":     "lvalue",
	}

	n := Merge(hMap, lMap)
	for k, v := range hMap {
		if n[k] != v {
			t.Errorf("expected key %s to be %s not %s", k, v, n[k])
		}
	}

	if n["lkey"] != "lvalue" {
		t.Errorf("epxected lkey to be lvalue not %s", n["lkey"])
	}
}
