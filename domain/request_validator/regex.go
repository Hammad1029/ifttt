package requestvalidator

import (
	"fmt"
	"strings"
)

const (
	alphaRegex   = "a-zA-Z *"
	numericRegex = "0-9"
	specialRegex = `!@#$%^&*(),.?":{}|<>`
	booleanRegex = `^(true|false)$`
)

func (t *textValue) minMax() string {
	if t.Minimum > 0 || t.Maximum > 0 {
		return fmt.Sprintf("{%d,%d}", t.Minimum, t.Maximum)
	} else {
		return "+"
	}
}

func (t *textValue) in() string {
	if len(t.In) == 0 {
		return ""
	}
	var patterns []string
	for _, v := range t.In {
		patterns = append(patterns, fmt.Sprint(v))
	}
	return fmt.Sprintf("^(%s)$", strings.Join(patterns, "|"))
}

func (n *numberValue) in() string {
	if len(n.In) == 0 {
		return ""
	}
	var patterns []string
	for _, v := range n.In {
		patterns = append(patterns, fmt.Sprint(v))
	}
	return fmt.Sprintf("^(%s)$", strings.Join(patterns, "|"))
}
