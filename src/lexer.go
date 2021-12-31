package main

import (
	"bufio"
  "fmt"
	"os"
)

// token types split into enum per type for easier parsing (on switch)
type TokenType int

const (
	EOL     TokenType = iota
	Iden              // any named identifier
	Literal           // can be split into the different literal types during lexing
	Signal            // input vs output, probably needs a better name here
	LParen
	RParen
	LCurly
	RCurly
  Comma
)

type Token struct {
	Type  TokenType
	Value string // store the direct value
	Pos   [2]int // store the position for error reporting
}

// stringer interface
func (t Token) String() string {
  return fmt.Sprintf("%v: %s at %d:%d", t.Type, t.Value, t.Pos[0], t.Pos[1])
}

// could be done functionally :)
// lexer consumes file line by line and sends tokens to a channel to be consumed by the parser
type Lexer struct {
	file   *os.File
	tokens chan Token // might need to be a locked array to enable peeking
}

func (lex Lexer) StartLexing(fileName string) error {
  var err error // throws non-name on left side error when using := operator
	lex.file, err = os.Open(fileName)
	if err != nil {
		return err
	}

  go lex.tokenizer()

	return nil
}

func (lex Lexer) tokenizer() {
  fileReader := bufio.NewReader(lex.File)


}
