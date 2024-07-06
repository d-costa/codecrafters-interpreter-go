package main

import (
	"fmt"
	"os"
	"strings"
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
		if strings.HasSuffix(lexeme, ".") { // Ends with a dot
			token.literal = fmt.Sprintf("%s0", lexeme)
		} else if !strings.Contains(lexeme, ".") { // Does not contain a dot
			token.literal = fmt.Sprintf("%s.0", lexeme)
		} else if strings.Contains(lexeme, ".") && strings.HasSuffix(lexeme, "0") { // Trailing decimal zeroes
			token.literal = strings.TrimRight(lexeme, "0") + "0"
		} else {
			token.literal = lexeme
		}
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

		switch {
		case '(' == fileContents[i]:
			newToken.setToken("LEFT_PAREN", "(")
		case ')' == fileContents[i]:
			newToken.setToken("RIGHT_PAREN", ")")
		case '{' == fileContents[i]:
			newToken.setToken("LEFT_BRACE", "{")
		case '}' == fileContents[i]:
			newToken.setToken("RIGHT_BRACE", "}")
		case ',' == fileContents[i]:
			newToken.setToken("COMMA", ",")
		case '.' == fileContents[i]:
			newToken.setToken("DOT", ".")
		case '-' == fileContents[i]:
			newToken.setToken("MINUS", "-")
		case '+' == fileContents[i]:
			newToken.setToken("PLUS", "+")
		case ';' == fileContents[i]:
			newToken.setToken("SEMICOLON", ";")
		case '*' == fileContents[i]:
			newToken.setToken("STAR", "*")
		case '=' == fileContents[i]:
			if match(fileContents, i+1, '=') {
				newToken.setToken("EQUAL_EQUAL", "==")
				i += 1
			} else {
				newToken.setToken("EQUAL", "=")
			}
		case '!' == fileContents[i]:
			if match(fileContents, i+1, '=') {
				newToken.setToken("BANG_EQUAL", "!=")
				i += 1
			} else {
				newToken.setToken("BANG", "!")
			}
		case '<' == fileContents[i]:
			if match(fileContents, i+1, '=') {
				newToken.setToken("LESS_EQUAL", "<=")
				i += 1
			} else {
				newToken.setToken("LESS", "<")
			}
		case '>' == fileContents[i]:
			if match(fileContents, i+1, '=') {
				newToken.setToken("GREATER_EQUAL", ">=")
				i += 1
			} else {
				newToken.setToken("GREATER", ">")
			}
		case '/' == fileContents[i]:
			if match(fileContents, i+1, '/') {
				for i < len(fileContents) && !match(fileContents, i+1, '\n') {
					i++
				}
			} else {
				newToken.setToken("SLASH", "/")
			}
		case ' ' == fileContents[i]:
			continue
		case '\r' == fileContents[i]:
			continue
		case '\t' == fileContents[i]:
			continue
		case '\n' == fileContents[i]:
			line_number++
			continue
		case '"' == fileContents[i]:
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
		case isDigit(fileContents[i]):
			numeric := ""
			numeric, i = extractNumeric(fileContents, i)

			// For the next dot character to be considered a decimal, it must be followed by a digit
			if i+1 < len(fileContents) && fileContents[i] == '.' && isDigit(fileContents[i+1]) {
				numeric += "."
				i++

				decimal := ""
				decimal, i = extractNumeric(fileContents, i)
				numeric += decimal
			}

			i-- // We are one step ahead
			newToken.setToken("NUMBER", numeric)

		case isIdentifierCharacter(fileContents[i]):
			identifier := ""
			for i < len(fileContents) && (isIdentifierCharacter(fileContents[i]) || isDigit(fileContents[i])) {
				identifier += string(fileContents[i])
				i++
			}
			newToken.setToken("IDENTIFIER", identifier)
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

func isIdentifierCharacter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
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
