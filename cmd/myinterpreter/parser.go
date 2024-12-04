package main

import (
	"fmt"
	"slices"
)

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

//	program        → statement* EOF ;
//
//	statement      → exprStmt
//					 | printStmt ;
//
//	exprStmt       → expression ";" ;
//	printStmt      → "print" expression ";" ;
//	expression     → equality ;
//	equality       → comparison ( ( "!=" | "==" ) comparison )* ;
//	comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
//	term           → factor ( ( "-" | "+" ) factor )* ;
//	factor         → unary ( ( "/" | "*" ) unary )* ;
//	unary          → ( "!" | "-" ) unary
//					 | primary ;
//	primary        → NUMBER | STRING | "true" | "false" | "nil"
//					 | "(" expression ")" ;

func (p *Parser) NextExpression() (Expression, error) {
	return p.parseExpression()
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseEquality()
}

func (p *Parser) parseSequence(parseFunc func() (Expression, error), matchers ...TokenType) (Expression, error) {
	e, err := parseFunc()
	if err != nil {
		return nil, err
	}

	for {
		token, ok := p.nextToken()
		if !ok || !slices.Contains(matchers, token.Type) {
			break
		}

		rightExpr, err := parseFunc()
		if err != nil {
			return nil, err
		}

		e = &BinaryExpr{
			Operator:  string(token.Type),
			LeftExpr:  e,
			RightExpr: rightExpr,
			Line:      token.Line,
		}
	}

	p.goBack()

	return e, nil
}

func (p *Parser) parseEquality() (Expression, error) {
	return p.parseSequence(p.parseComparison, BANG_EQUAL, EQUAL_EQUAL)
}

func (p *Parser) parseComparison() (Expression, error) {
	return p.parseSequence(p.parseTerm, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL)
}

func (p *Parser) parseTerm() (Expression, error) {
	return p.parseSequence(p.parseFactor, PLUS, MINUS)
}

func (p *Parser) parseFactor() (Expression, error) {
	return p.parseSequence(p.parseUnary, SLASH, STAR)
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
			Line:  token.Line,
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
	case token.Lexeme == "true":
		currExpr = &LiteralExpr{Literal: true, Line: token.Line}
	case token.Lexeme == "false":
		currExpr = &LiteralExpr{Literal: false, Line: token.Line}
	case token.Lexeme == "nil":
		currExpr = &LiteralExpr{Literal: nil, Line: token.Line}
	case token.Type == NUMBER || token.Type == STRING:
		currExpr = &LiteralExpr{Literal: token.Literal, Line: token.Line}
	case token.Type == LEFT_PAREN:
		e, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		n, exists := p.nextToken()
		if !exists || n.Type != RIGHT_PAREN {
			return nil, fmt.Errorf("[line %s] Unbalanced parentheses.", token.Line+1)
		}

		currExpr = &GroupingExpr{Expr: e, Line: token.Line}
	default:
		return nil, fmt.Errorf("[line %d] Error at '%s': Expect expression.", token.Line+1, token.Lexeme)
	}

	return currExpr, nil
}

func (p *Parser) nextToken() (*Token, bool) {
	p.pos++

	if p.pos >= len(p.tokens)-1 { // last token is EOF
		p.pos = len(p.tokens) - 1
		return nil, false
	}

	return p.tokens[p.pos], true
}

func (p *Parser) goBack() {
	if p.pos == -1 {
		return
	}

	p.pos--
}

var ErrNoMoreTokens = fmt.Errorf("no more tokens")
