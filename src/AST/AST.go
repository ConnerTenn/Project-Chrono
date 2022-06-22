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
	MathStmt struct {
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
func (x *MathStmt) IsComputable() bool  { return false }

func (x *BadExpr) GetPos() [2]int   { return x.Pos }
func (x *Ident) GetPos() [2]int     { return x.Pos }
func (x *Literal) GetPos() [2]int   { return x.Pos }
func (x *ParenExpr) GetPos() [2]int { return x.StartPos }
func (x *CallExpr) GetPos() [2]int  { return x.Pos }
func (x *MathStmt) GetPos() [2]int  { return x.Pos }

func (*BadExpr) exprNode()   {}
func (*Ident) exprNode()     {}
func (*Literal) exprNode()   {}
func (*ParenExpr) exprNode() {}
func (*CallExpr) exprNode()  {}
func (*MathStmt) exprNode()  {}

func (s *BadExpr) String() string { return "BAD EXPRESSION" }
func (x Ident) String() string {
	return x.Name
}

func (x Literal) String() string {
	return x.Value
}

func (x MathStmt) String() string {
	var str string

	str += x.LHS.String() + " "
	str += x.Op.String() + " "
	str += x.RHS.String()

	return str
}

//go:generate stringer -type=Operation
type Operation int

const (
	Asmt Operation = iota
	Add
	Sub
	Multi
	Div
	LShift
	RShift
	Bracket
)

/* --- Statements --- */

// Controls execution
type (
	Stmt interface {
		AST
		stmtNode()
		String() string
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

func (s *BadStmt) String() string { return "BAD STATEMENT" }
func (s AssignStmt) String() string {
	var str string

	str += s.LHS.String()
	str += " = "
	str += s.RHS.String()

	return str
}

func (s BlockStmt) String() string { return s.StringIdent(0) }
func (s BlockStmt) StringIdent(level int) string {
	var str string

	for _, stmt := range s.StmtList {
		str += strings.Repeat("  ", level) + stmt.String() + "\n"
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

	SignalDecl struct {
		Name  Ident
		Width int
		Type  SignalType
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

//go:generate stringer -type=SignalType
type SignalType int

const (
	Wire SignalType = iota //Default
	Reg
	//Var?
)

func (d ValueDecl) GetPos() [2]int  { return d.Name.GetPos() }
func (d ModuleDecl) GetPos() [2]int { return d.Name.GetPos() }

func (*ValueDecl) declNode()  {}
func (*ParamDecl) declNode()  {}
func (*ModuleDecl) declNode() {}

func (d SignalDecl) String() string {
	var str string
	str += d.Type.String() + " "
	if d.Width > 1 {
		str += "[" + fmt.Sprint(d.Width) + "] "
	}
	str += d.Name.Name

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

	str += d.Block.StringIdent(1)

	return str
}
