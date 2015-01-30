package stringutil

import (
	"bytes"
	"fmt"
)

func CapitalizeString(str string) string {
	byteStr := []byte(str)
	capletter := bytes.ToUpper([]byte{byteStr[0]})
	byteStr[0] = capletter[0]

	return string(byteStr)
}

func Wrap(str string, wrapper string) string {
	return fmt.Sprintf("%s%s%s", wrapper, str, wrapper)
}
