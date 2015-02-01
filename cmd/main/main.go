package main

import (
	"fmt"
	"os"

	"github.com/jarosser06/fastfood/cmd"
	"github.com/mitchellh/cli"
)

const ffVersion = "0.3.0"

func main() {
	c := cli.NewCLI("fastfood", ffVersion)

	// Pass remaining args minus the program name
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"gen": func() (cli.Command, error) {
			return &cmd.Generator{}, nil
		},
		"new": func() (cli.Command, error) {
			return &cmd.Creator{}, nil
		},
		"build": func() (cli.Command, error) {
			return &cmd.Builder{}, nil
		},
	}

	exitstatus, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(exitstatus)
}
