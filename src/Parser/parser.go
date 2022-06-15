package Parser

import (
	"strconv"

	AST "github.com/ConnerTenn/Project-Chrono/AST"
	L "github.com/ConnerTenn/Project-Chrono/Lexer"
)

// entry parsing function
func Parse(lex *L.Lexer) []AST.AST {
	var tree []AST.AST

	for lex.NextExists() {
		// dispatches to top level declarations
		curToken, _ := lex.GetNext()
		//nextToken, _ := lex.PeekNext()

		// FIXME: assume module
		tree = append(tree, parseModule(lex, curToken))

	}

	return tree
}

func parseModule(lex *L.Lexer, t L.Token) AST.ModuleDecl {
	newModule := AST.ModuleDecl{}

	newModule.Name = parseIdent(t)
	t, _ = lex.GetNext()
	if t.Type != L.LParen {
		// TODO : Add BadExpr?
		displayError("Did not find LParen to open module parameters", t, L.LParen)
	}

	// build parameters
	t, _ = lex.GetNext() // drop LParen
	for t.Type != L.RParen {
		newModule.Params = append(newModule.Params, parseParam(lex, t))

		t, _ = lex.GetNext() // get next
		if t.Type == L.Comma {
			t, _ = lex.GetNext() // drop comma
		}
	}

	//FIXME: Bypass block

	newModule.Block = parseBlock(lex)

	return newModule
}

func parseParam(lex *L.Lexer, t L.Token) AST.ParamDecl {
	curParam := AST.ParamDecl{}

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
			displayError("Bit width specifier not found", t, L.Literal)
		}

		curParam.Width, _ = strconv.Atoi(t.Value)

		t, _ = lex.GetNext()

		if t.Type != L.RBrace {
			displayError("Bit width closing brace not found", t, L.RBrace)
		}
	}

	// get / set name
	t, _ = lex.GetNext()

	curParam.Name = parseIdent(t)

	return curParam
}

func parseIdent(t L.Token) AST.Ident {
	if t.Type != L.Iden {
		displayError("Could not parse identifier", t, L.Iden)
	}

	return AST.Ident{t.Pos, t.Value}
}

func parseBlock(lex *L.Lexer) AST.BlockStmt {
	t, _ := lex.GetNext() //Consume LCurly
	blk := AST.BlockStmt{StartPos: t.Pos}

	//Run until end of block
	for t.Type != L.RCurly {
		//FIXME : Assuming blocks contain only statements
		blk.StmtList = append(blk.StmtList, parseStatement(lex))

		t, _ = lex.PeekNext()
	}
	lex.GetNext() //Consume RCurly

	blk.EndPos = t.Pos
	return blk
}

//FIXME : Definitely a lot to be added here
func parseStatement(lex *L.Lexer) AST.Stmt {
	var assign AST.AssignStmt

	t, _ := lex.GetNext()
	if t.Type != L.Iden {
		displayError("Expected identifier at the start of a statement", t, L.Iden)
	}
	assign.LHS = &AST.Ident{Name: t.Value, Pos: t.Pos}

	t, _ = lex.GetNext()
	if t.Type != L.Asmt {
		displayError("Expected '=' to follow an identifier", t, L.Iden)
	}

	var equation AST.MathStmt
	t, _ = lex.GetNext()
	if t.Type == L.Iden {
		equation.LHS = &AST.Ident{Name: t.Value, Pos: t.Pos}
	} else if t.Type == L.Literal {
		equation.LHS = &AST.Literal{Value: t.Value, Pos: t.Pos}
	} else {
		displayError("Expected identifier/literal in equation", t, L.Iden)
	}

	equation.Op = parseOperation(lex)

	t, _ = lex.GetNext()
	if t.Type == L.Iden {
		equation.RHS = &AST.Ident{Name: t.Value, Pos: t.Pos}
	} else if t.Type == L.Literal {
		equation.RHS = &AST.Literal{Value: t.Value, Pos: t.Pos}
	} else {
		displayError("Expected identifier/literal in equation", t, L.Iden)
	}

	assign.RHS = &equation

	t, _ = lex.GetNext()
	if t.Type != L.EOL {
		displayError("Expected ';' at the end of a statement", t, L.Iden)
	}

	return &assign
}

func parseOperation(lex *L.Lexer) AST.Operation {
	t, _ := lex.GetNext()
	if t.Type != L.Math {
		displayError("Expected operator in equation", t, L.Math)
	}

	var op AST.Operation

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
