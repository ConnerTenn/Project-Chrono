package main

import "fmt"

// TODO: split these out in a package

type AST interface {
	GetNext() *AST
}

type ParamDir int

const (
	In ParamDir = iota
	Out
	Inout
)

type ParamType int

const (
	Reg ParamType = iota
	Wire
	//Var?
	// Which is default?
)

type Parameter struct {
	Name  string
	Dir   ParamDir
	Width int
	Type  ParamType
}

type Module struct {
	Name   string
	Params []Parameter
	Block  *AST
}

func (m Module) String() string {
	var params string
	for i, param := range m.Params {
		params += fmt.Sprintf("%d: \n  Name: %v\n  Dir: %v\n  Width: %d\n  Type: %d\n",
			i, param.Name, param.Dir, param.Width, param.Type)
	}
	return "Name: " + m.Name + "\nParams: " + params
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

func (m MathExpression) GetLeft() *MathExpression {
	return nil
}

func (m MathExpression) GetRight() *MathExpression {
	return nil
}

func (m MathExpression) IsComputable() bool {
	return false
}
