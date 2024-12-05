package main

import (
	"errors"
	"fmt"
	"os"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(content []byte) error {
	scanner := NewScanner(content)

	var tokens []*Token
	for scanner.HasNext() {
		token, err := scanner.NextToken()
		if err != nil {
			return err
		}

		tokens = append(tokens, token)
	}

	parser := NewParser(tokens)

	for stmt, err := parser.NextStatement(); !errors.Is(err, ErrNoMoreTokens); stmt, err = parser.NextStatement() {
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(65)
		}

		_, err = stmt.Execute()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(70)
		}
	}

	return nil
}
