package main

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN  byte = '('
	RIGHT_PAREN byte = ')'
	LEFT_BRACE  byte = '{'
	RIGHT_BRACE byte = '}'
	COMMA       byte = ','
	DOT         byte = '.'
	SEMICOLON   byte = ';'
	PLUS        byte = '+'
	MINUS       byte = '-'
	STAR        byte = '*'
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

	for _, char := range fileContents {
		switch char {
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
		default:
			_, _ = fmt.Fprintf(os.Stderr, "[line 1] Error: Unexpected character: %s\n", string(char))
			lexErrFound = true
		}
	}

	fmt.Println("EOF  null")

	if lexErrFound {
		os.Exit(65)
	}
}
