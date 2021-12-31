package main

import "fmt"

func main() {
	var lex Lexer
	lex.Tokens = make(chan Token, 20)
	lex.StartLexing("./test")

	for token := range lex.Tokens {
		fmt.Println(token)
	}
}
