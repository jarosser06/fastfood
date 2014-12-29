package stringutil

import "bytes"

func CapitalizeString(str string) string {
	byteStr := []byte(str)
	capletter := bytes.ToUpper([]byte{byteStr[0]})
	byteStr[0] = capletter[0]

	return string(byteStr)
}
