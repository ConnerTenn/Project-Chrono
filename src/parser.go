package main

import (
	"fmt"
	"os"
)

func displayError(t Token) {

	fmt.Println("Unexpected Token: ", t)

	os.Exit(-1)
}

// use a prat / segmented parser design

func parseParam(lex *Lexer, head *AST) AST {

	return *head
}

// primary parsing function, can make assumptions about inital tokens
func Parse(lex *Lexer) AST {
	head := Module{}
	var t Token

	if !lex.ExpectNext(Iden) {
		t, _ = lex.GetNext()
		displayError(t)
	}

	t, _ = lex.GetNext()
	head.Name = t.Value

	if !lex.ExpectNext(LParen) {
		t, _ = lex.GetNext()
		displayError(t)
	}

	lex.GetNext() // drop LParen

	// build parameters
	//for {

	//}

	return head
}
