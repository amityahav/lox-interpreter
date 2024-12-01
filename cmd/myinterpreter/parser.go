package main

import (
	"fmt"
	"slices"
)

type Expression interface {
	Eval() (interface{}, error)
	String() string
}

type LiteralExpr struct {
	Literal interface{}
}

func (le *LiteralExpr) Eval() (interface{}, error) {
	return le.Literal, nil
}

func (le *LiteralExpr) String() string {
	if le.Literal == nil {
		return "nil"
	}

	if v, ok := le.Literal.(float64); ok {
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.1f", v)
		}
	}

	return fmt.Sprintf("%v", le.Literal)
}

type UnaryExpr struct {
	Unary string
	Expr  Expression
}

func (ue *UnaryExpr) Eval() (interface{}, error) {
	val, err := ue.Expr.Eval()
	if err != nil {
		return nil, err
	}

	switch TokenType(ue.Unary) {
	case MINUS:
		v, ok := val.(float64)
		if !ok {
			panic("not a num")
		}

		return -v, nil
	case BANG:
		return !isTrue(val), nil
	}

	// unreachable
	return nil, nil
}

func (ue *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", ue.Unary, ue.Expr.String())
}

type BinaryExpr struct {
	Operator  string
	LeftExpr  Expression
	RightExpr Expression
}

func (be *BinaryExpr) Eval() (interface{}, error) {
	leftVal, err := be.LeftExpr.Eval()
	if err != nil {
		return nil, err
	}

	rightVal, err := be.RightExpr.Eval()
	if err != nil {
		return nil, err
	}

	switch TokenType(be.Operator) {
	case SLASH:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv / rv, nil
	case STAR:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv * rv, nil
	case PLUS:
		lv, ok := leftVal.(float64)
		if ok {
			rv, ok2 := rightVal.(float64)
			if ok2 {
				return lv + rv, nil
			}
		}

		lvs, ok := leftVal.(string)
		if ok {
			rvs, ok2 := rightVal.(string)
			if ok2 {
				return lvs + rvs, nil
			}
		}

		panic("err")
	case MINUS:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv - rv, nil
	case LESS:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv < rv, nil
	case LESS_EQUAL:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv <= rv, nil
	case GREATER:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv > rv, nil
	case GREATER_EQUAL:
		lv, ok := leftVal.(float64)
		rv, ok2 := rightVal.(float64)
		if !ok || !ok2 {
			panic("not a num")
		}

		return lv >= rv, nil
	}

	// unreachable
	return nil, nil
}

func (be *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", be.Operator, be.LeftExpr, be.RightExpr)
}

type GroupingExpr struct {
	Expr Expression
}

func (ge *GroupingExpr) Eval() (interface{}, error) {
	return ge.Expr.Eval()
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
	case token.Lexeme == "true":
		currExpr = &LiteralExpr{Literal: true}
	case token.Lexeme == "false":
		currExpr = &LiteralExpr{Literal: false}
	case token.Lexeme == "nil":
		currExpr = &LiteralExpr{Literal: nil}
	case token.Type == NUMBER || token.Type == STRING:
		currExpr = &LiteralExpr{Literal: token.Literal}
	case token.Type == LEFT_PAREN:
		e, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		n, exists := p.nextToken()
		if !exists || n.Type != RIGHT_PAREN {
			return nil, fmt.Errorf("[line %s] Unbalanced parentheses.", token.Line+1)
		}

		currExpr = &GroupingExpr{Expr: e}
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

func isTrue(val interface{}) bool {
	if val == nil {
		return false
	}

	if v, ok := val.(bool); ok {
		return v
	}

	return true
}

var ErrNoMoreTokens = fmt.Errorf("no more tokens")
