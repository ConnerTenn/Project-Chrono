package main

import "fmt"

func main() {
	var lex Lexer
	lex.Tokens = NewQueue()
	lex.StartLexing("./test")

	for token, ok := lex.Tokens.GetNext(); ok; token, ok = lex.Tokens.GetNext() {
		fmt.Println(token)
	}
}
