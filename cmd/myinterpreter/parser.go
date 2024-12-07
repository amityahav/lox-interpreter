package main

import (
	"errors"
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

// Grammar Rules:
//
//	program        → declaration* EOF ;
//
//	declaration    → varDecl
//				     | statement ;
//  varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
//	statement      → exprStmt
//				     | forStmt
//					 | ifStmt
//				  	 | printStmt
//					 | whileStmt
//					 | block ;
//
//	forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//					  expression? ";"
//					  expression? ")" statement ;
//	whileStmt      → "while" "(" expression ")" statement ;
//	ifStmt         → "if" "(" expression ")" statement
//					  ( "else" statement )? ;
//
//	block          → "{" declaration* "}" ;
//	exprStmt       → expression ";" ;
//	printStmt      → "print" expression ";" ;
//	expression     → assignment ;
//	assignment     → IDENTIFIER "=" assignment
//					 | logic_or ;
//	logic_or       → logic_and ( "or" logic_and )* ;
//	logic_and      → equality ( "and" equality )* ;
//	equality       → comparison ( ( "!=" | "==" ) comparison )* ;
//	comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
//	term           → factor ( ( "-" | "+" ) factor )* ;
//	factor         → unary ( ( "/" | "*" ) unary )* ;
//	unary          → ( "!" | "-" ) unary | call ;
//	call           → primary ( "(" arguments? ")" )* ;
//  arguments      → expression ( "," expression )* ;
//	primary        → "true" | "false" | "nil"
//					 | NUMBER | STRING
//				     | "(" expression ")"
//				     | IDENTIFIER ;

func (p *Parser) NextDeclaration(state *Environment) (Statement, error) {
	return p.parseDeclaration(state)
}

func (p *Parser) parseDeclaration(state *Environment) (Statement, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	if token.Type.Is(VAR) {
		return p.parseVarDeclaration(state)
	}

	p.goBack()

	return p.parseStatement(state)
}

func (p *Parser) parseVarDeclaration(state *Environment) (Statement, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected IDENTIFIER, got EOF.")
	}

	if !token.Type.Is(IDENTIFIER) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected IDENTIFIER.", token.Line+1, token.Lexeme)
	}

	varName := token.Lexeme

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ';', got EOF.")
	}

	if token.Type.Is(SEMICOLON) {
		// uninitialized variable
		return &VarDeclStmt{
			Name:  varName,
			Expr:  &NoopExpr{},
			state: state,
		}, nil
	}

	if !token.Type.Is(EQUAL) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected ';' or '='.", token.Line+1, token.Lexeme)
	}

	expr, err := p.parseExpression(state)
	if err != nil {
		if errors.Is(err, ErrNoMoreTokens) {
			return nil, fmt.Errorf("[line %d] Error: Expected expression.", token.Line+1, token.Lexeme)
		}

		return nil, err
	}

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ';', got EOF.")
	}

	if !token.Type.Is(SEMICOLON) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected ';'.", token.Line+1, token.Lexeme)
	}

	return &VarDeclStmt{
		Name:  varName,
		Expr:  expr,
		state: state,
	}, nil
}

func (p *Parser) parseStatement(state *Environment) (Statement, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case PRINT:
		return p.parsePrintStatement(state)
	case LEFT_BRACE:
		return p.parseBlockStatement(state)
	case IF:
		return p.parseIfStatement(state)
	case WHILE:
		return p.parseWhileStatement(state)
	case FOR:
		return p.parseForStatement(state)
	}

	p.goBack()

	return p.parseExprStatement(state)
}

func (p *Parser) parsePrintStatement(state *Environment) (Statement, error) {
	expr, err := p.parseExpression(state)
	if err != nil {
		return nil, err
	}

	token, ok := p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ';', got EOF.")
	}

	if !token.Type.Is(SEMICOLON) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected ';'.", token.Line+1, token.Lexeme)
	}

	return &PrintStmt{Expr: expr}, nil
}

