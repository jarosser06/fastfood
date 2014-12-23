package cmd

import (
	"errors"
	"sync"
)

type Builder struct {
	Common
	WaitGroup sync.WaitGroup
}

func (b *Builder) Run(args []string) error {
	/*
		manifest, err := fastfood.NewManifest(path.Join(templatePack, "manifest.json"))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		args = cmdFlags.Args()
	*/
	return errors.New("not implemented")
}

func (b *Builder) Description() string {
	return "Creates a cookbook w/ providers from a config file"
}

func (b *Builder) Help() string {
	return "fastfood build [config_file]"
}
