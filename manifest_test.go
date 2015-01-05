package fastfood

import "testing"

var sampleFile = "tests/templatepack/manifest.json"

func TestNewManifest(t *testing.T) {
	manifest, err := NewManifest(sampleFile)

	if err != nil {
		t.Errorf("Did not expect error %v", err)
	}

	if len(manifest.StencilSets) != 2 {
		t.Errorf("Expected the length of the commands array to be 2")
	}

	eflag := false
	for _, s := range manifest.StencilSets {
		if s.Name == "database" {
			eflag = true
		}
	}

	if !eflag {
		t.Errorf("Expected one of the parsed commands to match the name 'db'")
	}

	if len(manifest.Base.Files) == 0 {
		t.Errorf("Expected more than 0 cookbook files")
	}
}

func TestValidManifest(t *testing.T) {
	m, _ := NewManifest(sampleFile)

	if ok, err := m.Valid(); !ok {
		t.Errorf("%s not a valid manifest %v", sampleFile, err)
	}
}
