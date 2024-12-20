package main

import (
	"errors"
	"fmt"

	"os"
)

func main() {
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
			if err != nil {
				errFound = true
				fmt.Fprintln(os.Stderr, err.Error()+"\n")
			} else {
				fmt.Println(expr)
			}
		}

		if errFound {
			os.Exit(65)
		}
	} else if command == "evaluate" {
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
			if err != nil {
				errFound = true
				fmt.Fprintln(os.Stderr, err.Error()+"\n")
			} else {
				v, err := expr.Eval(nil)
				if err != nil {
					errFound = true
					fmt.Fprintln(os.Stderr, err.Error()+"\n")
					os.Exit(70)
				}

				fmt.Println(strHelper(v))
			}
		}

		if errFound {
			os.Exit(65)
		}
	} else if command == "run" {
		_ = NewInterpreter().Interpret(fileContents)
	}
}

func strHelper(v interface{}) string {
	if v == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", v)
}
