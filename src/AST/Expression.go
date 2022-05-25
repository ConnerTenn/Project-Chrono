package AST

import (
	_ "fmt"
	_ "strings"
)

//== Generalized expressions

type Expression interface {
	IsComputable() bool // is CTE available?
}

type AssignmentExpression struct {
	Name string
	RHS  *AST
	Next *AST
}

func (ae AssignmentExpression) String() string {
	return ""
}

func (ae AssignmentExpression) GetNext() *AST {
	return (ae.Next)
}

func (ae AssignmentExpression) printAST(level int) {
	return
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
)

type MathExpression struct {
	LHS *MathExpression
	RHS *MathExpression
	Op  Operation
}

func (m MathExpression) String() string {
	return ""
}

func (m MathExpression) GetLeft() *MathExpression {
	return nil
}

func (m MathExpression) GetRight() *MathExpression {
	return nil
}

func (m MathExpression) IsComputable() bool {
	return false
}
