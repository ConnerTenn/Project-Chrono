package main

import "fmt"

func main() {
	lex, err := NewLexer("./test")
	if err != nil {
		fmt.Println("Error:", err)
	}
	go lex.Tokenizer()

	for token, ok := lex.GetNext(); ok; token, ok = lex.GetNext() {
		fmt.Println(token)
	}
}
