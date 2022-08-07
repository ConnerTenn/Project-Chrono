package AST

import (
	"fmt"
	"strings"
)

// Using Golangs breakdown of Expression, Statement, and Declaration

// base interface used by all blocks
type AST interface {
	GetPos() [2]int
}

/* --- Expressions --- */

// Specifies computation
type (
	Expr interface {
		AST
		IsComputable() bool // if the expression is evaluatable at compile time
		exprNode()
		String() string
	}

	// An Expression with Syntax errors
	BadExpr struct {
		Pos [2]int
	}

	// Represents an identifier used in computation
	Ident struct {
		Pos  [2]int
		Name string
	}

	// Represents a 'literal' or plain text value
	Literal struct {
		Pos   [2]int
		Type  string // to be switched to enum
		Value string
	}

	// Expression contained within parens (nested)
	ParenExpr struct {
		StartPos [2]int
		EndPos   [2]int
		X        Expr // Inner Expression
	}

	// Represents a function call
	CallExpr struct {
		Pos  [2]int
		Fn   string // Function identifier
		Args []Expr //List of function arguments
	}

	// Represents a math calculation
	MathExpr struct {
		Pos [2]int
		LHS Expr
		RHS Expr
		Op  Operation
	}
)

func (x *BadExpr) IsComputable() bool   { return false }
func (x *Ident) IsComputable() bool     { return false }
func (x *Literal) IsComputable() bool   { return true }
func (x *ParenExpr) IsComputable() bool { return x.X.IsComputable() }
func (x *CallExpr) IsComputable() bool  { return false }
func (x *MathExpr) IsComputable() bool  { return false }

func (x *BadExpr) GetPos() [2]int   { return x.Pos }
func (x *Ident) GetPos() [2]int     { return x.Pos }
func (x *Literal) GetPos() [2]int   { return x.Pos }
func (x *ParenExpr) GetPos() [2]int { return x.StartPos }
func (x *CallExpr) GetPos() [2]int  { return x.Pos }
func (x *MathExpr) GetPos() [2]int  { return x.Pos }

func (*BadExpr) exprNode()   {}
func (*Ident) exprNode()     {}
func (*Literal) exprNode()   {}
func (*ParenExpr) exprNode() {}
func (*CallExpr) exprNode()  {}
func (*MathExpr) exprNode()  {}

func (s *BadExpr) String() string { return "BAD EXPRESSION" }
func (x Ident) String() string {
	return x.Name
}

func (x Literal) String() string {
	return x.Value
}

func (x MathExpr) String() string {
	var str string

	str += "(" + x.LHS.String() + " "
	str += x.Op.String() + " "
	str += x.RHS.String() + ")"

	return str
}

//go:generate stringer -type=Operation
type Operation int

const (
	Asmt Operation = iota
	AsmtReg
	LShift
	RShift
	Add
	Sub
	Multi
	Div
	Bracket
	Equals
)

var Precedence = map[Operation]int{
	Asmt:    0,
	Equals:  1,
	LShift:  2,
	RShift:  2,
	Add:     3,
	Sub:     3,
	Multi:   4,
	Div:     4,
	Bracket: 5,
}

/* --- Statements --- */

// Controls execution
type (
	Stmt interface {
		AST
		stmtNode()
		String(indent int) string
	}

	// Represents a statement with incorrect syntax
	BadStmt struct {
		Pos [2]int
	}

	// A declaration linked to statements
	DeclStmt struct {
		Pos  [2]int
		Decl Decl
	}

	// Holds an expression to be executed
	ExprStmt struct {
		Pos [2]int
		X   Expr
	}

	// Represents an assignment
	AssignStmt struct {
		Pos [2]int
		Op  Operation
		LHS Expr
		RHS Expr
	}

	// Represents a block thats a sequence
	SequenceStmt struct {
		StartPos [2]int
		EndPos   [2]int
		Clk      string
		Inner    Stmt
	}

	// a function return
	ReturnStmt struct {
		Pos    [2]int
		Result Expr
	}

	// A {} block of statements
	BlockStmt struct {
		StartPos [2]int
		EndPos   [2]int
		StmtList []Stmt
	}

	// an If conditional
	IfStmt struct {
		Pos  [2]int
		Cond Expr
		Body Stmt
		Else Stmt
	}

	// a looped block
	LoopStmt struct {
		Pos  [2]int
		Cond Expr
		Body Stmt
	}
)

