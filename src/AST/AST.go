package AST

import (
	"fmt"
	"strings"
)

var IndentCharacters = "  "

func Ident(ident int) string {
	return strings.Repeat(IndentCharacters, ident)
}

type AST interface {
	printAST(level int)
	String() string
	WriteVerilog(ident int) string
}

func PrintAST(ast AST) {
	ast.printAST(0)
}

//============
//== Signal ==
//============

//go:generate stringer -type=SignalDir
type SignalDir int

const (
	In SignalDir = iota
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

type Signal struct {
	Name  string
	Dir   SignalDir
	Width int
	Type  SignalType
}

func (p Signal) String() string {
	return fmt.Sprintf("%s:%s,%s[%d]", p.Name, p.Dir, p.Type, p.Width)
}

//============
//== Module ==
//============

type Module struct {
	Block
	Name   string
	Params []Signal
}

// todo: add child and next
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

func (m *Module) printAST(level int) {
	fmt.Print(Ident(level))
	fmt.Println(m)

	for _, elem := range m.Elements {
		elem.printAST(level + 1)
	}
}

//===========
//== Block ==
//===========

// todo: convert to interface
type Block struct {
	Elements []AST
}

func (blk Block) String() string {
	return "Block"
}

func (blk *Block) printAST(level int) {
	fmt.Print(Ident(level))
	fmt.Println(blk)

	for _, elem := range blk.Elements {
		elem.printAST(level + 1)
	}
}
