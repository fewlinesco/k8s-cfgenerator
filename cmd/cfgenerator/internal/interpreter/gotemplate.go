package interpreter

import (
	"fmt"
	"strings"
	"text/template"
)

// GoTemplate represents the Go Template interpreter
type GoTemplate map[string]string

// NewGoTemplate builds a new Go Template interpreter
func NewGoTemplate() *GoTemplate {
	return &GoTemplate{}
}

// AddVar stores a new variable
func (g GoTemplate) AddVar(name string, value string) {
	g[name] = value
}

// Evaluate executes the template with all the variable previously stored accessible
func (g GoTemplate) Evaluate(tpl string) (string, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("can't parse gotpl template: %v", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, map[string]string(g)); err != nil {
		return "", fmt.Errorf("can't evaluate gotpl template: %v", err)
	}

	return buf.String(), nil
}
