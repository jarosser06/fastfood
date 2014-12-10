package fastfood

import (
	"fmt"
	"regexp"
)

type Helpers struct {
}

func (h *Helpers) QString(str string) string {
	if h.IsNodeAttr(str) || h.IsChefVar(str) {
		return str
	} else {
		return fmt.Sprintf("'%s'", str)
	}
}

func (h *Helpers) IsNodeAttr(str string) bool {
	reg, _ := regexp.Compile(`^node((\[\'([\w_-]+)\'\])+)`)

	return reg.MatchString(str)
}

func (h *Helpers) IsChefVar(str string) bool {
	reg, _ := regexp.Compile(`^node\.([\w_-]+)`)

	return reg.MatchString(str)
}
