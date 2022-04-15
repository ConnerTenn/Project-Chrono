package main

import (
	"fmt"
	"strings"
)

// TODO: split these out in a package

type AST interface {
	GetNext() []*AST
}

func printAST(ast *AST, level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(*ast)
	nextList := (*ast).GetNext()
	for _, element := range nextList {
		printAST(element, level+1)
	}
}

func PrintAST(ast *AST) {
	printAST(ast, 0)
}

//== Module ==

//go:generate stringer -type=ParamDir
type ParamDir int

const (
	In ParamDir = iota
	Out
	Inout
)

//go:generate stringer -type=ParamType
type ParamType int

const (
	Wire ParamType = iota //Default
	Reg
	//Var?
)

type Parameter struct {
	Name  string
	Dir   ParamDir
	Width int
	Type  ParamType
}

func (p Parameter) String() string {
  return fmt.Sprintf("%s:%s,%s[%d]", p.Name, p.Dir, p.Type, p.Width)
}

type Module struct {
	Name     string
	Params   []Parameter
	Children []*AST
}

func (m Module) String() string {
	var params string
	for i, param := range m.Params {
		//Add space between parameters
		if i != 0 {
			params += " "
		}
    params += param.String()
	}
	return "mod:" + m.Name + " (" + params + ")"
}

func (m Module) GetNext() []*AST {
	return m.Children
}

//== Block ==

type Block struct {
	Children []*AST
}

func (blk Block) GetNext() []*AST {
	return blk.Children
}

func (blk Block) String() string {
	return "Block"
}

//== Math ==

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
