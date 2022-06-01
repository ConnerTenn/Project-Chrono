package Parser

import (
	"fmt"
	"os"

	L "github.com/ConnerTenn/Project-Chrono/Lexer"
)

func displayError(t L.Token) {

	fmt.Println("Unexpected Token: ", t)

	os.Exit(-1)
}
