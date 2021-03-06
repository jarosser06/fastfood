package fastfood

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/template"
)

const (
	leftDelim  = "|{"
	rightDelim = "}|"
)

type Template struct {
	Content string
	Raw     string
}

// Given a name, an interface, and the template content returns a Template
func NewTemplate(name string, values interface{}, content []string) (*Template, error) {
	temp := template.New(name)

	temp.Delims(leftDelim, rightDelim)
	for _, cont := range content {
		_, err := temp.Parse(cont)

		if err != nil {
			err = fmt.Errorf("NewTemplate() error parsing content %v", err)
			return &Template{}, err
		}
	}

	//TODO: Print errors during template parsing
	var buffer bytes.Buffer
	err := temp.Execute(&buffer, values)
	if err != nil {
		return &Template{}, err
	}
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

func (t *Template) CleanNewlines() {
	reg, _ := regexp.Compile("[\n]{3,}")
	t.Content = reg.ReplaceAllString(t.Content, "\n\n")
}
