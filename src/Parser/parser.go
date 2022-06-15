package Parser

import (
  "strconv"
  L "github.com/ConnerTenn/Project-Chrono/Lexer"
  AST "github.com/ConnerTenn/Project-Chrono/AST"
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
