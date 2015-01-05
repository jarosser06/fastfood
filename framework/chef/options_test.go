package chef

import "testing"

var localOptions = `
{
	"berks_dependencies": {
		"couchdb": {
			"git": "git@github.com:jarosser06/couchdb-cookbook.git"
		}
	},
	"dependencies": [
		"couchdb"
	],
	"directories": [
		"recipes"
	],
	"files": {
		"recipes/testfile.rb": "recipes/testfile.rb"
	},
	"partials": [
		"partials/mypartial.rb"
	]
}`

var globalOptions = `
{
	"dependencies": [
		"rackspace_iptables"
	],
	"directories": [
		"templates"
	]
}`

func TestNewOptions(t *testing.T) {
	o, err := NewOptions([]byte(localOptions))

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if _, ok := o.BerksDeps["couchdb"]; !ok {
		t.Errorf("berks dependency couchdb does not exist")
	}

	if len(o.Directories) != 1 {
		t.Errorf("expected directories slice to have 1 element not %d", len(o.Directories))
	}
}

func TestMerge(t *testing.T) {
	l, _ := NewOptions([]byte(localOptions))
	g, _ := NewOptions([]byte(globalOptions))

	merged := Merge(g, l)

	if len(merged.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies after merge, got %d", len(merged.Dependencies))
	}

	if len(merged.BerksDeps) != 1 {
		t.Errorf("expected 1 berks dependency, got %d", len(merged.BerksDeps))
	}
}
