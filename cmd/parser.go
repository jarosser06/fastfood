package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jarosser06/fastfood/provider"
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

//TODO: This really shouldn't belong to the cmd parser
func ParseProviderFromFile(ckbk provider.Cookbook, path string) provider.Provider {
	provider := provider.NewProvider(ckbk)

	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file %s", path))
	}

	err = json.Unmarshal(f, &provider)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse json: %v", err))
	}

	return provider
}
