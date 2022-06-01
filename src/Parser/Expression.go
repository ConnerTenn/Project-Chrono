package Parser

import (
	AST "github.com/ConnerTenn/Project-Chrono/AST"
	L "github.com/ConnerTenn/Project-Chrono/Lexer"
)

func parseOperation(lex *L.Lexer) AST.Operation {
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

func parseValue(lex *L.Lexer) AST.ValueExpression {
	value := AST.ValueExpression{}
	t, _ := lex.GetNext()

	value.Value = t.Value

	if t.Type == L.Iden {
		value.Var = true
	} else if t.Type == L.Literal {
		value.Var = false
	} else {
		displayError(t)
	}

	return value
}

func parseMathExpression(lex *L.Lexer) AST.MathExpression {
	math := AST.MathExpression{}

	// assume lhs and rhs are values and proper syntax is given
	math.LHS = parseValue(lex)

	math.Op = parseOperation(lex)

	math.RHS = parseValue(lex)

	lex.GetNext() //FIXME: Consume Semicolon

	return math
}
