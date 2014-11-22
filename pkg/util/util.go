package util

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
)

// Utility function for appending to a file
func AppendFile(file string, text string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return errors.New(fmt.Sprintf("AppendFile(): %v", err))
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return errors.New(fmt.Sprintf("AppendFile(): %v", err))
	}

	return nil
}

// Assume 3 or more newlines should be compressed
func CollapseNewlines(str string) string {
	reg, _ := regexp.Compile("[\n]{3,}")

	return reg.ReplaceAllString(str, "\n\n")
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	} else {
		return false
	}
}

func IsNodeAttr(str string) bool {
	reg, _ := regexp.Compile(`^node((\[\'([\w_-]+)\'\])+)`)

	return reg.MatchString(str)
}

// Given a struct it will parse for strings
// and modify them with quotes if they are not
// node attributes
func FormatStrings(structT interface{}) {
	s := reflect.ValueOf(structT).Elem()
	for i := 0; i < s.NumField(); i++ {

		field := s.Field(i)
		if field.Kind() == reflect.String {
			if !IsNodeAttr(field.String()) {
				field.SetString(fmt.Sprintf("'%s'", field.String()))
			}
		}
	}
}
