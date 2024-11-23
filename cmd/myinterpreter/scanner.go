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
		lexErrFound   bool
		stringStarted bool
		numberStarted bool
		currToken     Token
	)

	for {
		lineNum, line, ok := s.nextLine()
		if !ok {
			break
		}

	LOOP:
		for i := 0; i < len(line); i++ {
			if stringStarted && TokenType(line[i]) != QUOTE {
				currToken.Literal += string(line[i])
				continue
			}

			if numberStarted {
				if strings.Contains("0123456789.", string(line[i])) {
					currToken.Lexeme += string(line[i])
					continue
				}

				numberStarted = false

				currToken.Literal = currToken.Lexeme
				if !strings.Contains(currToken.Literal, ".") {
					currToken.Literal = fmt.Sprintf("%s.0", currToken.Literal)
				}

				fmt.Println(currToken.String())
			}

			switch TokenType(line[i]) {
			case LEFT_PAREN:
				currToken = Token{
					Type:    LEFT_PAREN,
					Lexeme:  string(LEFT_PAREN),
					Literal: "null",
					Line:    lineNum,
				}
			case RIGHT_PAREN:
				currToken = Token{
					Type:    RIGHT_PAREN,
					Lexeme:  string(RIGHT_PAREN),
					Literal: "null",
					Line:    lineNum,
				}
			case LEFT_BRACE:
				currToken = Token{
					Type:    LEFT_BRACE,
					Lexeme:  string(LEFT_BRACE),
					Literal: "null",
					Line:    lineNum,
				}
			case RIGHT_BRACE:
				currToken = Token{
					Type:    RIGHT_BRACE,
					Lexeme:  string(RIGHT_BRACE),
					Literal: "null",
					Line:    lineNum,
				}
			case COMMA:
				currToken = Token{
					Type:    COMMA,
					Lexeme:  string(COMMA),
					Literal: "null",
					Line:    lineNum,
				}
			case DOT:
				currToken = Token{
					Type:    DOT,
					Lexeme:  string(DOT),
					Literal: "null",
					Line:    lineNum,
				}
			case SEMICOLON:
				currToken = Token{
					Type:    SEMICOLON,
					Lexeme:  string(SEMICOLON),
					Literal: "null",
					Line:    lineNum,
				}
			case PLUS:
				currToken = Token{
					Type:    PLUS,
					Lexeme:  string(PLUS),
					Literal: "null",
					Line:    lineNum,
				}
			case MINUS:
				currToken = Token{
					Type:    MINUS,
					Lexeme:  string(MINUS),
					Literal: "null",
					Line:    lineNum,
				}
			case STAR:
				currToken = Token{
					Type:    STAR,
					Lexeme:  string(STAR),
					Literal: "null",
					Line:    lineNum,
				}
			case EQUAL:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    EQUAL_EQUAL,
						Lexeme:  string(EQUAL_EQUAL),
						Literal: "null",
						Line:    lineNum,
					}

					i++
					break
				}

				currToken = Token{
					Type:    EQUAL,
					Lexeme:  string(EQUAL),
					Literal: "null",
					Line:    lineNum,
				}
			case BANG:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    BANG_EQUAL,
						Lexeme:  string(BANG_EQUAL),
						Literal: "null",
						Line:    lineNum,
					}

					i++
					break
				}

				currToken = Token{
					Type:    BANG,
					Lexeme:  string(BANG),
					Literal: "null",
					Line:    lineNum,
				}
			case LESS:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    LESS_EQUAL,
						Lexeme:  string(LESS_EQUAL),
						Literal: "null",
						Line:    lineNum,
					}

					i++
					break
				}

				currToken = Token{
					Type:    LESS,
					Lexeme:  string(LESS),
					Literal: "null",
					Line:    lineNum,
				}
			case GREATER:
				if i+1 < len(line) && TokenType(line[i+1]) == EQUAL {
					currToken = Token{
						Type:    GREATER_EQUAL,
						Lexeme:  string(GREATER_EQUAL),
						Literal: "null",
						Line:    lineNum,
					}

					i++
					break
				}

				currToken = Token{
					Type:    GREATER,
					Lexeme:  string(GREATER),
					Literal: "null",
					Line:    lineNum,
				}
			case SLASH:
				if i+1 < len(line) && TokenType(line[i+1]) == SLASH {
					break LOOP
				}

				currToken = Token{
					Type:    SLASH,
					Lexeme:  string(SLASH),
					Literal: "null",
					Line:    lineNum,
				}
			case SPACE, TAB:
				continue
			case QUOTE:
				if !stringStarted {
					stringStarted = true
					currToken = Token{
						Type: STRING,
						Line: lineNum,
					}

					continue
				}

				// we found the matching quote.
				stringStarted = false
				currToken.Lexeme = fmt.Sprintf("\"%s\"", currToken.Literal)
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				numberStarted = true
				currToken = Token{
					Type:   NUMBER,
					Lexeme: strconv.Itoa(int(line[i])),
					Line:   lineNum,
				}

				continue
			default:
				_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNum+1, string(line[i]))
				lexErrFound = true
				continue
			}

			fmt.Println(currToken.String())
		}
	}

	if stringStarted {
		// if we found an unterminated string
		_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.", currToken.Line+1)
	}

	fmt.Println((&Token{
		Type:    EOF,
		Lexeme:  string(EOF),
		Literal: "null",
	}).String())

	if lexErrFound || stringStarted {
		os.Exit(65)
	}
}

func (s *Scanner) nextLine() (int, string, bool) {
	if s.currLine >= len(s.lines) {
		return 0, "", false
	}

	c := s.currLine
	s.currLine++

	return c, s.lines[c], true
}
