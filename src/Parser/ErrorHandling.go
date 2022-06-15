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

func displayError(context string, recievedToken L.Token, expected L.TokenType) {

	fmt.Println("Error Parsing:", context,
		"\nRecieved:", recievedToken, "\nExpected:", expected)

	os.Exit(-1)
}
