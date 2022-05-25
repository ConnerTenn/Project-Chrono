package AST

import (
	"fmt"
	"strings"
)

//== Generalized expressions

type Expression interface {
	IsComputable() bool // is CTE available?
  GetNext() AST
  printAST(level int)
  String() string
}

type ValueExpression struct {
	Var   bool // is Variable?
	Value string
}

func (v ValueExpression) String() string {
	return v.Value
}

func (v ValueExpression) printAST(level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(v)
}

func (v ValueExpression) GetNext() AST {
	return nil
}

func (v ValueExpression) IsComputable() bool {
	return !v.Var
}

type AssignmentExpression struct {
	Name string
	RHS  Expression
	Next AST
}

func (ae AssignmentExpression) String() string {
	s := fmt.Sprintf("%s: %s", ae.Name, ae.RHS)
	return s
}

func (ae AssignmentExpression) GetNext() AST {
	return (ae.Next)
}

func (ae AssignmentExpression) printAST(level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(ae)
}

func (ae AssignmentExpression) IsComputable() bool {
	return false
}

//== Math ==

//go:generate stringer -type=Operation
type Operation int

const (
	Add Operation = iota
	Sub
	Multi
	Div
	LShift
	RShift
)

type MathExpression struct {
	LHS Expression
	RHS Expression
	Op  Operation
}

func (m MathExpression) String() string {
	return fmt.Sprintf("%s %s %s", m.LHS, m.Op, m.RHS)
}

func (m MathExpression) printAST(level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(m)
}

func (m MathExpression) GetNext() AST {
	return nil
}

func (m MathExpression) GetLeft() Expression {
	return m.LHS
}

func (m MathExpression) GetRight() Expression {
	return m.RHS
}

func (m MathExpression) IsComputable() bool {
	return m.LHS.IsComputable() && m.RHS.IsComputable()
}
