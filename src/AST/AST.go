package AST

import (
	"fmt"
	"strings"
)

type AST interface {
	GetNext() AST
	printAST(level int)
	String() string
}

func GetLast(ast AST) AST {
	next := ast
	for next != nil {
		ast = next
		next = ast.GetNext()
	}
	return ast
}

func PrintAST(ast AST) {
	ast.printAST(0)
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
	Block
	Name   string
	Params []Parameter
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
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(m)

	for _, elem := range m.Elements {
		elem.printAST(level + 1)
	}
}

//== Block ==
// todo: convert to interface
type Block struct {
	idx      int
	Elements []AST
}

func (blk *Block) GetNext() AST {
	if blk.idx < len(blk.Elements) {
		next := blk.Elements[blk.idx]
		blk.idx++
		return next
	}
	blk.idx = 0
	return nil
}

func (blk Block) String() string {
	return "Block"
}

func (blk *Block) printAST(level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(blk)

	for {
		next := blk.GetNext()
		if next == nil {
			break
		}
		next.printAST(level + 1)
	}
}
