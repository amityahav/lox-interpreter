package main

import (
	"errors"
	"fmt"
	"strconv"
)

type TokenType string

func (t TokenType) Is(t2 TokenType) bool { return t == t2 }

const (
	LEFT_PAREN    TokenType = "("
	RIGHT_PAREN   TokenType = ")"
	LEFT_BRACE    TokenType = "{"
	RIGHT_BRACE   TokenType = "}"
	COMMA         TokenType = ","
	DOT           TokenType = "."
	SEMICOLON     TokenType = ";"
	PLUS          TokenType = "+"
	MINUS         TokenType = "-"
	STAR          TokenType = "*"
	EQUAL         TokenType = "="
	EQUAL_EQUAL   TokenType = "=="
	BANG          TokenType = "!"
	BANG_EQUAL    TokenType = "!="
	LESS          TokenType = "<"
	LESS_EQUAL    TokenType = "<="
	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="
	SLASH         TokenType = "/"
	NEWLINE       TokenType = "\n"
	SPACE         TokenType = " "
	TAB           TokenType = "\t"
	STRING        TokenType = "<placeholder_str>"
	NUMBER        TokenType = "<placeholder_num>"
	QUOTE         TokenType = "\""
	IDENTIFIER    TokenType = "<placeholder_identifier>"
	AND           TokenType = "and"
	CLASS         TokenType = "class"
	ELSE          TokenType = "else"
	FALSE         TokenType = "false"
	FOR           TokenType = "for"
	FUN           TokenType = "fun"
	IF            TokenType = "if"
	NIL           TokenType = "nil"
	OR            TokenType = "or"
	PRINT         TokenType = "print"
	RETURN        TokenType = "return"
	SUPER         TokenType = "super"
	THIS          TokenType = "this"
	TRUE          TokenType = "true"
	VAR           TokenType = "var"
	WHILE         TokenType = "while"
	EOF           TokenType = ""
)

var reservedWords = map[TokenType]struct{}{
	AND:    {},
	CLASS:  {},
	ELSE:   {},
	FALSE:  {},
	FOR:    {},
	FUN:    {},
	IF:     {},
	NIL:    {},
	OR:     {},
	PRINT:  {},
	RETURN: {},
	SUPER:  {},
	THIS:   {},
	TRUE:   {},
	VAR:    {},
	WHILE:  {},
}

func (t TokenType) Type() string {
	switch t {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case DOT:
		return "DOT"
	case SEMICOLON:
		return "SEMICOLON"
	case STAR:
		return "STAR"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case SLASH:
		return "SLASH"
	case NEWLINE:
		return "NEWLINE"
	case SPACE:
		return "SPACE"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case TAB:
		return "TAB"
	case IDENTIFIER:
		return "IDENTIFIER"
	case EOF:
		return "EOF"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FOR:
		return "FOR"
	case FUN:
		return "FUN"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	default:
		return ""
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type.Type(), t.Lexeme, stringifyNumOrNil(t))
}

type Scanner struct {
	content []byte
	pos     int
	lineNum int
	done    bool
}

func NewScanner(content []byte) *Scanner {
	s := Scanner{
		content: content,
		pos:     -1,
	}

	return &s
}

