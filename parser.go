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
//	declaration    → funDecl
//					 | varDecl
//				     | statement ;
//	funDecl        → "fun" function ;
//	function       → IDENTIFIER "(" parameters? ")" block ;
// 	parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
//  varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
//	statement      → exprStmt
//					 | forStmt
//					 | ifStmt
//					 | printStmt
//					 | returnStmt
//					 | whileStmt
//					 | block ;
//
//	returnStmt     → "return" expression? ";" ;
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

func (p *Parser) NextDeclaration() (Statement, error) {
	return p.parseDeclaration()
}

func (p *Parser) parseDeclaration() (Statement, error) {
	token, ok := p.peek()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case FUN:
		return p.parseFunDeclaration()
	case VAR:
		return p.parseVarDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseFunDeclaration() (Statement, error) {
	_, err := p.match(FUN)
	if err != nil {
		return nil, err
	}

	return p.parseFunction()
}

func (p *Parser) parseFunction() (Statement, error) {
	token, err := p.match(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	funName := token.Lexeme

	token, err = p.match(LEFT_PAREN)
	if err != nil {
		return nil, err
	}

	var params []IdentifierExpr

	token, err = p.match(RIGHT_PAREN)
	if err != nil {
		params, err = p.parseParameters()
		if err != nil {
			return nil, err
		}

		_, err = p.match(RIGHT_PAREN)
		if err != nil {
			return nil, err
		}
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &FunDeclStmt{
		Name:   funName,
		Params: params,
		Body:   block,
	}, nil
}

func (p *Parser) parseParameters() ([]IdentifierExpr, error) {
	var params []IdentifierExpr

	token, err := p.match(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	params = append(params, IdentifierExpr{
		Name: token.Lexeme,
		Line: token.Line,
	})

	for {
		token, err := p.match(COMMA)
		if err != nil {
			break
		}

		token, err = p.match(IDENTIFIER)
		if err != nil {
			return nil, err
		}

		params = append(params, IdentifierExpr{
			Name: token.Lexeme,
			Line: token.Line,
		})
	}

	return params, nil
}

func (p *Parser) parseVarDeclaration() (Statement, error) {
	token, err := p.match(VAR)
	if err != nil {
		return nil, err
	}

	token, err = p.match(IDENTIFIER)
	if err != nil {
		return nil, err
	}

	varName := token.Lexeme
	var expr Expression = &NilExpr{}

	token, err = p.match(SEMICOLON)
	if err != nil {
		token, err = p.match(EQUAL)
		if err != nil {
			return nil, err
		}

		expr, err = p.parseExpression()
		if err != nil {
			if errors.Is(err, ErrNoMoreTokens) {
				return nil, fmt.Errorf("[line %d] Error: Expected expression.", token.Line+1, token.Lexeme)
			}

			return nil, err
		}

		token, err = p.match(SEMICOLON)
		if err != nil {
			return nil, err
		}
	}

	return &VarDeclStmt{
		Name: varName,
		Expr: expr,
	}, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	token, ok := p.peek()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case PRINT:
		return p.parsePrintStatement()
	case LEFT_BRACE:
		return p.parseBlockStatement()
	case IF:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForStatement()
	case RETURN:
		return p.parseReturnStatement()
	}

	return p.parseExprStatement()
}

func (p *Parser) parsePrintStatement() (Statement, error) {
	_, err := p.match(PRINT)
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.match(SEMICOLON)
	if err != nil {
		return nil, err
	}

	return &PrintStmt{Expr: expr}, nil
}

func (p *Parser) parseBlockStatement() (Statement, error) {
	_, err := p.match(LEFT_BRACE)
	if err != nil {
		return nil, err
	}

	var stmts []Statement

	for {
		_, err := p.match(RIGHT_BRACE)
		if err != nil {
			if errors.Is(err, UnexpectedEOF) {
				return nil, err
			}

			stmt, err := p.parseDeclaration()
			if err != nil {
				return nil, err
			}

			stmts = append(stmts, stmt)
		} else {
			break
		}
	}

	return &BlockStatement{
		Stmts: stmts,
	}, nil
}

func (p *Parser) parseIfStatement() (Statement, error) {
	_, err := p.match(IF)
	if err != nil {
		return nil, err
	}

	_, err = p.match(LEFT_PAREN)
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.match(RIGHT_PAREN)
	if err != nil {
		return nil, err
	}

	then, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var els Statement = &NilStmt{}

	_, err = p.match(ELSE)
	if err == nil {
		els, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{
		Condition: condition,
		Then:      then,
		Else:      els,
	}, nil
}

func (p *Parser) parseWhileStatement() (Statement, error) {
	_, err := p.match(WHILE)
	if err != nil {
		return nil, err
	}

	_, err = p.match(LEFT_PAREN)
	if err != nil {
		return nil, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.match(RIGHT_PAREN)
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parseForStatement() (Statement, error) {
	_, err := p.match(FOR)
	if err != nil {
		return nil, err
	}

	_, err = p.match(LEFT_PAREN)
	if err != nil {
		return nil, err
	}

	token, ok := p.peek()
	if !ok {
		return nil, fmt.Errorf("Error: Expected statment, got EOF.")
	}

	var initializer Statement = &NilStmt{}

	switch token.Type {
	case VAR:
		i, err := p.parseVarDeclaration()
		if err != nil {
			return nil, err
		}

		initializer = i
	case SEMICOLON:
		p.nextToken()
	default:
		i, err := p.parseExprStatement()
		if err != nil {
			return nil, err
		}

		initializer = i
	}

	var condition Expression = &NilExpr{}

	_, err = p.match(SEMICOLON)
	if err != nil {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.match(SEMICOLON)
		if err != nil {
			return nil, err
		}

		condition = expr
	}

	var increment Expression = &NilExpr{}

	_, err = p.match(RIGHT_PAREN)
	if err != nil {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.match(RIGHT_PAREN)
		if err != nil {
			return nil, err
		}

		increment = expr
	}

	body, err := p.parseStatement()
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
				},
			},
		},
	}, nil
}

func (p *Parser) parseReturnStatement() (Statement, error) {
	_, err := p.match(RETURN)
	if err != nil {
		return nil, err
	}

	var expr Expression = &NilExpr{}

	_, err = p.match(SEMICOLON)
	if err != nil {
		expr, err = p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.match(SEMICOLON)
		if err != nil {
			return nil, err
		}
	}

	return &ReturnStmt{Expr: expr}, nil
}

func (p *Parser) parseExprStatement() (Statement, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.match(SEMICOLON)
	if err != nil {
		return nil, err
	}

	return &ExprStmt{Expr: expr}, nil
}

func (p *Parser) NextExpression() (Expression, error) {
	return p.parseExpression()
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() (Expression, error) {
	token, err := p.match(IDENTIFIER)
	if err != nil {
		return p.parseLogicOr()
	}

	varName := token.Lexeme

	_, err = p.match(EQUAL)
	if err != nil {
		p.goBack()
		return p.parseLogicOr()
	}

	expr, err := p.parseAssignment()
	if err != nil {
		if errors.Is(err, ErrNoMoreTokens) {
			return nil, fmt.Errorf("[line %d] Error: Expected expression.", token.Line+1)
		}

		return nil, err
	}

	return &AssignmentExpr{
		Name: varName,
		Expr: expr,
		Line: token.Line,
	}, nil
}

func (p *Parser) parseSequenceBinary(parseFunc func() (Expression, error), matcher TokenType, matchers ...TokenType) (Expression, error) {
	e, err := parseFunc()
	if err != nil {
		return nil, err
	}

	for {
		token, err := p.match(matcher, matchers...)
		if err != nil {
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

	return e, nil
}

func (p *Parser) parseSequenceLogical(parseFunc func() (Expression, error), matcher TokenType, matchers ...TokenType) (Expression, error) {
	e, err := parseFunc()
	if err != nil {
		return nil, err
	}

	for {
		token, err := p.match(matcher, matchers...)
		if err != nil {
			break
		}

		rightExpr, err := parseFunc()
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

func (p *Parser) parseLogicOr() (Expression, error) {
	return p.parseSequenceLogical(p.parseLogicAnd, OR)
}

func (p *Parser) parseLogicAnd() (Expression, error) {
	return p.parseSequenceLogical(p.parseEquality, AND)
}

func (p *Parser) parseEquality() (Expression, error) {
	return p.parseSequenceBinary(p.parseComparison, BANG_EQUAL, EQUAL_EQUAL)
}

func (p *Parser) parseComparison() (Expression, error) {
	return p.parseSequenceBinary(p.parseTerm, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL)
}

func (p *Parser) parseTerm() (Expression, error) {
	return p.parseSequenceBinary(p.parseFactor, PLUS, MINUS)
}

func (p *Parser) parseFactor() (Expression, error) {
	return p.parseSequenceBinary(p.parseUnary, SLASH, STAR)
}

func (p *Parser) parseUnary() (Expression, error) {
	token, ok := p.peek()
	if !ok {
		return nil, ErrNoMoreTokens
	}

	switch token.Type {
	case BANG, MINUS:
		p.nextToken()

		u, err := p.parseUnary()
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

	return p.parseCall()
}

func (p *Parser) parseCall() (Expression, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		token, err := p.match(LEFT_PAREN)
		if err != nil {
			break
		}

		var args []Expression

		_, err = p.match(RIGHT_PAREN)
		if err != nil {
			args, err = p.parseArguments()
			if err != nil {
				return nil, err
			}

			_, err = p.match(RIGHT_PAREN)
			if err != nil {
				return nil, err
			}
		}

		expr = &CallExpr{
			Callee: expr,
			Args:   args,
			Line:   token.Line,
		}
	}

	return expr, nil
}

func (p *Parser) parseArguments() ([]Expression, error) {
	var args []Expression

	e, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	args = append(args, e)

	for {
		_, err = p.match(COMMA)
		if err != nil {
			break
		}

		exp, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		args = append(args, exp)
	}

	return args, nil
}

func (p *Parser) parsePrimary() (Expression, error) {
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
		currExpr = &LiteralExpr{Literal: &NilExpr{}, Line: token.Line}
	case NUMBER, STRING:
		currExpr = &LiteralExpr{Literal: token.Literal, Line: token.Line}
	case IDENTIFIER:
		currExpr = &IdentifierExpr{Name: token.Lexeme, Line: token.Line}
	case LEFT_PAREN:
		e, err := p.parseExpression()
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

func (p *Parser) match(tokenType TokenType, tokenTypes ...TokenType) (*Token, error) {
	token, ok := p.nextToken()
	if !ok {
		p.goBack()
		return nil, fmt.Errorf("Error: Expected '%s', got %w.", string(tokenType), UnexpectedEOF)
	}

	if !token.Type.Is(tokenType) && !slices.Contains(tokenTypes, token.Type) {
		p.goBack()
		return nil, fmt.Errorf("[line %d] Error at '%s': Expected '%s'.", token.Line+1, token.Lexeme)
	}

	return token, nil
}

func (p *Parser) peek() (*Token, bool) {
	if p.pos+1 >= len(p.tokens) {
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

var (
	ErrNoMoreTokens = fmt.Errorf("no more tokens")
	UnexpectedEOF   = fmt.Errorf("unexpected EOF")
)
