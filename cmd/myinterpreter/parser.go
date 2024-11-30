package main

import (
	"fmt"
)

// expression     → literal
//				  | unary
//				  | binary
//				  | grouping ;
//
// literal        → NUMBER | STRING | "true" | "false" | "nil" ;
// grouping       → "(" expression ")" ;
// unary          → ( "-" | "!" ) expression ;
// binary         → expression operator expression ;
// operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
//					| "+"  | "-"  | "*" | "/" ;

type Expression interface {
	String() string
}

type LiteralExpr struct {
	Literal string
}

func (e *LiteralExpr) String() string {
	return e.Literal
}

type UnaryExpr struct {
	Unary string
	Expr  Expression
}

type BinaryExpr struct {
	Operator  string
	LeftExpr  Expression
	RightExpr Expression
}

type GroupingExpr struct {
	LeftBrace  string
	RightBrace string
	Expr       Expression
}

type Parser struct {
	tokens []*Token
	pos    int
	done   bool
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    -1,
	}
}

var ErrNoMoreTokens = fmt.Errorf("no more tokens")

func (p *Parser) NextExpression() (Expression, error) {
	var currExpr Expression

	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch {
	case token.Lexeme == "true" || token.Lexeme == "false" || token.Lexeme == "nil":
		currExpr = &LiteralExpr{Literal: token.Lexeme}
	}

	return currExpr, nil
}

func (p *Parser) HasNext() bool {
	return !p.done
}

func (p *Parser) nextToken() (*Token, bool) {
	p.pos++

	if p.pos >= len(p.tokens)-1 { // last token is EOF
		p.pos = len(p.tokens) - 1
		return nil, false
	}

	return p.tokens[p.pos], true
}
