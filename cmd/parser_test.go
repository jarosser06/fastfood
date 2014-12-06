package cmd

import (
	"testing"

	"github.com/jarosser06/fastfood/provider"
)

func TestNewCommand(t *testing.T) {
	if c := NewCommand("testcommand"); c.Name != "testcommand" {
		t.Errorf("Expected new command with name testcommand")
	}
}

func TestParseCommandsFromFile(t *testing.T) {
	sampleFile := "../samples/manifest.json"

	commands := ParseCommandsFromFile(sampleFile)

	if len(commands) != 2 {
		t.Errorf("Expected the length of the commands array to be 2")
	}

	dbCmdExists := false
	for _, command := range commands {
		if command.Name == "db" {
			dbCmdExists = true
		}
	}

	if !dbCmdExists {
		t.Errorf("Expected one of the parsed commands to match the name 'db'")
	}
}

func TestParseProviderFromFile(t *testing.T) {
	ckbk := provider.Cookbook{
		Path:         "/tmp/testcookbook",
		Name:         "testcookbook",
		Dependencies: []string{"apt"},
	}

	sampleManifest := "../samples/database/manifest.json"

	p := ParseProviderFromFile(ckbk, sampleManifest)

	if _, ok := p.Types["mysql_master"]; !ok {
		t.Errorf("Expected on of the available types to be mysql_master\n")
	}
}
