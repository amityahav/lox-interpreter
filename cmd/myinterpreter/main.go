package main

import (
	"bytes"
	"fmt"
	"os"
)

const (
	LEFT_PAREN    byte   = '('
	RIGHT_PAREN   byte   = ')'
	LEFT_BRACE    byte   = '{'
	RIGHT_BRACE   byte   = '}'
	COMMA         byte   = ','
	DOT           byte   = '.'
	SEMICOLON     byte   = ';'
	PLUS          byte   = '+'
	MINUS         byte   = '-'
	STAR          byte   = '*'
	EQUAL         byte   = '='
	EQUAL_EQUAL   string = "=="
	BANG          byte   = '!'
	BANG_EQUAL    string = "!="
	LESS          byte   = '<'
	LESS_EQUAL    string = "<="
	GREATER       byte   = '>'
	GREATER_EQUAL string = ">="
	SLASH         byte   = '/'
	NEWLINE       byte   = '\n'
	SPACE         byte   = ' '
	TAB           byte   = '\t'
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var lexErrFound bool

	lines := bytes.Split(fileContents, []byte{NEWLINE})

	for lineNum, line := range lines {
	LOOP:
		for i := 0; i < len(line); i++ {
			switch fileContents[i] {
			case LEFT_PAREN:
				fmt.Printf("LEFT_PAREN ( null\n")
			case RIGHT_PAREN:
				fmt.Printf("RIGHT_PAREN ) null\n")
			case LEFT_BRACE:
				fmt.Printf("LEFT_BRACE { null\n")
			case RIGHT_BRACE:
				fmt.Printf("RIGHT_BRACE } null\n")
			case COMMA:
				fmt.Printf("COMMA , null\n")
			case DOT:
				fmt.Printf("DOT . null\n")
			case SEMICOLON:
				fmt.Printf("SEMICOLON ; null\n")
			case PLUS:
				fmt.Printf("PLUS + null\n")
			case MINUS:
				fmt.Printf("MINUS - null\n")
			case STAR:
				fmt.Printf("STAR * null\n")
			case EQUAL:
				if i+1 < len(fileContents) && fileContents[i+1] == EQUAL {
					fmt.Printf("EQUAL_EQUAL == null\n")
					i += 1
					continue
				}

				fmt.Printf("EQUAL = null\n")
			case BANG:
				if i+1 < len(fileContents) && fileContents[i+1] == EQUAL {
					fmt.Printf("BANG_EQUAL != null\n")
					i += 1
					continue
				}

				fmt.Printf("BANG ! null\n")
			case LESS:
				if i+1 < len(fileContents) && fileContents[i+1] == EQUAL {
					fmt.Printf("LESS_EQUAL <= null\n")
					i += 1
					continue
				}

				fmt.Printf("LESS < null\n")
			case GREATER:
				if i+1 < len(fileContents) && fileContents[i+1] == EQUAL {
					fmt.Printf("GREATER_EQUAL >= null\n")
					i += 1
					continue
				}

				fmt.Printf("GREATER > null\n")
			case SLASH:
				if i+1 < len(fileContents) && fileContents[i+1] == SLASH {
					break LOOP
				}

				fmt.Printf("SLASH / null\n")
			case SPACE, TAB:
			default:
				_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNum+1, string(fileContents[i]))
				lexErrFound = true
			}
		}
	}

	fmt.Println("EOF  null")

	if lexErrFound {
		os.Exit(65)
	}
}