func (p *Parser) parseBlockStatement(state *Environment) (Statement, error) {
	var stmts []Statement

	for {
		token, ok := p.nextToken()
		if !ok {
			return nil, fmt.Errorf("Error: Expected '}', got EOF.")
		}

		if token.Type.Is(RIGHT_BRACE) {
			return &BlockStatement{
				Stmts: stmts,
				state: state}, nil
		}

		p.goBack()

		stmt, err := p.parseDeclaration(state)
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}
}

func (p *Parser) parseIfStatement(state *Environment) (Statement, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected '(', got EOF.")
	}

	if !token.Type.Is(LEFT_PAREN) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected '('.", token.Line+1, token.Lexeme)
	}

	condition, err := p.parseExpression(state)
	if err != nil {
		return nil, err
	}

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ')', got EOF.")
	}

	if !token.Type.Is(RIGHT_PAREN) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected ')'.", token.Line+1, token.Lexeme)
	}

	then, err := p.parseStatement(state)
	if err != nil {
		return nil, err
	}

	token, ok = p.nextToken()
	if !ok || !token.Type.Is(ELSE) {
		if ok {
			p.goBack()
		}

		return &IfStmt{
			Condition: condition,
			Then:      then,
			Else:      &NoopStmt{},
		}, nil
	}

	els, err := p.parseStatement(state)
	if err != nil {
		return nil, err
	}

	return &IfStmt{
		Condition: condition,
		Then:      then,
		Else:      els,
	}, nil
}

func (p *Parser) parseWhileStatement(state *Environment) (Statement, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected '(', got EOF.")
	}

	if !token.Type.Is(LEFT_PAREN) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected '('.", token.Line+1, token.Lexeme)
	}

	condition, err := p.parseExpression(state)
	if err != nil {
		return nil, err
	}

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ')', got EOF.")
	}

	if !token.Type.Is(RIGHT_PAREN) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected ')'.", token.Line+1, token.Lexeme)
	}

	body, err := p.parseStatement(state)
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parseForStatement(state *Environment) (Statement, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected '(', got EOF.")
	}

	if !token.Type.Is(LEFT_PAREN) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected '('.", token.Line+1, token.Lexeme)
	}

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected statment, got EOF.")
	}

	var initializer Statement = &NoopStmt{}
	switch token.Type {
	case VAR:
		i, err := p.parseVarDeclaration(state)
		if err != nil {
			return nil, err
		}

		initializer = i
	case SEMICOLON:
	default:
		p.goBack()

		i, err := p.parseExprStatement(state)
		if err != nil {
			return nil, err
		}

		initializer = i
	}

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ';' or expression, got EOF.")
	}

	var condition Expression = &NoopExpr{}
	if !token.Type.Is(SEMICOLON) {
		p.goBack()

		expr, err := p.parseExpression(state)
		if err != nil {
			return nil, err
		}

		token, ok = p.nextToken()
		if !ok {
			return nil, fmt.Errorf("Error: Expected ';' or expression, got EOF.")
		}

		if !token.Type.Is(SEMICOLON) {
			return nil, fmt.Errorf("[line %d] Error at '%s': Expected ';'.", token.Line+1, token.Lexeme)
		}

		condition = expr
	}

	token, ok = p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ')' or expression, got EOF.")
	}

	var increment Expression = &NoopExpr{}
	if !token.Type.Is(RIGHT_PAREN) {
		p.goBack()

		expr, err := p.parseExpression(state)
		if err != nil {
			return nil, err
		}

		token, ok = p.nextToken()
		if !ok {
			return nil, fmt.Errorf("Error: Expected ')' or expression, got EOF.")
		}

		if !token.Type.Is(RIGHT_PAREN) {
			return nil, fmt.Errorf("[line %d] Error at '%s': Expected ')'.", token.Line+1, token.Lexeme)
		}

		increment = expr
	}

	body, err := p.parseStatement(state)
	if err != nil {
		return nil, err
	}

	// desugaring for-loop to while-loop
	return &BlockStatement{
		Stmts: []Statement{
			initializer,
			&WhileStmt{
				Condition: condition,
				Body: &BlockStatement{
					Stmts: []Statement{
						body,
						&ExprStmt{Expr: increment},
					},
					state: state,
				},
			},
		},
		state: state,
	}, nil
}

