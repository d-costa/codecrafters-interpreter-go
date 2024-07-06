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

func scanToken(ch string) (token Token, err error) {
	switch ch {
	case "(":
		token.setToken("LEFT_PAREN", "(")
	case ")":
		token.setToken("RIGHT_PAREN", ")")
	case "{":
		token.setToken("LEFT_BRACE", "{")
	case "}":
		token.setToken("RIGHT_BRACE", "}")
	case ",":
		token.setToken("COMMA", ",")
	case ".":
		token.setToken("DOT", ".")
	case "-":
		token.setToken("MINUS", "-")
	case "+":
		token.setToken("PLUS", "+")
	case ";":
		token.setToken("SEMICOLON", ";")
	case "*":
		token.setToken("STAR", "*")
	default:
		err = fmt.Errorf("Unexpected character: %s", string(ch))
	}
	return token, err
}

func tokenizeFile(fileContents []byte) ([]Token, bool) {
	tokens := []Token{}
	hasLexicalError := false
	line_number := 1
	for i := 0; i < len(fileContents); i++ {
		if fileContents[i] == '\n' {
			line_number++
		}
		newToken, err := scanToken(string(fileContents[i]))
		if err != nil {
			msg := fmt.Errorf("[line %d] Error: %s", line_number, err)
			fmt.Fprintln(os.Stderr, msg)
			hasLexicalError = true
		} else {
			tokens = append(tokens, newToken)
		}
	}

	tokens = append(tokens, EOF)
	return tokens, hasLexicalError
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

	tokens, hasLexicalError := tokenizeFile(fileContents)
	for _, token := range tokens {
		fmt.Println(token.toString())
	}

	if hasLexicalError {
		os.Exit(65)
	}
}
