package main

import (
	"fmt"
	"strings"
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
	Literal string
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type.Type(), t.Lexeme, t.Literal)
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
		if !ok {
			currToken = Token{
				Type:    EOF,
				Lexeme:  string(EOF),
				Literal: "null",
				Line:    s.lineNum,
			}

			s.done = true
		} else if TokenType(currChar) == LEFT_PAREN {
			currToken = Token{
				Type:    LEFT_PAREN,
				Lexeme:  string(LEFT_PAREN),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == RIGHT_PAREN {
			currToken = Token{
				Type:    RIGHT_PAREN,
				Lexeme:  string(RIGHT_PAREN),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == LEFT_BRACE {
			currToken = Token{
				Type:    LEFT_BRACE,
				Lexeme:  string(LEFT_BRACE),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == RIGHT_BRACE {
			currToken = Token{
				Type:    RIGHT_BRACE,
				Lexeme:  string(RIGHT_BRACE),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == COMMA {
			currToken = Token{
				Type:    COMMA,
				Lexeme:  string(COMMA),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == DOT {
			currToken = Token{
				Type:    DOT,
				Lexeme:  string(DOT),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == SEMICOLON {
			currToken = Token{
				Type:    SEMICOLON,
				Lexeme:  string(SEMICOLON),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == PLUS {
			currToken = Token{
				Type:    PLUS,
				Lexeme:  string(PLUS),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == MINUS {
			currToken = Token{
				Type:    MINUS,
				Lexeme:  string(MINUS),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == STAR {
			currToken = Token{
				Type:    STAR,
				Lexeme:  string(STAR),
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == NEWLINE {
			s.lineNum++
			continue
		} else if TokenType(currChar) == EQUAL {
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    EQUAL_EQUAL,
					Lexeme:  string(EQUAL_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
			} else {
				currToken = Token{
					Type:    EQUAL,
					Lexeme:  string(EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}
			}
		} else if TokenType(currChar) == BANG {
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    BANG_EQUAL,
					Lexeme:  string(BANG_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
			} else {
				currToken = Token{
					Type:    BANG,
					Lexeme:  string(BANG),
					Literal: "null",
					Line:    s.lineNum,
				}
			}
		} else if TokenType(currChar) == LESS {
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    LESS_EQUAL,
					Lexeme:  string(LESS_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
			} else {
				currToken = Token{
					Type:    LESS,
					Lexeme:  string(LESS),
					Literal: "null",
					Line:    s.lineNum,
				}
			}
		} else if TokenType(currChar) == GREATER {
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    GREATER_EQUAL,
					Lexeme:  string(GREATER_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
			} else {
				currToken = Token{
					Type:    GREATER,
					Lexeme:  string(GREATER),
					Literal: "null",
					Line:    s.lineNum,
				}
			}
		} else if TokenType(currChar) == SLASH {
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
				Literal: "null",
				Line:    s.lineNum,
			}
		} else if TokenType(currChar) == SPACE {
			continue
		} else if TokenType(currChar) == TAB {
			continue
		} else if TokenType(currChar) == QUOTE {
			currToken = Token{
				Type: STRING,
				Line: s.lineNum,
			}

			for {
				n, e := s.peek()
				if !e {
					return nil, fmt.Errorf("[line %d] Error: Unterminated string.", currToken.Line+1)
				}

				if TokenType(n) == QUOTE {
					currToken.Lexeme = fmt.Sprintf("\"%s\"", currToken.Literal)
					s.nextChar()
					break
				}

				currToken.Literal += string(n)

				s.nextChar()
			}
		} else if isNumeric(currChar) {
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

			currToken.Literal = currToken.Lexeme
			if !strings.Contains(currToken.Literal, ".") {
				currToken.Literal = fmt.Sprintf("%s.0", currToken.Literal)
			} else {
				idx := strings.Index(currToken.Literal, ".")

				var i int
				for i = len(currToken.Literal) - 1; i > idx && currToken.Literal[i] == '0'; i-- {
				}

				currToken.Literal = currToken.Literal[:i+1]
				if i == idx {
					currToken.Literal = fmt.Sprintf("%s0", currToken.Literal)
				}
			}
		} else if isAlphabet(currChar) || currChar == '_' {
			currToken = Token{
				Type:    IDENTIFIER,
				Literal: "null",
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

			if _, found := reservedWords[TokenType(currToken.Lexeme)]; found {
				currToken.Type = TokenType(currToken.Lexeme)
			}
		} else {
			return nil, fmt.Errorf("[line %d] Error: Unexpected character: %s\n", s.lineNum+1, string(currChar))
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