func (p *Parser) parseExprStatement(state *Environment) (Statement, error) {
	expr, err := p.parseExpression(state)
	if err != nil {
		return nil, err
	}

	token, ok := p.nextToken()
	if !ok {
		return nil, fmt.Errorf("Error: Expected ';', got EOF.")
	}

	if !token.Type.Is(SEMICOLON) {
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected ';'.", token.Line+1, token.Lexeme)
	}

	return &ExprStmt{Expr: expr}, nil
}

func (p *Parser) NextExpression(state *Environment) (Expression, error) {
	return p.parseExpression(state)
}

func (p *Parser) parseExpression(state *Environment) (Expression, error) {
	return p.parseAssignment(state)
}

func (p *Parser) parseAssignment(state *Environment) (Expression, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	if token.Type.Is(IDENTIFIER) {
		varName := token.Lexeme

		token, ok = p.nextToken()
		if !ok {
			return nil, fmt.Errorf("Error: Expected '=' or ';', got EOF.")
		}

		if !token.Type.Is(EQUAL) {
			p.goBack()
			p.goBack()
			return p.parseEquality(state)
		}

		expr, err := p.parseAssignment(state)
		if err != nil {
			if errors.Is(err, ErrNoMoreTokens) {
				return nil, fmt.Errorf("[line %d] Error: Expected expression.", token.Line+1)
			}

			return nil, err
		}

		return &AssignmentExpr{
			Name:  varName,
			Expr:  expr,
			Line:  token.Line,
			state: state,
		}, nil
	}

	p.goBack()

	return p.parseLogicOr(state)
}

func (p *Parser) parseSequenceBinary(state *Environment, parseFunc func(state *Environment) (Expression, error), matchers ...TokenType) (Expression, error) {
	e, err := parseFunc(state)
	if err != nil {
		return nil, err
	}

	for {
		token, ok := p.nextToken()
		if !ok || !slices.Contains(matchers, token.Type) {
			if ok {
				p.goBack()
			}

			break
		}

		rightExpr, err := parseFunc(state)
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

	return e, nil
}

func (p *Parser) parseSequenceLogical(state *Environment, parseFunc func(state *Environment) (Expression, error), matchers ...TokenType) (Expression, error) {
	e, err := parseFunc(state)
	if err != nil {
		return nil, err
	}

	for {
		token, ok := p.nextToken()
		if !ok || !slices.Contains(matchers, token.Type) {
			if ok {
				p.goBack()
			}

			break
		}

		rightExpr, err := parseFunc(state)
		if err != nil {
			return nil, err
		}

		e = &LogicalExpr{
			Operator:  string(token.Type),
			LeftExpr:  e,
			RightExpr: rightExpr,
		}
	}

	return e, nil
}

func (p *Parser) parseLogicOr(state *Environment) (Expression, error) {
	return p.parseSequenceLogical(state, p.parseLogicAnd, OR)
}

func (p *Parser) parseLogicAnd(state *Environment) (Expression, error) {
	return p.parseSequenceLogical(state, p.parseEquality, AND)
}

func (p *Parser) parseEquality(state *Environment) (Expression, error) {
	return p.parseSequenceBinary(state, p.parseComparison, BANG_EQUAL, EQUAL_EQUAL)
}

func (p *Parser) parseComparison(state *Environment) (Expression, error) {
	return p.parseSequenceBinary(state, p.parseTerm, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL)
}

func (p *Parser) parseTerm(state *Environment) (Expression, error) {
	return p.parseSequenceBinary(state, p.parseFactor, PLUS, MINUS)
}

func (p *Parser) parseFactor(state *Environment) (Expression, error) {
	return p.parseSequenceBinary(state, p.parseUnary, SLASH, STAR)
}

func (p *Parser) parseUnary(state *Environment) (Expression, error) {
	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case BANG, MINUS:
		u, err := p.parseUnary(state)
		if err != nil {
			if errors.Is(err, ErrNoMoreTokens) {
				return nil, fmt.Errorf("[line %d] Error at '%s': Expect expression.", token.Line+1, token.Lexeme)
			}

			return nil, err
		}

		return &UnaryExpr{
			Unary: string(token.Type),
			Expr:  u,
			Line:  token.Line,
		}, nil
	}

	p.goBack()

	return p.parseCall(state)
}

