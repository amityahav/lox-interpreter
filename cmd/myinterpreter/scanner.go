package main

import (
	"fmt"
	"os"
	"strconv"
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
	EOF           TokenType = ""
)

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
	case EOF:
		return "EOF"
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
}

func NewScanner(content []byte) *Scanner {
	s := Scanner{content: content}

	return &s
}

func (s *Scanner) Scan() {
	var (
		lexErrFound bool
		currToken   Token
	)

LOOP:
	for {
		currChar, ok := s.nextChar()
		if !ok {
			currChar = string(EOF)
		}

		switch TokenType(currChar) {
		case LEFT_PAREN:
			currToken = Token{
				Type:    LEFT_PAREN,
				Lexeme:  string(LEFT_PAREN),
				Literal: "null",
				Line:    s.lineNum,
			}
		case RIGHT_PAREN:
			currToken = Token{
				Type:    RIGHT_PAREN,
				Lexeme:  string(RIGHT_PAREN),
				Literal: "null",
				Line:    s.lineNum,
			}
		case LEFT_BRACE:
			currToken = Token{
				Type:    LEFT_BRACE,
				Lexeme:  string(LEFT_BRACE),
				Literal: "null",
				Line:    s.lineNum,
			}
		case RIGHT_BRACE:
			currToken = Token{
				Type:    RIGHT_BRACE,
				Lexeme:  string(RIGHT_BRACE),
				Literal: "null",
				Line:    s.lineNum,
			}
		case COMMA:
			currToken = Token{
				Type:    COMMA,
				Lexeme:  string(COMMA),
				Literal: "null",
				Line:    s.lineNum,
			}
		case DOT:
			currToken = Token{
				Type:    DOT,
				Lexeme:  string(DOT),
				Literal: "null",
				Line:    s.lineNum,
			}
		case SEMICOLON:
			currToken = Token{
				Type:    SEMICOLON,
				Lexeme:  string(SEMICOLON),
				Literal: "null",
				Line:    s.lineNum,
			}
		case PLUS:
			currToken = Token{
				Type:    PLUS,
				Lexeme:  string(PLUS),
				Literal: "null",
				Line:    s.lineNum,
			}
		case MINUS:
			currToken = Token{
				Type:    MINUS,
				Lexeme:  string(MINUS),
				Literal: "null",
				Line:    s.lineNum,
			}
		case STAR:
			currToken = Token{
				Type:    STAR,
				Lexeme:  string(STAR),
				Literal: "null",
				Line:    s.lineNum,
			}
		case NEWLINE:
			s.lineNum++
			continue
		case EQUAL:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    EQUAL_EQUAL,
					Lexeme:  string(EQUAL_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    EQUAL,
				Lexeme:  string(EQUAL),
				Literal: "null",
				Line:    s.lineNum,
			}
		case BANG:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    BANG_EQUAL,
					Lexeme:  string(BANG_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    BANG,
				Lexeme:  string(BANG),
				Literal: "null",
				Line:    s.lineNum,
			}
		case LESS:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    LESS_EQUAL,
					Lexeme:  string(LESS_EQUAL),
					Literal: "null",
					Line:    s.lineNum,
				}

				s.nextChar()
				break
			}

			currToken = Token{
				Type:    LESS,
				Lexeme:  string(LESS),
				Literal: "null",
				Line:    s.lineNum,
			}
		case GREATER:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == EQUAL {
				currToken = Token{
					Type:    GREATER_EQUAL,
					Lexeme:  string(GREATER_EQUAL),
					Literal: "null",
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
		case SLASH:
			if nextChar, exist := s.peek(); exist && TokenType(nextChar) == SLASH {
				// comment encountered
				s.nextChar()

				for n, e := s.peek(); e && TokenType(n) != NEWLINE; s.nextChar() {
				}

				continue
			}

			currToken = Token{
				Type:    SLASH,
				Lexeme:  string(SLASH),
				Literal: "null",
				Line:    s.lineNum,
			}
		case SPACE, TAB:
			continue
		case QUOTE:
			currToken = Token{
				Type: STRING,
				Line: s.lineNum,
			}

			for {
				n, e := s.peek()
				if !e {
					_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.", currToken.Line+1)
					lexErrFound = true
					break LOOP
				}

				if TokenType(n) == QUOTE {
					currToken.Lexeme = fmt.Sprintf("\"%s\"", currToken.Literal)
					break
				}

				currToken.Literal += n
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			currToken = Token{
				Type:   NUMBER,
				Lexeme: currChar,
				Line:   s.lineNum,
			}

			for n, e := s.peek(); e && strings.Contains("0123456789.", n); s.nextChar() {
				currToken.Lexeme += n
			}

			currToken.Literal = currToken.Lexeme
			if !strings.Contains(currToken.Literal, ".") {
				currToken.Literal = fmt.Sprintf("%s.0", currToken.Literal)
			} else {
				idx := strings.Index(currToken.Literal, ".")
				d, err := strconv.Atoi(currToken.Literal[idx+1:])
				if err != nil {
					panic(err)
				}

				if d == 0 {
					currToken.Literal = fmt.Sprintf("%s.0", currToken.Literal[:idx])
				}
			}
		case EOF:
			currToken = Token{
				Type:    EOF,
				Lexeme:  string(EOF),
				Literal: "null",
				Line:    s.lineNum,
			}

			fmt.Println(currToken.String())
			break LOOP
		default:
			_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", s.lineNum+1, currChar)
			lexErrFound = true
			continue
		}

		fmt.Println(currToken.String())
	}

	if lexErrFound {
		os.Exit(65)
	}
}

func (s *Scanner) nextChar() (string, bool) {
	if s.pos >= len(s.content) {
		return "", false
	}

	c := s.content[s.pos]
	s.pos++

	return string(c), true
}

func (s *Scanner) peek() (string, bool) {
	if s.pos+1 >= len(s.content) {
		return "", false
	}

	return string(s.content[s.pos+1]), true
}
