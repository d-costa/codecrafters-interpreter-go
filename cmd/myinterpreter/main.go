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

func tokenizeFile(fileContents []byte) ([]Token, bool) {
	tokens := []Token{}
	hasLexicalError := false
	line_number := 1
	for i := 0; i < len(fileContents); i++ {
		if fileContents[i] == '\n' {
			line_number++
		}

		newToken := Token{}

		switch fileContents[i] {
		case '(':
			newToken.setToken("LEFT_PAREN", "(")
		case ')':
			newToken.setToken("RIGHT_PAREN", ")")
		case '{':
			newToken.setToken("LEFT_BRACE", "{")
		case '}':
			newToken.setToken("RIGHT_BRACE", "}")
		case ',':
			newToken.setToken("COMMA", ",")
		case '.':
			newToken.setToken("DOT", ".")
		case '-':
			newToken.setToken("MINUS", "-")
		case '+':
			newToken.setToken("PLUS", "+")
		case ';':
			newToken.setToken("SEMICOLON", ";")
		case '*':
			newToken.setToken("STAR", "*")
		case '=':
			if fileContents[i+1] == '=' {
				newToken.setToken("EQUAL_EQUAL", "==")
				i += 1
			} else {
				newToken.setToken("EQUAL", "=")
			}
		default:
			msg := fmt.Errorf("[line %d] Error: Unexpected character: %c", line_number, fileContents[i])
			fmt.Fprintln(os.Stderr, msg)
			hasLexicalError = true
		}

		if newToken.TokenType != "" {
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
