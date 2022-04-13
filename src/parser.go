package main

import (
	"fmt"
	"os"
	"strconv"
)

func displayError(t Token) {

	fmt.Println("Unexpected Token: ", t)

	os.Exit(-1)
}

// use a prat / segmented parser design

func parseParam(lex *Lexer) Parameter {
	t, _ := lex.GetNext()
	curParam := Parameter{}

	// set / get direction
	switch t.Value {
	case "in":
		{
			curParam.Dir = In
		}
	case "out":
		{
			curParam.Dir = Out
		}
	case "inout":
		{
			curParam.Dir = Inout
		}
	}

	// set / get param type
	if lex.ExpectNext(Spec) {
		t, _ = lex.GetNext()
		switch t.Value {
		case "reg":
			{
				curParam.Type = Reg
			}
		}
	}

	// set / get bit width
	if lex.ExpectNext(LBrace) {
		lex.GetNext()
		t, _ = lex.GetNext()

		if t.Type != Literal {
			displayError(t)
		}

		curParam.Width, _ = strconv.Atoi(t.Value)

		t, _ = lex.GetNext()

		if t.Type != RBrace {
			displayError(t)
		}
	}

	// get / set name
	t, _ = lex.GetNext()

	if t.Type != Iden {
		displayError(t)
	}

	curParam.Name = t.Value

	return curParam
}

func parseModule(lex *Lexer) Module {
	head := Module{}
	var t Token

	t, _ = lex.GetNext()
	head.Name = t.Value

	if !lex.ExpectNext(LParen) {
		t, _ = lex.GetNext()
		displayError(t)
	}

	lex.GetNext() // drop LParen

	// build parameters
	for {
		if !lex.ExpectNext(Direction) {
			break // if the next value isn't a direction, there are no parameters
		}

		// add to params list
		head.Params = append(head.Params, parseParam(lex))

		t, _ = lex.GetNext()

		// check if end of param list
		if t.Type == RParen {
			break
		} else if t.Type != Comma {
			displayError(t)
		}
	}

	return head
}

// todo: remove assumptions & error handling
// primary parsing function, can make assumptions about inital tokens
func Parse(lex *Lexer) AST {

	if !lex.ExpectNext(Iden) {
		t, _ := lex.GetNext()
		displayError(t)
	}

	head := parseModule(lex)

	return head
}
