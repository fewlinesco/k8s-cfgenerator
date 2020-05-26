package internal_test

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/fewlinesco/k8s-cfgenerator/cmd/cfgenerator/internal"
	"github.com/fewlinesco/k8s-cfgenerator/cmd/cfgenerator/internal/interpreter"
)

func getRuntime(t *testing.T, name string) interpreter.Interpreter {
	runtime, found := interpreter.Get(name)
	if !found {
		t.Fatalf("can't get interpreter")
	}

	return runtime
}

func openInput(t *testing.T, path string) io.Reader {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("can't open input file: %v", err)
	}
	return f
}

func readExpectedOutput(t *testing.T, path string) string {
	s, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("can't read expected output file: %v", err)
	}
	return string(s)
}

func TestValidTemplates(t *testing.T) {
	tcs := []struct {
		Name               string
		RuntimeName        string
		InputPath          string
		Volumes            []string
		ExpectedOutputPath string
	}{
		{
			RuntimeName: "jsonnet",
			InputPath:   "../examples/jsonnet/config.jsonnet",
			Volumes: []string{
				"../examples/jsonnet/volumes/config",
				"../examples/jsonnet/volumes/secrets",
			},
			ExpectedOutputPath: "../examples/jsonnet/expected-config.json",
		},
		{
			RuntimeName: "plain",
			InputPath:   "../examples/plain/config.conf.tpl",
			Volumes: []string{
				"../examples/jsonnet/volumes/config",
				"../examples/jsonnet/volumes/secrets",
			},
			ExpectedOutputPath: "../examples/plain/expected-config.conf",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.RuntimeName, func(t *testing.T) {
			runtime := getRuntime(t, tc.RuntimeName)
			input := openInput(t, tc.InputPath)
			expectedOutput := readExpectedOutput(t, tc.ExpectedOutputPath)
			var output strings.Builder

			if err := internal.Generate(runtime, input, &output, tc.Volumes); err != nil {
				t.Fatal(err)
			}

			if expectedOutput != output.String() {
				t.Fatalf("invalid output\nexpected:\n'%s'\nactual:\n'%s'\n", expectedOutput, output.String())
			}
		})
	}

}
