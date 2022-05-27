package main

import (
	"fmt"
	"os"
	"strconv"

	AST "github.com/ConnerTenn/Project-Chrono/AST"
)

func displayError(t Token) {

	fmt.Println("Unexpected Token: ", t)

	os.Exit(-1)
}

// use a prat / segmented parser design

func parseParam(lex *Lexer) AST.Parameter {
	t, _ := lex.GetNext()
	curParam := AST.Parameter{}

	// set / get direction
	switch t.Value {
	case "in":
		{
			curParam.Dir = AST.In
		}
	case "out":
		{
			curParam.Dir = AST.Out
		}
	case "inout":
		{
			curParam.Dir = AST.Inout
		}
	}

	// set / get param type
	if lex.ExpectNext(Spec) {
		t, _ = lex.GetNext()
		curParam.Type = AST.Wire
		switch t.Value {
		case "reg":
			curParam.Type = AST.Reg
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

func parseMathExpression(lex *Lexer) AST.MathExpression {
	math := AST.MathExpression{}

	parseValue := func() AST.ValueExpression {
		value := AST.ValueExpression{}
		t, _ := lex.GetNext()

		value.Value = t.Value

		if t.Type == Iden {
			value.Var = true
		} else if t.Type == Literal {
			value.Var = false
		} else {
			displayError(t)
		}

		return value
	}

	parseOperation := func() AST.Operation {
		var op AST.Operation
		t, _ := lex.GetNext()

		switch t.Value {
		case "+":
			op = AST.Add
		case "-":
			op = AST.Sub
		case "*":
			op = AST.Multi
		case "/":
			op = AST.Div
		case "<<":
			op = AST.LShift
		case ">>":
			op = AST.RShift
		}

		return op
	}

	// assume lhs and rhs are values and proper syntax is given
	math.LHS = parseValue()

	math.Op = parseOperation()

	math.RHS = parseValue()

	return math
}

func parseModule(lex *Lexer) AST.Module {
	head := AST.Module{}
	var t Token

	t, _ = lex.GetNext() // assumtion first token will be an iden
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

	// parse code
	t, _ = lex.GetNext()
	if t.Type != LCurly {
		displayError(t)
	}

	// FIX ME: massive assumptions made to test expression structs
	var expTop AST.Expression
	var expHead AST.Expression

	t, _ = lex.GetNext() // drop LCurly

	for {
		if t.Type == Iden && lex.ExpectNext(Asmt) { // parse assignment
			asmt := AST.AssignmentExpression{}
			asmt.Name = t.Value

			lex.GetNext() // drop asmt token

			asmt.RHS = parseMathExpression(lex)

			expHead = &asmt
		} else {
			break
		}

		if expTop == nil {
			expTop = expHead
		}
	}

	head.Elements = append(head.Elements, expTop)

	// === FIX ME!!! ===
	// Temp code
	//Consume all until next RCurly
	for !lex.ExpectNext(RCurly) {
		_, _ = lex.GetNext()
		//fmt.Println(t)
	}
	//Consume RCurly
	lex.GetNext()

	return head
}

// todo: remove assumptions & add error handling
// primary parsing function, can make assumptions about inital tokens
func parseBlock(lex *Lexer) AST.Block {
	block := AST.Block{}

	for !(lex.ExpectNext(EOL) || lex.ExpectNext(RCurly)) {
		var next AST.AST
		if lex.ExpectNext(Iden) {
			m := parseModule(lex)
			next = &m
		} else {
			break
		}

		block.Elements = append(block.Elements, next)
	}

	//Consume RCurly
	if lex.ExpectNext(RCurly) {
		lex.GetNext()
	}

	return block
}

func Parse(lex *Lexer) AST.AST {
	ast := parseBlock(lex)
	return &ast
}
