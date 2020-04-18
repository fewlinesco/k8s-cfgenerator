package interpreter

import (
	"fmt"

	"github.com/google/go-jsonnet"
)

// Jsonnet represents the JSONNET interpreter
type Jsonnet struct {
	vm *jsonnet.VM
}

// NewJsonnet builds a new JSONNET interpreter
func NewJsonnet() *Jsonnet {
	return &Jsonnet{vm: jsonnet.MakeVM()}
}

// AddVar stores a new variable as ExtVar
func (j *Jsonnet) AddVar(name string, value string) {
	j.vm.ExtVar(name, value)
}

// Evaluate executes the template with all the variable previously stored accessible using std.extVar
func (j *Jsonnet) Evaluate(tpl string) (string, error) {
	json, err := j.vm.EvaluateSnippet("", string(tpl))
	if err != nil {
		return "", fmt.Errorf("can't evaluate jsonnet template: %v", err)
	}

	return json, nil
}
