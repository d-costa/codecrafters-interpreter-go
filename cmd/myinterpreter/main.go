package main

import (
	"fmt"
	"os"
)

type Token struct {
	TokenType string
	lexeme    string
}

var EOF Token = Token{TokenType: "EOF", lexeme: ""}

func (token *Token) setToken(tokenType string, lexeme string) {
	token.TokenType = tokenType
	token.lexeme = lexeme
}

func (token *Token) toString() string {
	return fmt.Sprintf("%s %s null", token.TokenType, token.lexeme)

}

func scanToken(ch string) Token {
	token := Token{}
	switch ch {
	case "(":
		token.setToken("LEFT_PAREN", "(")
	case ")":
		token.setToken("RIGHT_PAREN", ")")
	default:
		fmt.Println("Unknown char: ", ch)
	}
	return token
}

func tokenizeFile(fileContents []byte) []Token {
	tokens := []Token{}
	for i := 0; i < len(fileContents); i++ {
		newToken := scanToken(string(fileContents[i]))
		tokens = append(tokens, newToken)
	}

	tokens = append(tokens, EOF)
	return tokens
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

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

	if len(fileContents) > 0 {
		tokens := tokenizeFile(fileContents)
		for _, token := range tokens {
			fmt.Println(token.toString())
		}
	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
}
