package interpreter

import (
	"fmt"
	"strings"
	"text/template"
)

// Plain represents the Go Template interpreter
type Plain map[string]string

// NewPlain builds a new Go Template interpreter
func NewPlain() *Plain {
	return &Plain{}
}

// AddVar stores a new variable
func (g Plain) AddVar(name string, value string) {
	g[name] = value
}

// Evaluate executes the template with all the variable previously stored accessible
func (g Plain) Evaluate(tpl string) (string, error) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("can't parse plain template: %v", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, map[string]string(g)); err != nil {
		return "", fmt.Errorf("can't evaluate plain template: %v", err)
	}

	return buf.String(), nil
}
