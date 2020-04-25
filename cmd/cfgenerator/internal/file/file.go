package file

import (
	"fmt"
	"os"
)

// OpenInput opens the file for reading and ensures it's not empty.
// If path is `-` it reads from STDIN
func OpenInput(path string) (*os.File, error) {
	var input *os.File

	switch path {
	case "-":
		input = os.Stdin
	default:
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("can't open file")
		}

		input = f
	}

	stat, err := input.Stat()
	if err != nil {
		return input, fmt.Errorf("can't read from file: %v", err)
	}

	if stat.Size() <= 0 {
		return input, fmt.Errorf("empty file")
	}

	return input, nil
}

// OpenOutput opens the file for writing.
// If path is `-` it writes to STDOUT
func OpenOutput(path string) (*os.File, error) {
	switch path {
	case "-":
		return os.Stdout, nil
	default:
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("can't open file: %v", err)
		}

		return f, nil
	}
}
