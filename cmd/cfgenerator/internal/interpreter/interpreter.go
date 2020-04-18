package interpreter

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is an error returned when no interpreter match the givent criteria
	ErrNotFound = errors.New("not found")

	interpreters = make(map[string]BuilderFunc)
)

func init() {
	Register("jsonnet", func() Interpreter { return NewJsonnet() })
	Register("plain", func() Interpreter { return NewGoTemplate() })
}

// BuilderFunc represents a function that initialize a new Interpreter
type BuilderFunc func() Interpreter

// Register registers a new interpreter
func Register(name string, builderFunc BuilderFunc) {
	interpreters[name] = builderFunc
}

// Get builds a new interpreter from its name or return a ErrNotFound error
func Get(name string) (Interpreter, error) {
	builder, found := interpreters[name]
	if !found {
		return nil, fmt.Errorf("%w", ErrNotFound)
	}

	return builder(), nil
}

// Interpreter represents something able to aggregate variables and render templates
type Interpreter interface {
	AddVar(name string, value string)
	Evaluate(tpl string) (string, error)
}
