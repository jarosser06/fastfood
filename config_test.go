package fastfood

import "testing"

func TestNewConfig(t *testing.T) {
	sampleFile := "tests/template.json"

	config, err := NewConfig(sampleFile)

	if err != nil {
		t.Errorf("NewConfig returned unexpected error: %v", err)
	}

	if config.Name != "testcookbook" {
		t.Errorf("Expected name to be testcookbook")
	}

	if config.Providers[0]["provider"] != "database" {
		t.Errorf("Expected first provider to be database")
	}
}
