package fileutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Utility function for appending to a file
func AppendFile(file string, text string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return fmt.Errorf("AppendFile(): %v", err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return fmt.Errorf("AppendFile(): %v", err)
	}

	return nil
}

func Copy(source string, dest string) error {
	if FileExist(dest) {
		return errors.New("file already exists")
	} else {
		sourceBytes, err := ioutil.ReadFile(source)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(dest, sourceBytes, 0644)
		if err != nil {
			return err
		}

		return nil
	}
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	} else {
		return false
	}
}
