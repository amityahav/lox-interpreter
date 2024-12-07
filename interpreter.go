package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Caller interface {
	Call(args ...interface{}) (interface{}, error)
}

type NativeClock struct{}

func (nc *NativeClock) Call(_ ...interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

type Scope struct {
	Bindings map[string]interface{}
}

func (s *Scope) SetBinding(name string, value interface{}) {
	s.Bindings[name] = value
}

type Environment struct {
	Scopes []*Scope
}

func (e *Environment) GrowScopes() {
	e.Scopes = append(e.Scopes, &Scope{make(map[string]interface{})})
}

func (e *Environment) GetInnermostScope() *Scope {
	return e.Scopes[len(e.Scopes)-1]
}

func (e *Environment) CloseInnermostScope() {
	if isGlobalScope := len(e.Scopes) == 1; isGlobalScope {
		return
	}

	e.Scopes = e.Scopes[:len(e.Scopes)-1]
}

func (e *Environment) GetScopeFor(name string) (*Scope, bool) {
	// search for the binding through all existing scopes, starting from the innermost one
	for i := len(e.Scopes) - 1; i >= 0; i-- {
		_, ok := e.Scopes[i].Bindings[name]
		if ok {
			return e.Scopes[i], ok
		}
	}

	return nil, false
}

type Interpreter struct {
	env Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: Environment{Scopes: []*Scope{{map[string]interface{}{
			"clock": &NativeClock{},
		}}}}, // first scope is the global scope
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
