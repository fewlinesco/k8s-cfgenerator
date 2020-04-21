package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kdisneur/k8s-cfgenerator/cmd/cfgenerator/internal/interpreter"
	"github.com/kdisneur/k8s-cfgenerator/cmd/cfgenerator/internal/volume"
)

const usageFmt = `Synopsis

	%[1]s [-interpreter=plain|jsonnet] [volume-paths ...]

Description

	Reads a content (plain text or JSONNET) from STDIN and output the result
	to STDOUT (as a JSON or plain text).

	It reads all files present in 'the volume-paths' folders and for each of
	these files, sets the file name as variable name and the content of the
	file as value.

Flags

	-interpreter=plain|jsonnet
	   When plain, interprets the input as plain text and use gotpl as
	   variable system.

	   When jsonnet, interprets the input as JSONNET and use extVar as
	   variable system.

	   By default it is set to jsonnet

Arguments

	[volume-paths ...]
	   a list of folder or files.

	   When file: the content of the file will be loaded and set in a JSONNET
	   extVar named with the file name.

	   When folder: the content of each of the file of the folder will be
	   loaded and set in a JSONNET extVar named with the file name.
	   The script doesn't load files in sub folders.

Examples

	1. read all files in /data/configmap and /data/secrets and use their name
	   as JSONNET extvar. Then evaluates STDIN and generate a JSON in STDOUT

	   $> %[1]s /data/configmap /data/secrets <<EOF
	   {
	   	api: {
	   		address: '0.0.0.0:' + std.extVar('API_PORT'),
	   	},
	   	database: {
	   		password: std.extVar("DATABASE_PASSWORD"),
	   		username: std.extVar("DATABASE_USERNAME"),
	   	},
	   }
	   EOF

	2. read all files in /data/configmap and /data/secrets and use their name
	   as JSONNET extvar. Then evaluates /app/config.jsonnet and generate a
	   JSON in /app/config.json

	   $> %[1]s /data/configmap /data/secrets < /app/confg.jsonnet > /app/config.json

`

func main() {
	var cfg = struct {
		InterpreterName string
	}{
		InterpreterName: "jsonnet",
	}

	flag.Usage = func() { fmt.Fprintf(flag.CommandLine.Output(), usageFmt, filepath.Base(os.Args[0])) }
	flag.StringVar(&cfg.InterpreterName, "interpreter", cfg.InterpreterName, "")

	flag.Parse()

	if err := run(cfg.InterpreterName, os.Stdin, os.Stdout, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(interpreterName string, input *os.File, output io.Writer, volumes []string) error {
	runtime, found := interpreter.Get(interpreterName)
	if !found {
		return fmt.Errorf("unsupported interpreter '%s'", interpreterName)
	}

	stat, err := input.Stat()
	if err != nil {
		return fmt.Errorf("can't read from input file: %v", err)
	}

	if stat.Size() <= 0 {
		return fmt.Errorf("empty input file")
	}

	for _, root := range volumes {
		if err := volume.LoadAllVariables(runtime, root); err != nil {
			return fmt.Errorf("can't read volume variables '%s': %v", root, err)
		}
	}

	tpl, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("can't read template from STDIN: %v", err)
	}

	content, err := runtime.Evaluate(string(tpl))
	if err != nil {
		return fmt.Errorf("can't evaluate template: %v", err)
	}

	fmt.Fprint(output, content)

	return nil
}
