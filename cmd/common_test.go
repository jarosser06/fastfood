package cmd

import "testing"

func TestNewManifest(t *testing.T) {
	sampleFile := "../samples/manifest.json"

	manifest, err := NewManifest(sampleFile)

	if err != nil {
		t.Errorf("Did not expect error %v", err)
	}

	if len(manifest.Providers) != 2 {
		t.Errorf("Expected the length of the commands array to be 2")
	}

	dbCmdExists := false
	for _, provider := range manifest.Providers {
		if provider.Name == "database" {
			dbCmdExists = true
		}
	}

	if !dbCmdExists {
		t.Errorf("Expected one of the parsed commands to match the name 'db'")
	}

	if len(manifest.Cookbook.Files) == 0 {
		t.Errorf("Expected more than 0 cookbook files")
	}
}

func TestMapArgs(t *testing.T) {
	args := []string{"name:testing"}

	mappedArgs := MapArgs(args)

	if val, ok := mappedArgs["name"]; !ok && val == "testing" {
		t.Errorf("Expected a map with a key of name and value of testing")
	}
}
