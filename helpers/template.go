package helpers

import (
	"fmt"

	"github.com/jarosser06/fastfood/util"
)

type Template struct {
}

func (t *Template) QString(str string) string {
	if util.IsNodeAttr(str) {
		return str
	} else {
		return fmt.Sprintf("'%s'", str)
	}
}
