package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Command struct {
	Name          string `json:"name"`
	Manifest      string `json:"manifest"`
	Help          string `json:"help"`
	templatesPath string
}

func NewCommand(name string) Command {
	return Command{Name: name}
}

func ParseCommandsFromFile(path string) map[string]Command {

	cmdsStruct := struct {
		Commands map[string]Command
	}{}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file %s", path))
	}

	err = json.Unmarshal(f, &cmdsStruct)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse json: %v", err))
	}

	return cmdsStruct.Commands
}

func NewCommandFromFile(name string, path string) Command {
	c := Command{Name: name}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file %s", path))
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse json: %v", err))
	}

	return c
}
