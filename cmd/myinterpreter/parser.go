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

func (ue *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", ue.Unary, ue.Expr.String())
}

type BinaryExpr struct {
	Operator  string
	LeftExpr  Expression
	RightExpr Expression
}

func (be *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Operator, be.LeftExpr, be.RightExpr)
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
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    -1,
	}
}

var ErrNoMoreTokens = fmt.Errorf("no more tokens")

//	expression     → equality ;
//	equality       → comparison ( ( "!=" | "==" ) comparison )* ;
//	comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
//	term           → factor ( ( "-" | "+" ) factor )* ;
//	factor         → unary ( ( "/" | "*" ) unary )* ;
//	unary          → ( "!" | "-" ) unary
//					 | primary ;
//	primary        → NUMBER | STRING | "true" | "false" | "nil"
//					 | "(" expression ")" ;

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseEquality()
}

func (p *Parser) parseEquality() (Expression, error) {
	return p.parseComparison()
}

func (p *Parser) parseComparison() (Expression, error) {
	return p.parseTerm()
}

func (p *Parser) parseTerm() (Expression, error) {
	return p.parseFactor()
}

func (p *Parser) parseFactor() (Expression, error) {
	e, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		token, ok := p.nextToken()
		if !ok {
			break
		}

		if token.Type != SLASH && token.Type != STAR {
			p.goBack()
			break
		}

		rightExpr, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		e = &BinaryExpr{
			Operator:  string(token.Type),
			LeftExpr:  e,
			RightExpr: rightExpr,
		}
	}

	return e, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case BANG, MINUS:
		u, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return &UnaryExpr{
			Unary: string(token.Type),
			Expr:  u,
		}, nil
	}

	p.goBack()
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expression, error) {
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
		e, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		ge.Expr = e

		n, exists := p.peek()
		if !exists || n.Type != RIGHT_PAREN {
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

func (p *Parser) NextExpression() (Expression, error) {
	return p.parseExpression()
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
	if p.pos+1 >= len(p.tokens)-1 {
		return nil, false
	}

	return p.tokens[p.pos+1], true
}

func (p *Parser) goBack() {
	if p.pos == -1 {
		return
	}

	p.pos--
}

func isOperator(token *Token) bool {
	switch token.Type {
	case PLUS, MINUS, STAR, EQUAL, EQUAL_EQUAL, BANG, BANG_EQUAL, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL, SLASH:
		return true
	}

	return false
}
