package util

import (
	"errors"
	"fmt"
	"os"
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
