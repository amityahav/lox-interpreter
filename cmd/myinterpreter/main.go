package main

import (
	"errors"
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

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var (
		errFound bool
		tokens   []*Token
	)

	if command == "tokenize" {
		s := NewScanner(fileContents)
		for s.HasNext() {
			token, err := s.NextToken()
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error()+"\n")
				errFound = true
			} else {
				fmt.Println(token)
			}
		}

		if errFound {
			os.Exit(65)
		}
	} else if command == "parse" {
		s := NewScanner(fileContents)
		for s.HasNext() {
			token, err := s.NextToken()
			if err != nil {
				errFound = true
			} else {
				tokens = append(tokens, token)
			}
		}

		p := NewParser(tokens)
		for expr, err := p.NextExpression(); !errors.Is(err, ErrNoMoreTokens); expr, err = p.NextExpression() {
			fmt.Println(expr.String())
		}
	}
}
