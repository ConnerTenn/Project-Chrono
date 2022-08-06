package Parser

import (
	"fmt"
	"os"

	L "github.com/ConnerTenn/Project-Chrono/Lexer"
)

func displayTokenError(t L.Token) {

	fmt.Println("Unexpected Token: ", t)

	os.Exit(-1)
}

func displayError(context string, recievedToken L.Token, expected ...L.TokenType) {

	fmt.Println("Error Parsing:", context,
		"\nRecieved:", recievedToken)

	fmt.Print("Expected: ")
	for _, expect := range expected {
		fmt.Print(expect, " ")
	}
	fmt.Println()

	os.Exit(-1)
}

func displayAndCheckError(context string, recievedToken L.Token, expected ...L.TokenType) {
	for _, token := range expected {
		if recievedToken.Type == token {
			return
		}
	}

	displayError(context, recievedToken, expected...)
}
