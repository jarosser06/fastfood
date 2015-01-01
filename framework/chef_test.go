package framework

import (
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
}

func TestcapitalizeOptions(t *testing.T) {

}
