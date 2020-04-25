package internal

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kdisneur/k8s-cfgenerator/cmd/cfgenerator/internal/interpreter"
	"github.com/kdisneur/k8s-cfgenerator/cmd/cfgenerator/internal/volume"
)

// Generate reads all the volumes to collect the variables and execute the template
func Generate(runtime interpreter.Interpreter, input io.Reader, output io.Writer, volumes []string) error {
	for _, root := range volumes {
		if err := volume.LoadAllVariables(runtime, root); err != nil {
			return fmt.Errorf("can't read volume variables '%s': %v", root, err)
		}
	}

	tpl, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("can't read template: %v", err)
	}

	content, err := runtime.Evaluate(string(tpl))
	if err != nil {
		return fmt.Errorf("can't evaluate template: %v", err)
	}

	fmt.Fprint(output, content)

	return nil
}
