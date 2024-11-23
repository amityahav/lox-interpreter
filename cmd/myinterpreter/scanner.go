package main

import (
	"fmt"
	"os"
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
	Literal string
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %s %s", t.Type.Type(), t.Type, t.Literal)
}

type Scanner struct {
	lines    []string
	currLine int
}

func NewScanner(content []byte) *Scanner {
	contentStr := string(content)
	lines := strings.Split(contentStr, string(NEWLINE))
	s := Scanner{lines: lines}

	return &s
}

func (s *Scanner) Scan() {
	var (
		lexErrFound bool
		currToken   Token
	)

	for {
		lineNum, line, ok := s.nextLine()
		if !ok {
			break
		}

	LOOP:
		for i := 0; i < len(line); i++ {
			switch TokenType(line[i]) {
			case LEFT_PAREN:
				currToken = Token{
					Type:    LEFT_PAREN,
					Literal: "null",
				}
			case RIGHT_PAREN:
				currToken = Token{
					Type:    RIGHT_PAREN,
					Literal: "null",
				}
			case LEFT_BRACE:
				currToken = Token{
					Type:    LEFT_BRACE,
					Literal: "null",
				}
			case RIGHT_BRACE:
				currToken = Token{
					Type:    RIGHT_BRACE,
					Literal: "null",
				}
			case COMMA:
				currToken = Token{
					Type:    COMMA,
					Literal: "null",
				}
			case DOT:
				currToken = Token{
					Type:    DOT,
					Literal: "null",
				}
			case SEMICOLON:
				currToken = Token{
					Type:    SEMICOLON,
					Literal: "null",
				}
			case PLUS:
				currToken = Token{
					Type:    PLUS,
					Literal: "null",
				}
			case MINUS:
				currToken = Token{
					Type:    MINUS,
					Literal: "null",
				}
			case STAR:
				currToken = Token{
					Type:    STAR,
					Literal: "null",
				}
			case EQUAL:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    EQUAL_EQUAL,
						Literal: "null",
					}

					i += 1
					continue
				}

				currToken = Token{
					Type:    EQUAL,
					Literal: "null",
				}
			case BANG:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    BANG_EQUAL,
						Literal: "null",
					}

					i += 1
					continue
				}

				currToken = Token{
					Type:    BANG,
					Literal: "null",
				}
			case LESS:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    LESS_EQUAL,
						Literal: "null",
					}

					i += 1
					continue
				}

				currToken = Token{
					Type:    LESS,
					Literal: "null",
				}
			case GREATER:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    GREATER_EQUAL,
						Literal: "null",
					}

					i += 1
					continue
				}

				currToken = Token{
					Type:    GREATER,
					Literal: "null",
				}
			case SLASH:
				if i+1 < len(line) && TokenType(line[i+1]) == SLASH {
					break LOOP
				}

				currToken = Token{
					Type:    SLASH,
					Literal: "null",
				}
			case SPACE, TAB:
			default:
				_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNum+1, string(line[i]))
				lexErrFound = true
			}

			fmt.Println(currToken.String())
		}
	}

	fmt.Println((&Token{
		Type:    EOF,
		Literal: "null",
	}).String())

	if lexErrFound {
		os.Exit(65)
	}
}

func (s *Scanner) nextLine() (int, string, bool) {
	if s.currLine >= len(s.lines) {
		return 0, "", false
	}

	return s.currLine, s.lines[s.currLine], true
}
