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

	return AST.Ident{Pos: t.Pos, Name: t.Value}
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

//fpn: Forward polish notation
func createExpression(fpn chan L.Token) AST.Expr {
	head := <-fpn
	if isOperation(head) {
		op := parseOperation(head)
		//Recursively collect the RHS
		rhs := createExpression(fpn)
		//Recursively collect the LHS
		lhs := createExpression(fpn)
		return &AST.MathStmt{Pos: head.Pos, LHS: lhs, RHS: rhs, Op: op}
	}

	//If this isn't an operation, it must be an Iden or a Literal
	switch head.Type {
	case L.Iden:
		return &AST.Ident{Pos: head.Pos, Name: head.Value}
	case L.Literal:
		return &AST.Literal{Pos: head.Pos, Type: head.Type.String(), Value: head.Value}
	default:
		return &AST.BadExpr{}
	}
}

//fpn: Forward polish notation
func createStatement(fpn chan L.Token) AST.Stmt {
	head := <-fpn

	op := parseOperation(head)
	switch op {
	case AST.Asmt:
		//Recursively collect the RHS
		rhs := createExpression(fpn)
		//Recursively collect the LHS
		lhs := createExpression(fpn)
		return &AST.AssignStmt{Pos: head.Pos, LHS: lhs, RHS: rhs}
	}
	return &AST.BadStmt{}
}

//FIXME : Definitely a lot to be added here
func parseStatement(lex *L.Lexer) AST.Stmt {
	//Reverse polish notation buffer
	var rpn []L.Token
	//Stack for storing the operators
	var opStack []L.Token

	t, _ := lex.GetNext()
	for t.Type != L.EOL {

		if t.Type == L.Iden {
			rpn = append(rpn, t)

		} else if t.Type == L.Literal {
			rpn = append(rpn, t)

		} else if t.Type == L.Asmt {
			//Assignments always go directly onto the opStack (Since it is the lowest precedence)
			opStack = append(opStack, t)

		} else if t.Type == L.LParen {
			//LParen always goes directly onto stack as a marker for when RParen is found
			opStack = append(opStack, t)

		} else if t.Type == L.RParen {
			//Place all operations till the previous LParen onto the rpn stack
			op := opStack[len(opStack)-1]
			for op.Type != L.LParen {
				//Pop off opStack and place into rpn
				rpn = append(rpn, op)
				opStack = opStack[:len(opStack)-1]

				op = opStack[len(opStack)-1]
			}
			//Remove LParen from opStack
			opStack = opStack[:len(opStack)-1]

		} else if t.Type == L.Math {
			if len(opStack) > 0 {
				op1 := parseOperation(opStack[len(opStack)-1])
				op2 := parseOperation(t)
				//If op on the opStack is higher (or equal) precedence, place into the rpn immediately
				if op1 != AST.Bracket && OpCmp(op1, op2) >= 0 {
					rpn = append(rpn, opStack[len(opStack)-1])
					opStack = opStack[:len(opStack)-1]
				}
			}
			opStack = append(opStack, t)

		} else {
			displayError("Unknown token", t, L.EOL)
		}

		t, _ = lex.GetNext()
	}

	for i := len(opStack) - 1; i >= 0; i-- {
		rpn = append(rpn, opStack[i])
	}

	// fmt.Println(rpn)

	//FIXME: Assume expression is an assignment
	if rpn[len(rpn)-1].Type != L.Asmt {
		displayError("Expected assignment expression", rpn[len(rpn)-1], L.Asmt)
	}

	//Convert rpn into a forward polish notation as a channel acting as a fifo
	fpn := make(chan L.Token, len(rpn))
	for i := len(rpn) - 1; i >= 0; i-- {
		fpn <- rpn[i]
	}

	return createStatement(fpn)
}

func isOperation(t L.Token) bool {
	if t.Type != L.Math && t.Type != L.Asmt && t.Type != L.LParen && t.Type != L.RParen {
		return false
	}
	return true
}

func parseOperation(t L.Token) AST.Operation {
	if !isOperation(t) {
		displayError("Expected operator in equation", t, L.Math)
	}

	var op AST.Operation

	switch t.Value {
	case "=":
		op = AST.Asmt
	case "<=":
		op = AST.Asmt
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
	case "(":
		op = AST.Bracket
	case ")":
		op = AST.Bracket
	}

	return op
}

func OpCmp(op1 AST.Operation, op2 AST.Operation) int {
	return AST.Precedence[op1] - AST.Precedence[op2]
}
