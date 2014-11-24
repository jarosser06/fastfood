package application

import (
	"testing"

	"github.com/jarosser06/fastfood/provider/cookbook"
)

func FakeCookbook() cookbook.Cookbook {
	return cookbook.NewCookbook(
		"/tmp",
		"testcookbook",
	)
}

func TestNew(t *testing.T) {
	app := New("testapp", FakeCookbook())

	if app.Name != "testapp" {
		t.Errorf("Expected new application to have name testapp")
	}
}

func TestPath(t *testing.T) {
	app := New("testapp", FakeCookbook())

	if app.Path() != "/var/www/testapp" {
		t.Errorf("Expected Path() to return /var/www/testapp not %v", app.Path())
	}
}
