package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/jarosser06/fastfood/cmd"
)

const ffVersion = "0.1.2"

type Command interface {
	Help() string
	Run([]string) error
	Description() string
	LoadManifest() error
	SetCookbookPath(string)
	SetTemplatePack(string)
}

func main() {
	app := cli.NewApp()
	app.Name = "fastfood"
	app.Usage = "Generates chef things from templates"
	app.Version = ffVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "template-pack, p",
			Value: func() string {
				home := os.Getenv("HOME")
				return path.Join(home, "fastfood")
			}(),
			Usage:  "path to the template pack",
			EnvVar: cmd.EnvTempPack,
		},
		cli.StringFlag{
			Name: "cookbooks-path, c",
			Value: func() string {
				home := os.Getenv("HOME")
				return path.Join(home, "cookbooks")
			}(),
			Usage:  "directory to store the new cookbook",
			EnvVar: cmd.EnvCookbooksPath,
		},
	}

	app.Commands = setUpCommands()

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("Error executing command: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func setUpCommands() []cli.Command {
	// Create a new map of commands
	cmds := map[string]Command{
		"gen":   &cmd.Generator{},
		"new":   &cmd.Creator{},
		"build": &cmd.Builder{},
	}

	commands := []cli.Command{
		{
			Name:        "gen",
			Description: cmds["gen"].Description(),
			Usage:       cmds["gen"].Help(),
			Action: func(c *cli.Context) {
				cmds["gen"].SetTemplatePack(c.GlobalString("template-pack"))
				cmds["gen"].SetCookbookPath(c.GlobalString("cookbooks-path"))

				if err := cmds["gen"].LoadManifest(); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					if err := cmds["gen"].Run(c.Args()); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}
			},
		},
		{
			Name:        "new",
			Description: cmds["new"].Description(),
			Usage:       cmds["new"].Help(),
			Action: func(c *cli.Context) {
				cmds["new"].SetTemplatePack(c.GlobalString("template-pack"))
				cmds["new"].SetCookbookPath(c.GlobalString("cookbooks-path"))

				if err := cmds["new"].LoadManifest(); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					if err := cmds["new"].Run(c.Args()); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}
			},
		},
		{
			Name:        "build",
			Description: cmds["build"].Description(),
			Usage:       cmds["build"].Help(),
			Action: func(c *cli.Context) {
				cmds["build"].SetTemplatePack(c.GlobalString("template-pack"))
				cmds["build"].SetCookbookPath(c.GlobalString("cookbooks-path"))

				if err := cmds["build"].LoadManifest(); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					if err := cmds["build"].Run(c.Args()); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}
			},
		},
	}

	return commands
}
