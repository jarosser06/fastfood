package provider

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/GeertJohan/go.rice"
	"github.com/jarosser06/fastfood/util"
)

type Provider interface {
	Exist() bool
	Files() map[string]string
	Partials() []string
	TemplateBox() string
}

func GenFiles(p Provider) error {

	templateBox, _ := rice.FindBox(
		fmt.Sprintf("templates/%s", p.TemplateBox()),
	)

	for cookbookFile, templateFile := range p.Files() {
		tmpStr, _ := templateBox.String(templateFile)

		var loadedPartials []string
		for _, partial := range p.Partials() {
			partialStr, _ := templateBox.String(fmt.Sprintf("partials/%s", partial))

			loadedPartials := append(loadedPartials, partialStr)

		}
	}

}

func GenDirs(dirs []string, cookbookPath string) error {
	for _, dir := range dirs {
		fullPath := path.Join(cookbookPath, dir)

		if !util.FileExist(fullPath) {
			err := os.MkdirAll(fullPath, 0755)

			if err != nil {
				return errors.New(fmt.Sprintf("GenDirs(): %v", err))
			}
		}
	}

	return nil
}
