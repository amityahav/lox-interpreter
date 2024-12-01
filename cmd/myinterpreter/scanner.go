package main

import (
	"errors"
	"fmt"
	"strconv"
)

type TokenType string

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
	EOF           TokenType = ""
)

var reservedWords = map[TokenType]struct{}{
	"and":    {},
	"class":  {},
	"else":   {},
	"false":  {},
	"for":    {},
	"fun":    {},
	"if":     {},
	"nil":    {},
	"or":     {},
	"print":  {},
	"return": {},
	"super":  {},
	"this":   {},
	"true":   {},
	"var":    {},
	"while":  {},
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
	case "and":
		return "AND"
	case "class":
		return "CLASS"
	case "else":
		return "ELSE"
	case "false":
		return "FALSE"
	case "for":
		return "FOR"
	case "fun":
		return "FUN"
	case "if":
		return "IF"
	case "nil":
		return "NIL"
	case "or":
		return "OR"
	case "print":
		return "PRINT"
	case "return":
		return "RETURN"
	case "super":
		return "SUPER"
	case "this":
		return "THIS"
	case "true":
		return "TRUE"
	case "var":
		return "VAR"
	case "while":
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
		case TokenType(currChar) == LEFT_PAREN ||
			TokenType(currChar) == RIGHT_PAREN ||
			TokenType(currChar) == LEFT_BRACE ||
			TokenType(currChar) == RIGHT_BRACE ||
			TokenType(currChar) == COMMA ||
			TokenType(currChar) == DOT ||
			TokenType(currChar) == SEMICOLON ||
			TokenType(currChar) == PLUS ||
			TokenType(currChar) == MINUS ||
			TokenType(currChar) == STAR:
			currToken = Token{
				Type:    TokenType(currChar),
				Lexeme:  string(currChar),
				Literal: nil,
				Line:    s.lineNum,
			}
		case TokenType(currChar) == NEWLINE:
			s.lineNum++
			continue
		case TokenType(currChar) == EQUAL:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
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
		case TokenType(currChar) == BANG:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
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
		case TokenType(currChar) == LESS:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
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
		case TokenType(currChar) == GREATER:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
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
		case TokenType(currChar) == SLASH:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == SLASH {
				// comment encountered
				s.nextChar()

				for {
					n, e := s.peek()
					if !e || TokenType(n) == NEWLINE {
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
		case TokenType(currChar) == SPACE ||
			TokenType(currChar) == TAB:
			continue
		case TokenType(currChar) == QUOTE:
			currToken = Token{
				Type: STRING,
				Line: s.lineNum,
			}

			var str string

			for {
				n, e := s.peek()
				if !e {
					return nil, fmt.Errorf("[line %d] Error: Unterminated string.", currToken.Line+1)
				}

				if TokenType(n) == QUOTE {
					currToken.Lexeme = fmt.Sprintf("\"%s\"", str)
					currToken.Literal = str
					s.nextChar()
					break
				}

				str += string(n)

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
				if !e || (!isNumeric(n) && TokenType(n) != DOT) {
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
