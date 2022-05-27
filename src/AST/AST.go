package AST

import (
	"fmt"
	"strings"
)

type AST interface {
	GetNext() AST      //Get immediately following
	GetLast() AST      //Get last of all submodules
	GetAfter() AST     //Get next after all submodules
	SetNext(next AST)  //Set immediately following
	SetAfter(next AST) //Set the next after all submodules
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
	Name   string
	Params []Parameter
	Next   AST
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

func (m *Module) GetNext() AST {
	return m.Next
}

func (m *Module) GetLast() AST {
	last := (AST)(m)

	if m.GetNext() != nil {
		last = m.GetNext()
	}

	return last
}

func (m *Module) GetAfter() AST {
	return m.GetNext()
}

func (m *Module) SetNext(next AST) {
	m.Next = next
}

func (m *Module) SetAfter(next AST) {
	m.GetLast().SetNext(next)
}

func (m *Module) printAST(level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(m)
}

//== Block ==
// todo: convert to interface
type Block struct {
	NumElements int
	Next        AST
}

func (blk *Block) GetNext() AST {
	return blk.Next
}

func (blk *Block) GetLast() AST {
	last := (AST)(blk)
	if blk.GetNext() != nil {
		last = blk.GetNext()
	}
	for i := 0; i < blk.NumElements-1; i++ {
		last = last.GetLast()
	}
	return last
}

func (blk *Block) GetAfter() AST {
	return blk.GetLast().GetNext()
}

func (b *Block) SetNext(next AST) {
	b.Next = next
}

func (b *Block) SetAfter(next AST) {
	b.GetLast().SetNext(next)
	b.NumElements++
}

func (blk Block) String() string {
	return "Block [" + fmt.Sprint(blk.NumElements) + "]"
}

func (blk *Block) printAST(level int) {
	fmt.Print(strings.Repeat(" ", level*2))
	fmt.Println(blk)

	next := blk.Next
	for i := 0; i < blk.NumElements; i++ {
		next.printAST(level + 1)
		next = next.GetAfter()
	}
}
