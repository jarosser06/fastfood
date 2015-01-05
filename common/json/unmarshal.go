/*
 * Thanks to github.com/mitchellh for already figuring this out
 */

package json

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Unmarshal(data []byte, i interface{}) error {
	err := json.Unmarshal(data, i)
	if err != nil {
		serr, ok := err.(*json.SyntaxError)
		if !ok {
			return err
		}

		nline := []byte{'\x0a'}

		start := bytes.LastIndex(data[:serr.Offset], nline) + 1
		end := len(data)
		if idx := bytes.Index(data[start:], nline); idx >= 0 {
			end = start + idx
		}

		line := bytes.Count(data[:start], nline) + 1
		pos := int(serr.Offset) - start - 1

		return fmt.Errorf("Error in line %d, char %d: %s\n%s",
			line, pos, serr, data[start:end])
	}

	return nil
}
