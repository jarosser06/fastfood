package cmd

import "testing"

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
