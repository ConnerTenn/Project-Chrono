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

	// TODO : Add BadExpr
	displayAndCheckError("Did not find LParen to open module parameters", t, L.LParen)

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

// TODO: Split out a parseSignal function
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

		displayAndCheckError("Bit width specifier not found", t, L.Literal)

		curParam.Width, _ = strconv.Atoi(t.Value)

		t, _ = lex.GetNext()

		displayAndCheckError("Bit width closing brace not found", t, L.RBrace)
	}

	// get / set name
	t, _ = lex.GetNext()

	curParam.Name = parseIdent(t)

	// check if tied to a clock
	if lex.ExpectNext(L.Atmark) {
		// drop Atmark
		_, _ = lex.GetNext()

		// get clock info
		t, _ = lex.GetNext()

		displayAndCheckError("Clock Declaration Incorrect", t, L.Iden, L.Math)
		if t.Type == L.Math && t.Value != "!" {
			displayError("Clocks Can Only Be Negated", t, L.Iden, L.Math)
		} else if t.Type == L.Math {
			curParam.Clock.Neg = true

			t, _ = lex.GetNext()
		}

		displayAndCheckError("Clock Declaration Incorrect", t, L.Iden)
		curParam.Clock.Name = parseIdent(t)
	}

	return curParam
}

func parseIdent(t L.Token) AST.Ident {
	displayAndCheckError("Could not parse identifier", t, L.Iden)

	return AST.Ident{Pos: t.Pos, Name: t.Value}
}

func parseBlock(lex *L.Lexer) AST.BlockStmt {
	t, _ := lex.GetNext() //Consume LCurly
	blk := AST.BlockStmt{StartPos: t.Pos}

	t, _ = lex.PeekNext()

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

	if lex.ExpectNext(L.Iden) {
		//FIXME: Assume expression is an assignment
		lhs, _ := lex.GetNext()

		asmt, _ := lex.GetNext()
		displayAndCheckError("Expected assignment statement", asmt, L.Asmt)

		op := AST.Asmt
		if asmt.Value == "<-" {
			op = AST.AsmtReg
		}

		rhs := ParseExpression(lex)

		return &AST.AssignStmt{Pos: lhs.Pos, Op: op, LHS: &AST.Ident{Pos: lhs.Pos, Name: lhs.Value}, RHS: rhs}

	} else if lex.ExpectNext(L.If) {
		//Consume if
		ifToken, _ := lex.GetNext()

		cond := ParseExpression(lex)

		body := parseBlock(lex)

		return &AST.IfStmt{
			Pos:  ifToken.Pos,
			Cond: cond,
			Body: &body,
			Else: nil,
		}
	}

	next, _ := lex.GetNext()
	displayError("Could not parse statement", next, L.Iden, L.If, L.RCurly)
	return &AST.BadStmt{Pos: next.Pos}
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
		return &AST.MathExpr{Pos: head.Pos, LHS: lhs, RHS: rhs, Op: op}
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

func ParseExpression(lex *L.Lexer) AST.Expr {
	//Reverse polish notation buffer
	var rpn []L.Token
	//Stack for storing the operators
	var opStack []L.Token

	expectNext := true
	t, _ := lex.GetNext()
	for expectNext {
		expectNext = false

		switch t.Type {
		case L.Iden:
			rpn = append(rpn, t)

		case L.Literal:
			rpn = append(rpn, t)

		case L.Asmt:
			//Assignments always go directly onto the opStack (Since it is the lowest precedence)
			opStack = append(opStack, t)
			expectNext = true

		case L.LParen:
			//LParen always goes directly onto stack as a marker for when RParen is found
			opStack = append(opStack, t)
			expectNext = true

		case L.RParen:
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

		case L.Math, L.Cmp:
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
			expectNext = true

		default:
			displayError("Unknown token", t, L.Iden, L.Literal, L.Asmt, L.LParen, L.RParen, L.Math, L.Cmp)
		}

		if !expectNext {
			//Check the next token to see if a new operator exists to continue the expression
			n, _ := lex.PeekNext()
			if n.Type == L.Math || n.Type == L.Cmp || n.Type == L.RParen {
				expectNext = true
			}
		}

		//Collect the next token if required
		if expectNext {
			t, _ = lex.GetNext()
		}
	}

	for i := len(opStack) - 1; i >= 0; i-- {
		rpn = append(rpn, opStack[i])
	}

	// fmt.Println(rpn)

	//Convert rpn into a forward polish notation as a channel acting as a fifo
	fpn := make(chan L.Token, len(rpn))
	for i := len(rpn) - 1; i >= 0; i-- {
		fpn <- rpn[i]
	}

	return createExpression(fpn)
}

func isOperation(t L.Token) bool {
	if t.Type != L.Math && t.Type != L.Cmp && t.Type != L.LParen && t.Type != L.RParen {
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
	case "==":
		op = AST.Equals
	}

	return op
}

func OpCmp(op1 AST.Operation, op2 AST.Operation) int {
	return AST.Precedence[op1] - AST.Precedence[op2]
}
