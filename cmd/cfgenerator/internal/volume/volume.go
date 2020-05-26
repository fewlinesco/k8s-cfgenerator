package volume

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fewlinesco/k8s-cfgenerator/cmd/cfgenerator/internal/interpreter"
)

// LoadAllVariables reads all the files in the root folder (or just the root file if it's
// a file) and load all the variables into the runtime.
//
// The name of each file define the variable name and its content the value
func LoadAllVariables(runtime interpreter.Interpreter, root string) error {
	var buf bytes.Buffer

	return filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p == root && info.IsDir() {
			return nil
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		if strings.HasPrefix(info.Name(), ".") {
			// Skip hidden files. It's mostly due to the way ConfigMap and Secrets are
			// mounted on Kubernetes
			return nil
		}

		file, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("can't open file %s: %v", p, err)
		}
		defer file.Close()

		buf.Reset()
		if _, err := io.Copy(&buf, file); err != nil {
			return fmt.Errorf("can't read external variable: %s", p)
		}

		extVarName := filepath.Base(p)
		extVarValue := strings.TrimSpace(buf.String())

		runtime.AddVar(extVarName, extVarValue)

		return nil
	})
}
