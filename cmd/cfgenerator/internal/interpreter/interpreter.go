package interpreter

import (
	"errors"
)

var (
	// ErrNotFound is an error returned when no interpreter match the givent criteria
	ErrNotFound = errors.New("not found")

	interpreters = make(map[string]BuilderFunc)
)

func init() {
	Register("jsonnet", func() Interpreter { return NewJsonnet() })
	Register("plain", func() Interpreter { return NewPlain() })
}

// BuilderFunc represents a function that initialize a new Interpreter
type BuilderFunc func() Interpreter

// Register registers a new interpreter
func Register(name string, builderFunc BuilderFunc) {
	interpreters[name] = builderFunc
}

// Get builds a new interpreter from its name and return a boolean indicating wether the interpreter
// has been found
func Get(name string) (Interpreter, bool) {
	builder, found := interpreters[name]
	if !found {
		return nil, false
	}

	return builder(), true
}

// Interpreter represents something able to aggregate variables and render templates
type Interpreter interface {
	AddVar(name string, value string)
	Evaluate(tpl string) (string, error)
}
