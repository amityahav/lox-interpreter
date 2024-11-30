package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var errFound bool

	s := NewScanner(fileContents)

	for s.HasNext() {
		token, err := s.NextToken()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			errFound = true
		} else {
			fmt.Println(token)
		}
	}

	if errFound {
		os.Exit(65)
	}
}
