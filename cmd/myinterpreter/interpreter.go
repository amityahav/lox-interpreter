package main

import (
	"errors"
	"fmt"
	"os"
)

type State struct {
	Globals map[string]interface{}
}

func (s *State) AddGlobal(key string, value interface{}) {
	s.Globals[key] = value
}

func (s *State) GetGlobal(key string) (interface{}, bool) {
	val, ok := s.Globals[key]
	return val, ok
}

type Interpreter struct {
	state State
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		state: State{Globals: make(map[string]interface{})},
	}
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

	for stmt, err := parser.NextDeclaration(&i.state); !errors.Is(err, ErrNoMoreTokens); stmt, err = parser.NextDeclaration(&i.state) {
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
