package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestWithJsonnet(t *testing.T) {
	input, err := os.Open("./examples/jsonnet/config.jsonnet")
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput, err := ioutil.ReadFile("examples/jsonnet/expected-config.json")
	if err != nil {
		t.Fatal(err)
	}

	var actualOutput bytes.Buffer
	err = run("jsonnet", input, &actualOutput, []string{
		"./examples/jsonnet/volumes/config",
		"./examples/jsonnet/volumes/secrets",
	})

	if err != nil {
		t.Fatal(err)
	}

	if string(expectedOutput) != actualOutput.String() {
		t.Fatalf("invalid output\nexpected:\n'%s'\nactual:\n'%s'\n", string(expectedOutput), actualOutput.String())
	}
}

func TestWithGotemplate(t *testing.T) {
	input, err := os.Open("./examples/plain/config.conf.tpl")
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput, err := ioutil.ReadFile("examples/plain/expected-config.conf")
	if err != nil {
		t.Fatal(err)
	}

	var actualOutput bytes.Buffer
	err = run("plain", input, &actualOutput, []string{
		"./examples/jsonnet/volumes/config",
		"./examples/jsonnet/volumes/secrets",
	})

	if err != nil {
		t.Fatal(err)
	}

	if string(expectedOutput) != actualOutput.String() {
		t.Fatalf("invalid output\nexpected:\n'%s'\nactual:\n'%s'\n", string(expectedOutput), actualOutput.String())
	}
}
