package Parser

import (
	"strconv"

	AST "github.com/ConnerTenn/Project-Chrono/AST"
	L "github.com/ConnerTenn/Project-Chrono/Lexer"
)

// use a prat / segmented parser design

func parseParam(lex *L.Lexer) AST.Signal {
	t, _ := lex.GetNext()
	curParam := AST.Signal{}

	// set / get direction
	switch t.Value {
	case "in":
		curParam.Dir = AST.In
	case "out":
		curParam.Dir = AST.Out
	case "inout":
		curParam.Dir = AST.Inout
	}

	// set / get param type
	if lex.ExpectNext(L.Spec) {
		t, _ = lex.GetNext()
		curParam.Type = AST.Wire
		switch t.Value {
		case "reg":
			curParam.Type = AST.Reg
		}
	}

	// set / get bit width
	if lex.ExpectNext(L.LBrace) {
		lex.GetNext()
		t, _ = lex.GetNext()

		if t.Type != L.Literal {
			displayError(t)
		}

		curParam.Width, _ = strconv.Atoi(t.Value)

		t, _ = lex.GetNext()

		if t.Type != L.RBrace {
			displayError(t)
		}
	}

	// get / set name
	t, _ = lex.GetNext()

	if t.Type != L.Iden {
		displayError(t)
	}

	curParam.Name = t.Value

	return curParam
}

func parseModule(lex *L.Lexer) AST.Module {
	head := AST.Module{}
	var t L.Token

	t, _ = lex.GetNext() // assumtion first token will be an iden
	head.Name = t.Value

	if !lex.ExpectNext(L.LParen) {
		t, _ = lex.GetNext()
		displayError(t)
	}

	lex.GetNext() // drop LParen

	// build parameters
	for {
		if !lex.ExpectNext(L.Direction) {
			break // if the next value isn't a direction, there are no parameters
		}

		// add to params list
		head.Params = append(head.Params, parseParam(lex))

		t, _ = lex.GetNext()

		// check if end of param list
		if t.Type == L.RParen {
			break
		} else if t.Type != L.Comma {
			displayError(t)
		}
	}

	// parse code

	t, _ = lex.PeekNext()
	if !lex.ExpectNext(L.LCurly) {
		displayError(t)
	}
	// drop LCurly
	lex.GetNext()

	// FIX ME: massive assumptions made to test expression structs
	// var expTop AST.Expression
	var expHead AST.Expression

	for lex.ExpectNext(L.Iden) {
		t, _ = lex.GetNext()

		if t.Type == L.Iden && lex.ExpectNext(L.Asmt) { // parse assignment
			asmt := AST.AssignmentExpression{}
			asmt.Name = t.Value

			lex.GetNext() // drop asmt token

			asmt.RHS = parseMathExpression(lex)

			expHead = &asmt
		} else {
			break
		}

		head.Elements = append(head.Elements, expHead)

	}

	// === FIX ME!!! ===
	// Temp code
	//Consume all until next RCurly
	for !lex.ExpectNext(L.RCurly) {
		_, _ = lex.GetNext()
		//fmt.Println(t)
	}
	//Consume RCurly
	lex.GetNext()

	return head
}

// todo: remove assumptions & add error handling
// primary parsing function, can make assumptions about inital tokens
func parseBlock(lex *L.Lexer) AST.Block {
	block := AST.Block{}

	for !(lex.ExpectNext(L.EOL) || lex.ExpectNext(L.RCurly)) {
		var next AST.AST
		if lex.ExpectNext(L.Iden) {
			m := parseModule(lex)
			next = &m
		} else {
			break
		}

		block.Elements = append(block.Elements, next)
	}

	//Consume RCurly
	if lex.ExpectNext(L.RCurly) {
		lex.GetNext()
	}

	return block
}

func Parse(lex *L.Lexer) AST.AST {
	ast := parseBlock(lex)
	return &ast
}
