package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

const usageFmt = `Synopsis

	%[1]s [volume-paths ...]

Description

	Reads a JSONNET from STDIN and output a JSON in STDOUT.

	It reads all files present in 'the volume-paths' folders and for each of
	these files, sets the file name as a JSONNET extVar and the content of
	the file as the value of the extVar.

	These JSONNET extVar variables are available in the JSONNET template
	received from STDIN.

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
	flag.Usage = func() { fmt.Fprintf(flag.CommandLine.Output(), usageFmt, filepath.Base(os.Args[0])) }
	flag.Parse()

	if err := run(os.Stdin, os.Stdout, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(input *os.File, output io.Writer, volumes []string) error {
	stat, err := input.Stat()
	if err != nil {
		return fmt.Errorf("can't read from input file: %v", err)
	}

	if stat.Size() <= 0 {
		return fmt.Errorf("empty input file")
	}

	var buf bytes.Buffer
	vm := jsonnet.MakeVM()
	for _, root := range volumes {
		err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if p == root && info.IsDir() {
				return nil
			}

			if info.IsDir() {
				return filepath.SkipDir
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

			vm.ExtVar(extVarName, extVarValue)

			return nil
		})

		if err != nil {
			return fmt.Errorf("can't read external variables: %v", err)
		}
	}

	tpl, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("can't read jsonnet template from STDIN: %v", err)
	}

	json, err := vm.EvaluateSnippet("", string(tpl))
	if err != nil {
		return fmt.Errorf("can't evaluate jsonnet template: %v", err)
	}

	fmt.Fprint(output, json)

	return nil
}
