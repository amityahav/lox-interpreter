package main

import (
	"fmt"
	"slices"
)

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
	Expr Expression
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
		if !ok {
			break
		}

		if slices.Contains(matchers, token.Type) {
			rightExpr, err := parseFunc()
			if err != nil {
				return nil, err
			}

			e = &BinaryExpr{
				Operator:  string(token.Type),
				LeftExpr:  e,
				RightExpr: rightExpr,
			}

			continue
		}

		break
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
		e, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		n, exists := p.nextToken()
		if !exists || n.Type != RIGHT_PAREN {
			return nil, fmt.Errorf("unbalanced parenthesis")
		}

		currExpr = &GroupingExpr{Expr: e}
	case token.Type == RIGHT_PAREN:
		// TODO: we get here if there's an empty group or an unbalanced parenthesis
		return nil, fmt.Errorf("something")
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