func (s *Scanner) NextToken() (*Token, error) {
	var currToken Token

	for {
		currChar, ok := s.nextChar()
		switch {
		case !ok:
			currToken = Token{
				Type:    EOF,
				Lexeme:  string(EOF),
				Literal: nil,
				Line:    s.lineNum,
			}

			s.done = true
		case TokenType(currChar).Is(LEFT_PAREN) ||
			TokenType(currChar).Is(RIGHT_PAREN) ||
			TokenType(currChar).Is(LEFT_BRACE) ||
			TokenType(currChar).Is(RIGHT_BRACE) ||
			TokenType(currChar).Is(COMMA) ||
			TokenType(currChar).Is(DOT) ||
			TokenType(currChar).Is(SEMICOLON) ||
			TokenType(currChar).Is(PLUS) ||
			TokenType(currChar).Is(MINUS) ||
			TokenType(currChar).Is(STAR):
			currToken = Token{
				Type:    TokenType(currChar),
				Lexeme:  string(currChar),
				Literal: nil,
				Line:    s.lineNum,
			}
		case TokenType(currChar).Is(NEWLINE):
			s.lineNum++
			continue
		case TokenType(currChar).Is(EQUAL):
			if nextChar, exist := s.peek(); exist && TokenType(nextChar).Is(EQUAL) {
				currToken = Token{
					Type:    EQUAL_EQUAL,
					Lexeme:  string(EQUAL_EQUAL),
					Literal: nil,
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    EQUAL,
				Lexeme:  string(EQUAL),
				Literal: nil,
				Line:    s.lineNum,
			}
		case TokenType(currChar).Is(BANG):
			if nextChar, exist := s.peek(); exist && TokenType(nextChar).Is(EQUAL) {
				currToken = Token{
					Type:    BANG_EQUAL,
					Lexeme:  string(BANG_EQUAL),
					Literal: nil,
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    BANG,
				Lexeme:  string(BANG),
				Literal: nil,
				Line:    s.lineNum,
			}
		case TokenType(currChar).Is(LESS):
			if nextChar, exist := s.peek(); exist && TokenType(nextChar).Is(EQUAL) {
				currToken = Token{
					Type:    LESS_EQUAL,
					Lexeme:  string(LESS_EQUAL),
					Literal: nil,
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    LESS,
				Lexeme:  string(LESS),
				Literal: nil,
				Line:    s.lineNum,
			}
		case TokenType(currChar).Is(GREATER):
			if nextChar, exist := s.peek(); exist && TokenType(nextChar).Is(EQUAL) {
				currToken = Token{
					Type:    GREATER_EQUAL,
					Lexeme:  string(GREATER_EQUAL),
					Literal: nil,
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    GREATER,
				Lexeme:  string(GREATER),
				Literal: "null",
				Line:    s.lineNum,
			}
		case TokenType(currChar).Is(SLASH):
			if nextChar, exist := s.peek(); exist && TokenType(nextChar).Is(SLASH) {
				// comment encountered
				s.nextChar()

				for {
					n, e := s.peek()
					if !e || TokenType(n).Is(NEWLINE) {
						break
					}

					s.nextChar()
				}

				continue
			}

			currToken = Token{
				Type:    SLASH,
				Lexeme:  string(SLASH),
				Literal: nil,
				Line:    s.lineNum,
			}
		case TokenType(currChar).Is(SPACE) ||
			TokenType(currChar).Is(TAB):
			continue
		case TokenType(currChar).Is(QUOTE):
			currToken = Token{
				Type: STRING,
				Line: s.lineNum,
			}

			var bytes []byte

			for {
				n, e := s.peek()
				if !e {
					return nil, fmt.Errorf("[line %d] Error: Unterminated string.", currToken.Line+1)
				}

				if TokenType(n).Is(QUOTE) {
					currToken.Lexeme = fmt.Sprintf("\"%s\"", string(bytes))
					currToken.Literal = string(bytes)
					s.nextChar()
					break
				}

				bytes = append(bytes, n)

				s.nextChar()
			}
		case isNumeric(currChar):
			currToken = Token{
				Type:   NUMBER,
				Lexeme: string(currChar),
				Line:   s.lineNum,
			}

			for {
				n, e := s.peek()
				if !e || (!isNumeric(n) && !TokenType(n).Is(DOT)) {
					break
				}

				currToken.Lexeme += string(n)

				s.nextChar()
			}

			num, err := strconv.ParseFloat(currToken.Lexeme, 64)
			if err != nil {

				return nil, err
			}

			currToken.Literal = num
		case isAlphabet(currChar) || currChar == '_':
			currToken = Token{
				Type:    IDENTIFIER,
				Literal: nil,
				Lexeme:  string(currChar),
				Line:    s.lineNum,
			}

			for {
				n, e := s.peek()
				if !e || (!isAlphaNumeric(n) && n != '_') {
					break
				}

				currToken.Lexeme += string(n)

				s.nextChar()
			}

			if _, isKeyword := reservedWords[TokenType(currToken.Lexeme)]; isKeyword {
				currToken.Type = TokenType(currToken.Lexeme)
			}
		default:
			errStr := fmt.Sprintf("[line %d] Error: Unexpected character: ", s.lineNum+1) + string(currChar)
			return nil, errors.New(errStr)
		}

		return &currToken, nil
	}
}

func (s *Scanner) HasNext() bool {
	return !s.done
}

func (s *Scanner) nextChar() (byte, bool) {
	s.pos++

	if s.pos >= len(s.content) {
		s.pos = len(s.content)
		return 0, false
	}

	c := s.content[s.pos]

	return c, true
}

func (s *Scanner) peek() (byte, bool) {
	if s.pos+1 >= len(s.content) {
		return 0, false
	}

	return s.content[s.pos+1], true
}

func isNumeric(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphabet(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func isAlphaNumeric(b byte) bool {
	return isAlphabet(b) || isNumeric(b)
}

func stringifyNumOrNil(token *Token) string {
	if token.Literal == nil {
		return "null"
	}

	if v, ok := token.Literal.(float64); ok {
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.1f", v)
		}
	}

	return fmt.Sprintf("%v", token.Literal)
}
