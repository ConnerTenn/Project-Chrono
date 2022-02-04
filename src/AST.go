package main

// TODO: split these out in a package

type AST interface {
	GetNext() *AST
}

type ParamType int

const (
	In = iota
	Out
	Inout
)

type Parameters struct {
	Name  string
	Type  ParamType
	Width int
	Reg   bool
}

type Module struct {
	Name   string
	Params []Parameters
	Block  *AST
}

func (m Module) GetNext() *AST {
	return m.Block
}

type Block struct {
}

type Operation int

const (
	Add = iota
	Sub
	Multi
	Div
)

type MathExpression struct {
	LHS *MathExpression
	RHS *MathExpression
	Op  Operation
}

func GetLeft() *MathExpression {
	return nil
}

func GetRight() *MathExpression {
	return nil
}

func IsComputable() bool {
	return false
}
