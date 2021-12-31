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
	Signal            // input or output, probably needs a better name here
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
	Tokens chan Token // might need to be a locked array to enable peeking
}

func (lex Lexer) StartLexing(fileName string) error {
  var err error // throws non-name on left side error (for lex.file) when using := operator
	lex.file, err = os.Open(fileName)
	if err != nil {
		return err
	}

  //lex.Tokens = make(chan Token, 20)

  go lex.tokenizer()

	return nil
}

func (lex Lexer) tokenizer() {
  scanner := bufio.NewScanner(lex.file)
  scanner.Split(bufio.ScanRunes)

  pos := [2]int{1,1} // position / head tracker for error reporting

  for scanner.Scan() {
    var t TokenType
    val := scanner.Text()
    // currently only deal with single char tokens

    // single char tokenizing
    switch val {
      case " ": {pos[1]+=1; continue} // ignore spaces
      case "\n": {pos[0]+=1; pos[1]=1; continue} // ignore new lines
      case "\r": continue // ignore carriage returns (don't need to count for position tracking, always bundled with new line)
      case ",": t = Comma
      case "{": t = LCurly
      case "}": t = RCurly
      case "(": t = LParen
      case ")": t = RParen
      case ";": t = EOL
    }

    lex.Tokens <- Token{t, val, pos}
    pos[1]+=1;
  }

  close(lex.Tokens)
}
