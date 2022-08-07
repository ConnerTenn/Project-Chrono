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
		curToken := lex.GetNext()
		//nextToken, _ := lex.PeekNext()

		// FIXME: assume module
		tree = append(tree, parseModule(lex, curToken))

	}

	return tree
}

func parseModule(lex *L.Lexer, t L.Token) AST.ModuleDecl {
	newModule := AST.ModuleDecl{}

	newModule.Name = parseIdent(t)
	t = lex.GetNext()

	// TODO : Add BadExpr
	displayAndCheckError("Did not find LParen to open module parameters", t, L.LParen)

	// build parameters
	for !lex.ExpectNext(")") {
		t = lex.GetNext() //get parameters

		newModule.Params = append(newModule.Params, parseParam(lex, t))

		if lex.ExpectNext(",") {
			lex.GetNext() //drop comma
		}
	}
	lex.GetNext() //drop RParen

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

	// set / get bit width
	if lex.ExpectNext("[") {
		lex.GetNext()
		t = lex.GetNext()

		displayAndCheckError("Bit width specifier not found", t, L.Literal)

		curParam.Width, _ = strconv.Atoi(t.Value)

		t = lex.GetNext()

		displayAndCheckError("Bit width closing brace not found", t, L.RBrace)
	}

	// get / set name
	t = lex.GetNext()

	curParam.Name = parseIdent(t)

	// check if tied to a clock
	if lex.ExpectNext("@") {
		// drop Atmark
		_ = lex.GetNext()

		// get clock info
		t = lex.GetNext()

		displayAndCheckError("Clock Declaration Incorrect", t, L.Iden, L.Math)
		if t.IsMath() {
			if t.Is("!") {
				curParam.Clock.Neg = true

				t = lex.GetNext()
			} else {
				displayError("Clocks Can Only Be Negated", t, L.Iden, L.Math)
			}
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
	t := lex.GetNext() //Consume LCurly
	displayAndCheckError("Block Statement improperly started", t, L.LCurly)

	blk := AST.BlockStmt{StartPos: t.Pos}

	//Run until end of block
	for t = lex.PeekNext(); !t.IsRCurly(); t = lex.PeekNext() {
		//FIXME : Assuming blocks contain only statements
		blk.StmtList = append(blk.StmtList, parseStatement(lex))
	}
	displayAndCheckError("Block Statement improperly terminated", t, L.RCurly)
	lex.GetNext() //Consume RCurly

	blk.EndPos = t.Pos
	return blk
}

//FIXME : Definitely a lot to be added here
func parseStatement(lex *L.Lexer) AST.Stmt {
	next := lex.PeekNext()
	displayAndCheckError("Invalid beginning of a statement", next, L.Iden, L.If, L.LCurly)

	if next.IsIden() {
		//FIXME: Assume expression is an assignment
		lhs := lex.GetNext()

		asmt := lex.GetNext()
		displayAndCheckError("Expected assignment statement", asmt, L.Asmt)

		op := AST.Asmt
		if asmt.Is("<-") {
			op = AST.AsmtReg
		}

		rhs := ParseExpression(lex)

		return &AST.AssignStmt{Pos: lhs.Pos, Op: op, LHS: &AST.Ident{Pos: lhs.Pos, Name: lhs.Value}, RHS: rhs}

	} else if next.Is("if") {
		//Consume if
		ifToken := lex.GetNext()

		cond := ParseExpression(lex)

		body := parseBlock(lex)

		newStmt := AST.IfStmt{
			Pos:  ifToken.Pos,
			Cond: cond,
			Body: &body,
			Else: nil,
		}

		if lex.ExpectNext("else") {
			// drop else
			_ = lex.GetNext()
			t := lex.PeekNext()

			displayAndCheckError("else cannot be an arbitrary statement", t, L.If, L.LCurly)

			newStmt.Else = parseStatement(lex)
		}

		return &newStmt
	} else if next.IsLCurly() {
		newBlock := parseBlock(lex)
		return &newBlock
	} else {
		return &AST.BadStmt{Pos: next.Pos}
	}
}

//fpn: Forward polish notation
func createExpression(fpn chan L.Token) AST.Expr {
	head := <-fpn
	if head.IsOperator() || head.IsParen() {
		op := parseOperation(head)
		//Recursively collect the RHS
		rhs := createExpression(fpn)
		//Recursively collect the LHS
		lhs := createExpression(fpn)
		return &AST.MathExpr{Pos: head.Pos, LHS: lhs, RHS: rhs, Op: op}
	}

	//If this isn't an operation, it must be an Iden or a Literal
	if head.IsIden() {
		return &AST.Ident{Pos: head.Pos, Name: head.Value}
	} else if head.IsLiteral() {
		return &AST.Literal{Pos: head.Pos, Value: head.Value}
	} else {
		return &AST.BadExpr{}
	}
}

func ParseExpression(lex *L.Lexer) AST.Expr {
	//Reverse polish notation buffer
	var rpn []L.Token
	//Stack for storing the operators
	var opStack []L.Token

	expectNext := true
	t := lex.GetNext()
	for expectNext {
		expectNext = false

		if t.IsIden() {
			rpn = append(rpn, t)

		} else if t.IsLiteral() {
			rpn = append(rpn, t)

		} else if t.IsAssignment() {
			//Assignments always go directly onto the opStack (Since it is the lowest precedence)
			opStack = append(opStack, t)
			expectNext = true

		} else if t.IsLParen() {
			//LParen always goes directly onto stack as a marker for when RParen is found
			opStack = append(opStack, t)
			expectNext = true

		} else if t.IsRParen() {
			//Place all operations till the previous LParen onto the rpn stack
			op := opStack[len(opStack)-1]
			for !op.IsLParen() {
				//Pop off opStack and place into rpn
				rpn = append(rpn, op)
				opStack = opStack[:len(opStack)-1]

				op = opStack[len(opStack)-1]
			}
			//Remove LParen from opStack
			opStack = opStack[:len(opStack)-1]

		} else if t.IsOperator() {
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

		} else {
			displayError("Unknown token", t, L.Iden, L.Literal, L.Asmt, L.LParen, L.RParen, L.Math, L.Cmp)
		}

		if !expectNext {
			//Check the next token to see if a new operator exists to continue the expression
			n := lex.PeekNext()
			if n.IsOperator() || n.IsRParen() {
				expectNext = true
			}
		}

		//Collect the next token if required
		if expectNext {
			t = lex.GetNext()
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

func parseOperation(t L.Token) AST.Operation {
	if !(t.IsOperator() || t.IsParen()) {
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
