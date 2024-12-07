package main

import (
	"errors"
	"fmt"
	"os"
)

type Environment struct {
	Bindings map[string]interface{}
	parent   *Environment
}

func (e *Environment) SetBinding(name string, value interface{}) {
	e.Bindings[name] = value
}

func (e *Environment) Lookup(name string) (*Environment, bool) {
	for curr := e; curr != nil; curr = curr.parent {
		if _, ok := curr.Bindings[name]; ok {
			return curr, true
		}
	}

	return nil, false
}

func ExpandEnv(parentEnv *Environment) *Environment {
	return &Environment{
		Bindings: make(map[string]interface{}),
		parent:   parentEnv,
	}
}

type Interpreter struct {
	env Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: Environment{
			Bindings: map[string]interface{}{
				"clock": &NativeClock{},
			},
		},
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

	for stmt, err := parser.NextDeclaration(); !errors.Is(err, ErrNoMoreTokens); stmt, err = parser.NextDeclaration() {
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(65)
		}

		_, err = stmt.Execute(&i.env)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(70)
		}
	}

	return nil
}
