package main

import (
	"errors"
	"fmt"
	"os"
)

type Scope struct {
	Bindings map[string]interface{}
}

func (s *Scope) SetBinding(name string, value interface{}) {
	s.Bindings[name] = value
}

type State struct {
	Scopes []*Scope
}

func (s *State) GrowScopes() {
	s.Scopes = append(s.Scopes, &Scope{make(map[string]interface{})})
}

func (s *State) GetInnermostScope() *Scope {
	return s.Scopes[len(s.Scopes)-1]
}

func (s *State) CloseInnermostScope() {
	if isGlobalScope := len(s.Scopes) == 1; isGlobalScope {
		return
	}

	s.Scopes = s.Scopes[:len(s.Scopes)-1]
}

func (s *State) GetScopeFor(name string) (*Scope, bool) {
	// search for the binding through all existing scopes, starting from the innermost one
	for i := len(s.Scopes) - 1; i >= 0; i-- {
		_, ok := s.Scopes[i].Bindings[name]
		if ok {
			return s.Scopes[i], ok
		}
	}

	return nil, false
}

type Interpreter struct {
	state State
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		state: State{Scopes: []*Scope{{make(map[string]interface{})}}}, // first scope is the global scope
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
