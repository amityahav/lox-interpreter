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

func (le *LiteralExpr) String() string {
	return le.Literal
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
	LeftParen  string
	RightParen string
	Expr       Expression
}

func (ge *GroupingExpr) String() string {
	return fmt.Sprintf("(group %s)", ge.Expr.String())
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
	case token.Type == NUMBER || token.Type == STRING:
		currExpr = &LiteralExpr{Literal: token.Literal}
	case token.Type == LEFT_PAREN:
		var ge GroupingExpr
		e, err := p.NextExpression()
		if err != nil {
			return nil, err
		}

		ge.Expr = e

		n, exists := p.peek()
		if !exists {
			return nil, ErrNoMoreTokens
		}

		if n.Type != RIGHT_PAREN {
			return nil, fmt.Errorf("unbalanced parenthesis")
		}

		currExpr = &ge
		p.nextToken()
	case token.Type == RIGHT_PAREN:
		// TODO: we get here if there's an empty group or an unbalanced parenthesis
		return nil, fmt.Errorf("something")
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

func (p *Parser) peek() (*Token, bool) {
	if p.pos >= len(p.tokens) {
		return nil, false
	}

	return p.tokens[p.pos+1], true
}
