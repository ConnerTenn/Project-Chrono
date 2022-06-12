package main

import (
	"fmt"
	"os"

	L "github.com/ConnerTenn/Project-Chrono/Lexer"
)

// todo: add type to hold CLI options, with description for help menu
func ShowHelp() {
	fmt.Print(`Usage:
./Project-Chrono [-h] [<file>]
    -h --help       Show the help menu.
                    This argument is optional and will cause the program to
                    exit immediately.
`)

	os.Exit(-1)
}

func main() {
	// parse CLI command
	args := os.Args

	if len(args) == 1 {
		fmt.Println("Please specify a file for compiling.")
		ShowHelp()
	}

	if args[1] == "-h" || args[1] == "--help" {
		ShowHelp()
	}

	filename := args[1]

	// for loop / switch over compiler options?

	lex, err := L.NewLexer(filename)
	if err != nil {
		fmt.Println("Error:", err)
	}

	go lex.Tokenizer()

	// doing this sync for now
	//tree := P.Parse(&lex)

	//AST.PrintAST(tree)
}
