package verilog

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	AST "github.com/ConnerTenn/Project-Chrono/AST"
)

func Indent(level int) string {
	return strings.Repeat("\t", level)
}

func displayError(msg string) {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("Error: "+msg+"\n -- at  %s:%d\n", fn, line)
	os.Exit(-1)
}

var outfile *os.File

func createFile(filename string) {
	outfile, _ = os.Create(filename)
}
func writeToFile(str string) {
	outfile.WriteString(str)
}
func closeFile() {
	outfile.Close()
}

func emitModuleDecl(mod AST.ModuleDecl) {
	writeToFile("\nmodule " + mod.Name.Name + " (\n")

	//Write parameters
	for i, param := range mod.Params {
		str := Indent(1)
		//Dir
		switch param.Dir {
		case AST.In:
			str += "input"
		case AST.Out:
			str += "output"
		case AST.Inout:
			str += "inout"
		}
		//wire/reg
		if param.Clock != nil {
			str += " reg"
		}
		//Bus width
		if param.Width > 1 {
			str += fmt.Sprintf(" [%d:0]", param.Width-1)
		}
		//Name
		str += " " + param.Name.Name
		if i < len(mod.Params)-1 {
			str += ","
		}
		str += "\n"
		writeToFile(str)
	}
	writeToFile(");\n")

	emitBlock(mod.Block, 0, false)

	writeToFile("endmodule\n")
}

func emitBlock(blk AST.BlockStmt, ident int, surround bool) {
	if surround {
		writeToFile(Indent(ident) + "begin\n")
	}

	//Write each statement
	str := ""
	for _, stmt := range blk.StmtList {
		emitStatement(stmt, ident+1)
	}
	writeToFile(str)

	if surround {
		writeToFile(Indent(ident) + "end\n")
	}
}

func emitStatement(stmt AST.Stmt, ident int) {
}

func GenerateVerilog(ast []AST.AST) {
	createFile("generated.sv")
	defer closeFile()

	for _, elem := range ast {
		switch obj := elem.(type) {
		case AST.ModuleDecl:
			fmt.Println("BlockStmt")
			emitModuleDecl(obj)
		default:
			displayError("Unexpected AST element: " + fmt.Sprint(reflect.TypeOf(elem)))
		}
	}
	// 	filescope := ast.(*AST.BlockStmt)

	// for stmt := range filescope.StmtList
	// 	switch obj := filescope.(type) {
	// 	case AST.Expr:
	// 		fmt.Println("Expr", obj.(AST.Expr))
	// 	case AST.Stmt:
	// 		fmt.Println("Stmt")
	// 	case *AST.BlockStmt:
	// 		fmt.Println("BlockStmt")
	// 	case AST.Decl:
	// 		fmt.Println("Decl")
	// 	}
}
