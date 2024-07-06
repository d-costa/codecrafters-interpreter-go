package main

import (
	"fmt"
	"os"
)

type Token struct {
	TokenType string
	lexeme    string
	literal   string
}

var EOF Token = Token{TokenType: "EOF", lexeme: "", literal: "null"}

func (token *Token) setToken(tokenType string, lexeme string) {
	token.TokenType = tokenType
	token.lexeme = lexeme
	if tokenType == "STRING" {
		token.literal = lexeme[1 : len(lexeme)-1]
	} else if tokenType == "NUMBER" {
		token.literal = lexeme
	} else {
		token.literal = "null"
	}
}

func (token *Token) toString() string {
	return fmt.Sprintf("%s %s %s", token.TokenType, token.lexeme, token.literal)

}

func match(arr []byte, index int, ch byte) bool {
	return index < len(arr) && arr[index] == ch
}

func tokenizeFile(fileContents []byte) ([]Token, bool) {
	tokens := []Token{}
	hasLexicalError := false
	line_number := 1

	for i := 0; i < len(fileContents); i++ {
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
			if match(fileContents, i+1, '=') {
				newToken.setToken("EQUAL_EQUAL", "==")
				i += 1
			} else {
				newToken.setToken("EQUAL", "=")
			}
		case '!':
			if match(fileContents, i+1, '=') {
				newToken.setToken("BANG_EQUAL", "!=")
				i += 1
			} else {
				newToken.setToken("BANG", "!")
			}
		case '<':
			if match(fileContents, i+1, '=') {
				newToken.setToken("LESS_EQUAL", "<=")
				i += 1
			} else {
				newToken.setToken("LESS", "<")
			}
		case '>':
			if match(fileContents, i+1, '=') {
				newToken.setToken("GREATER_EQUAL", ">=")
				i += 1
			} else {
				newToken.setToken("GREATER", ">")
			}
		case '/':
			if match(fileContents, i+1, '/') {
				for i < len(fileContents) && !match(fileContents, i+1, '\n') {
					i++
				}
			} else {
				newToken.setToken("SLASH", "/")
			}
		case ' ':
			continue
		case '\r':
			continue
		case '\t':
			continue
		case '\n':
			line_number++
			continue
		case '"':
			terminated := false
			literal := ""
			for i = i + 1; i < len(fileContents); i++ {
				if match(fileContents, i, '"') {
					newToken.setToken("STRING", fmt.Sprintf("\"%s\"", literal))
					terminated = true
					break
				} else {
					literal += string(fileContents[i])
				}
			}

			if !terminated {
				msg := fmt.Errorf("[line %d] Error: Unterminated string.", line_number)
				fmt.Fprintln(os.Stderr, msg)
				hasLexicalError = true
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			numeric := ""
			numeric, i = extractNumeric(fileContents, i)

			//DECIMAL
			if i < len(fileContents) && fileContents[i] == '.' {
				numeric += "."
				decimal := ""
				decimal, i = extractNumeric(fileContents, i+1)
				numeric += decimal
			}
			newToken.setToken("NUMBER", numeric)

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

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func extractNumeric(fileContents []byte, i int) (string, int) {
	accum := ""
	for i < len(fileContents) && isDigit(fileContents[i]) {
		accum += string(fileContents[i])
		i++
	}
	return accum, i
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