func (s *BadStmt) GetPos() [2]int      { return s.Pos }
func (s *DeclStmt) GetPos() [2]int     { return s.Pos }
func (s *ExprStmt) GetPos() [2]int     { return s.Pos }
func (s *AssignStmt) GetPos() [2]int   { return s.Pos }
func (s *SequenceStmt) GetPos() [2]int { return s.StartPos }
func (s *ReturnStmt) GetPos() [2]int   { return s.Pos }
func (s *BlockStmt) GetPos() [2]int    { return s.StartPos }
func (s *IfStmt) GetPos() [2]int       { return s.Pos }
func (s *LoopStmt) GetPos() [2]int     { return s.Pos }

func (*BadStmt) stmtNode()      {}
func (*DeclStmt) stmtNode()     {}
func (*ExprStmt) stmtNode()     {}
func (*AssignStmt) stmtNode()   {}
func (*SequenceStmt) stmtNode() {}
func (*ReturnStmt) stmtNode()   {}
func (*BlockStmt) stmtNode()    {}
func (*IfStmt) stmtNode()       {}
func (*LoopStmt) stmtNode()     {}

func Indent(level int) string {
	return strings.Repeat("  ", level)
}

func (s *BadStmt) String(indent int) string { return "BAD STATEMENT" }
func (s AssignStmt) String(indent int) string {
	var str string

	str += Indent(indent)
	str += s.LHS.String()
	if s.Op == AsmtReg {
		str += " <- "
	} else {
		str += " = "
	}
	str += s.RHS.String()

	return str
}
func (s *IfStmt) String(indent int) string {
	var str string
	str += Indent(indent)
	str += "if " + s.Cond.String() + "\n"

	str += s.Body.String(indent + 1)

	if s.Else != nil {
		str += Indent(indent) + "else\n" + s.Else.String(indent+1)
	}

	return str
}

func (s BlockStmt) String(indent int) string {
	var str string

	for _, stmt := range s.StmtList {
		str += stmt.String(indent) + "\n"
	}

	return str
}

/* --- Declarations --- */

// Binding of identifiers
type (
	Decl interface {
		AST
		declNode()
	}

	ValueDecl struct {
		Name  Ident
		Type  string // FIXME
		Value Expr
	}

	ClockDecl struct {
		Name Ident
		Neg  bool
	}

	SignalDecl struct {
		Name  Ident
		Width int
		Clock *ClockDecl
	}

	ParamDecl struct { //Extends SignalDecl
		SignalDecl
		Dir ParamDir
	}

	ModuleDecl struct {
		Name   Ident
		Params []ParamDecl
		Block  BlockStmt
	}
)

//go:generate stringer -type=ParamDir
type ParamDir int

const (
	In ParamDir = iota
	Out
	Inout
)

func (d ValueDecl) GetPos() [2]int  { return d.Name.GetPos() }
func (d ModuleDecl) GetPos() [2]int { return d.Name.GetPos() }
func (d ClockDecl) GetPos() [2]int  { return d.Name.GetPos() }

func (*ValueDecl) declNode()  {}
func (*ParamDecl) declNode()  {}
func (*ModuleDecl) declNode() {}
func (*ClockDecl) declNode()  {}

func (d ClockDecl) String() string {
	var str string
	if d.Name.Name != "" {
		str += " @ "

		if d.Neg {
			str += "!"
		}

		str += d.Name.Name
	}

	return str
}

func (d SignalDecl) String() string {
	var str string
	if d.Width > 1 {
		str += "[" + fmt.Sprint(d.Width) + "] "
	}
	str += d.Name.Name

	if d.Clock != nil {
		str += d.Clock.String()
	}

	return str
}

func (d ParamDecl) String() string {
	var str string
	str += d.Dir.String() + " "
	str += d.SignalDecl.String()
	return str
}

func (d ModuleDecl) String() string {
	var str string
	str += d.Name.Name
	str += "("
	for i, param := range d.Params {
		str += param.String()
		if i < len(d.Params)-1 {
			str += ", "
		}
	}
	str += ")\n"

	str += d.Block.String(1)

	return str
}
