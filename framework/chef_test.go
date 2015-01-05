package framework

import (
	"os"
	"path"
	"testing"

	"github.com/jarosser06/fastfood"
)

var fopts = fastfood.FrameworkOptions{
	Destination: "/tmp",
	BaseFiles: []string{
		"Berksfile",
		"CHANGELOG.md",
		"Gemfile",
		"metadata.rb",
		"README.md",
		"recipes/default.rb",
		"test/unit/spec/default_spec.rb",
		"test/unit/spec/spec_helper.rb",
	},
	BaseDirs: []string{
		"attributes",
		"files",
		"libraries",
		"providers",
		"recipes",
		"resources",
		"templates",
		"test/unit/spec",
	},
	Name:        "testcookbook",
	TemplateDir: "../tests/templatepack/frameworks/chef",
}

func TestChefInit(t *testing.T) {
	c := Chef{}
	err := c.Init(fopts)
	if err != nil {
		t.Errorf("Init() returned unexepected error %v")
	}

	if c.cookbook.Name != fopts.Name {
		t.Errorf("expecting cookbook name %s not %s", fopts.Name, c.cookbook.Name)
	}

	cPath := path.Join(fopts.Destination, fopts.Name)
	if c.cookbook.Path != cPath {
		t.Errorf("expected cookbook path to be %s not %s", cPath, c.cookbook.Path)
	}
}

func TestChefGenerateBase(t *testing.T) {
	c := Chef{}
	c.Init(fopts)

	updatedFiles, err := c.GenerateBase()
	defer os.RemoveAll("/tmp/testcookbook")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(updatedFiles) != len(fopts.BaseFiles) {
		t.Errorf("expected all base files to be updated")
	}
}

func TestChefGenerateStencil(t *testing.T) {
	sset, _ := fastfood.NewStencilSet("../tests/templatepack/stencils/database/manifest.json")
	c := Chef{}
	c.Init(fopts)
	c.GenerateBase()
	defer os.RemoveAll("/tmp/testcookbook")

	opts := map[string]string{
		"name":    "my_mysql",
		"openfor": "myapp",
	}

	updatedFiles, err := c.GenerateStencil("mysql_master", sset, opts)
	if err != nil {
		t.Errorf("unexpected error occured generating stencil: %v", err)
	}

	if len(updatedFiles) != 2 {
		t.Errorf("expected 2 files to be updated not %i", len(updatedFiles))
	}
}

func TestChefGenerateStencil_forceworks(t *testing.T) {
	sset, _ := fastfood.NewStencilSet("../tests/templatepack/stencils/database/manifest.json")
	c := Chef{}
	// Turn force on
	fopts.Force = true
	c.Init(fopts)
	c.GenerateBase()
	defer os.RemoveAll("/tmp/testcookbook")

	opts := map[string]string{
		"name":    "my_mysql",
		"openfor": "myapp",
	}

	// Expect that all files will be updated every time with force on
	for i := 0; i < 2; i++ {
		updatedFiles, err := c.GenerateStencil("mysql_master", sset, opts)
		if err != nil {
			t.Errorf("unexpected error occured generating stencil: %v", err)
		}

		if len(updatedFiles) != 2 {
			t.Errorf("expected 2 files to be updated not %i", len(updatedFiles))
		}
	}
}

func TestCapitalizeOptions(t *testing.T) {
	o := CapitalizeOptions(map[string]string{"option1": "value1"})

	if len(o) != 1 {
		t.Errorf("expected returned map to have only 1 item")
	}

	if _, ok := o["Option1"]; !ok {
		t.Errorf("expected option1 to become Option1")
	}
}
