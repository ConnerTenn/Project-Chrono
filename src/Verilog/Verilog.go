package verilog

import (
	"fmt"
	"os"
	"reflect"
	"runtime"

	AST "github.com/ConnerTenn/Project-Chrono/AST"
)

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

func EmitModuleDecl(obj AST.ModuleDecl) {
	writeToFile("\nmodule " + obj.Name.Name + " (\n")

	//Write parameters
	for i, param := range obj.Params {
		str := "\t" + param.Dir.String()
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
		if i < len(obj.Params)-1 {
			str += ","
		}
		str += "\n"
		writeToFile(str)
	}
	writeToFile(");\n")

	writeToFile("endmodule\n")
}

func GenerateVerilog(ast []AST.AST) {
	createFile("generated.sv")
	defer closeFile()

	for _, elem := range ast {
		switch obj := elem.(type) {
		case AST.ModuleDecl:
			fmt.Println("BlockStmt")
			EmitModuleDecl(obj)
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
