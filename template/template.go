package template

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/jarosser06/fastfood/util"
)

const (
	leftDelim  = "|{"
	rightDelim = "}|"
)

type Template struct {
	Content string
}

// Given a name, an interface, and the template content returns a Template
func NewTemplate(name string, values interface{}, content ...string) (*Template, error) {
	temp := template.New(name)

	temp.Delims(leftDelim, rightDelim)
	for _, cont := range content {
		_, err := temp.Parse(cont)

		if err != nil {
			err = errors.New(fmt.Sprintf("NewTemplate() error parsing content %v", err))
			return &Template{}, err
		}
	}

	var buffer bytes.Buffer
	temp.Execute(&buffer, values)
	return &Template{Content: buffer.String()}, nil
}

// Flush the template to a file
func (t *Template) Flush(fileName string) error {
	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		return err
	}

	io.WriteString(f, t.Content)
	return nil
}

// Helper calls util.CollapseNewlines on template string
func (t *Template) CleanNewlines() {
	t.Content = util.CollapseNewlines(t.Content)
}
