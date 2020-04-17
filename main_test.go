package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestWithConfigAndSecrets(t *testing.T) {
	input, err := os.Open("./examples/config_and_secrets/config.jsonnet")
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput, err := ioutil.ReadFile("examples/config_and_secrets/expected-config.json")
	if err != nil {
		t.Fatal(err)
	}

	var actualOutput bytes.Buffer
	err = run(input, &actualOutput, []string{
		"./examples/config_and_secrets/volumes/config",
		"./examples/config_and_secrets/volumes/secrets",
	})

	if err != nil {
		t.Fatal(err)
	}

	if string(expectedOutput) != actualOutput.String() {
		t.Fatalf("invalid output\nexpected:\n'%s'\nactual:\n'%s'\n", string(expectedOutput), actualOutput.String())
	}
}
