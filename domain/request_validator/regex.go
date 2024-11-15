package requestvalidator

import (
	"fmt"
	"strings"
)

const (
	alphaRegex   = "a-zA-Z"
	numericRegex = "0-9"
	specialRegex = `[!@#$%^&*(),.?":{}|<>]`
	numberRegex  = `^-?\d+(\.\d+)?$`
	booleanRegex = `^(true|false)$`
)

func (t *textValue) minMax(min int, max int) string {
	return fmt.Sprintf("{%d,%d}", min, max)
}

func (t *textValue) in(in []any) string {
	if len(in) == 0 {
		return ""
	}
	var patterns []string
	for _, v := range in {
		patterns = append(patterns, fmt.Sprintf("^%v$", v))
	}
	return fmt.Sprintf("(?:%s)", strings.Join(patterns, "|"))
}

func (t *textValue) notIn(in []any) string {
	if len(in) == 0 {
		return ""
	}
	var patterns []string
	for _, v := range in {
		patterns = append(patterns, fmt.Sprintf("^(?!%v$)", v))
	}
	return fmt.Sprintf("(?:%s)", strings.Join(patterns, ""))
}

func (n *numberValue) minMax(min int, max int) string {
	var conditions []string
	if min != 0 {
		conditions = append(conditions, fmt.Sprintf("(?:%d <= n)", min))
	}
	if max != 0 {
		conditions = append(conditions, fmt.Sprintf("(?:n <= %d)", max))
	}
	return strings.Join(conditions, " && ")
}

func (n *numberValue) in(in []any) string {
	if len(in) == 0 {
		return ""
	}
	var patterns []string
	for _, v := range in {
		patterns = append(patterns, fmt.Sprintf("^%v$", v))
	}
	return fmt.Sprintf("(?:%s)", strings.Join(patterns, "|"))
}

func (n *numberValue) notIn(in []any) string {
	if len(in) == 0 {
		return ""
	}
	var patterns []string
	for _, v := range in {
		patterns = append(patterns, fmt.Sprintf("^(?!%v$)", v))
	}
	return fmt.Sprintf("(?:%s)", strings.Join(patterns, ""))
}
