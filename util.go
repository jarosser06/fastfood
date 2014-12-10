package fastfood

import (
	"bytes"
	"errors"
	"fmt"
	"os"
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

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	} else {
		return false
	}
}

func CapitalizeString(str string) string {
	byteStr := []byte(str)
	capletter := bytes.ToUpper([]byte{byteStr[0]})
	byteStr[0] = capletter[0]

	return string(byteStr)
}