// unary          → ( "!" | "-" ) unary | call ;
// call           → primary ( "(" arguments? ")" )* ;
func (p *Parser) parseCall(state *Environment) (Expression, error) {
	expr, err := p.parsePrimary(state)
	if err != nil {
		return nil, err
	}

	for {
		token, ok := p.nextToken()
		if !ok || !token.Type.Is(LEFT_PAREN) {
			if ok {
				p.goBack()
			}

			break
		}

		token, ok = p.nextToken()
		if !ok {
			return nil, fmt.Errorf("Error: Expected ')' or arguments, got EOF.")
		}

		if token.Type.Is(RIGHT_PAREN) {
			expr = &CallExpr{
				Callee: expr,
				Args:   nil,
			}

			continue
		}

		args, err := p.parseArguments(state)
		if err != nil {
			return nil, err
		}

		token, ok = p.nextToken()
		if !ok {
			return nil, fmt.Errorf("Error: Expected ')' or arguments, got EOF.")
		}

		if !token.Type.Is(RIGHT_PAREN) {
			return nil, fmt.Errorf("[line %d] Error at '%s': Expect ')'.", token.Line+1, token.Lexeme)
		}

		expr = &CallExpr{
			Callee: expr,
			Args:   args,
		}
	}

	return expr, nil
}

// arguments      → expression ( "," expression )* ;
func (p *Parser) parseArguments(state *Environment) ([]Expression, error) {
	var args []Expression

	e, err := p.parseExpression(state)
	if err != nil {
		return nil, err
	}

	args = append(args, e)

	for {
		token, ok := p.nextToken()
		if !ok || !token.Type.Is(COMMA) {
			if ok {
				p.goBack()
			}

			break
		}

		exp, err := p.parseExpression(state)
		if err != nil {
			return nil, err
		}

		args = append(args, exp)
	}

	return args, nil
}

func (p *Parser) parsePrimary(state *Environment) (Expression, error) {
	var currExpr Expression

	token, ok := p.nextToken()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case TRUE:
		currExpr = &LiteralExpr{Literal: true, Line: token.Line}
	case FALSE:
		currExpr = &LiteralExpr{Literal: false, Line: token.Line}
	case NIL:
		currExpr = &LiteralExpr{Literal: nil, Line: token.Line}
	case NUMBER, STRING:
		currExpr = &LiteralExpr{Literal: token.Literal, Line: token.Line}
	case IDENTIFIER:
		currExpr = &IdentifierExpr{Name: token.Lexeme, Line: token.Line, state: state}
	case LEFT_PAREN:
		e, err := p.parseExpression(state)
		if err != nil {
			if errors.Is(err, ErrNoMoreTokens) {
				return nil, fmt.Errorf("[line %d] Unbalanced parentheses.", token.Line+1)
			}

			return nil, err
		}

		n, exists := p.nextToken()
		if !exists || !n.Type.Is(RIGHT_PAREN) {
			return nil, fmt.Errorf("[line %d] Unbalanced parentheses.", token.Line+1)
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
