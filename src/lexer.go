package main

import (
	"bufio"
	"fmt"
	"os"
)

// token types split into enum per type for easier parsing (on switch)
type TokenType int

const (
	EOL       TokenType = iota
	Iden                // any named identifier
	Literal             // can be split into the different literal types during lexing
	Direction           // signal direction, input or output
	Keyword             // generic keyword, to be swapped out with more specific groupings (such as 'direction')
	LParen
	RParen
	LCurly
	RCurly
	Comma
	Asmt
	Cmp
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

// keyword list
var keywordMap = map[string]TokenType{
	"in":  Direction,
	"out": Direction,
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

func MultiToken(first rune, next rune) bool {
	// check if valid begining using Rune / Unicode values
	checkVal := func(val rune) bool {
		return (val >= '0' && val <= '9') || // test if number
			(val >= 'A' && val <= 'Z') || // test if uppercase character
			(val >= 'a' && val <= 'z') || // test if lowercase character
			val == '_' // test if underscore
	}

	checkCmp := func(val rune) bool {
		return (val == '<' || val == '=' || val == '>')
	}

	if checkVal(first) && checkVal(next) {
		return true
	}

	if checkCmp(first) && checkCmp(next) {
		return true
	}

	return false
}

func (lex Lexer) tokenizer() {
	reader := bufio.NewReader(lex.file)

	defer lex.file.Close()
	defer close(lex.Tokens)

	pos := [2]int{1, 1} // position / head tracker for error reporting

	// // check if valid begining using Rune / Unicode values
	// checkVal := func(val rune) bool {
	// 	return (val >= '0' && val <= '9') || // test if number
	// 		(val >= 'A' && val <= 'Z') || // test if uppercase character
	// 		(val >= 'a' && val <= 'z') || // test if lowercase character
	// 		val == '_' // test if underscore
	// }

	for true {
		var charAdd int = 1
		var val string = ""

		var nextRune rune
		var firstRune rune

		for buildVal := true; buildVal; buildVal = MultiToken(firstRune, nextRune) {
			newRune, _, err := reader.ReadRune()
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Println("Error:", err) // temp
				}
				return
			}
			val += string(newRune)
			charAdd++

			nextVal, err := reader.Peek(1)
			if err != nil {
				if err.Error() != "EOF" {
					fmt.Println("Error:", err) // temp
				}
				break
			}

			firstRune = rune(val[0])
			nextRune = rune(nextVal[0])
		}

		var t TokenType

		// keywords tokenizing
		if keyType, ok := keywordMap[val]; ok { // if val is in keyword map
			t = keyType
		}

		// iden tokenizing
		if !(val[0] >= '0' && val[0] <= '9') && t == 0 {
			t = Iden
		}

		// literals tokenizing
		if val[0] >= '0' && val[0] <= '9' {
			t = Literal
		}

		// single & multi char tokenizing
		switch val {
		case " ":
			{
				pos[1] += 1
				continue
			} // ignore spaces
		case "\n":
			{
				pos[0] += 1
				pos[1] = 1
				continue
			} // ignore new lines
		case "\r":
			continue // ignore carriage returns (don't need to count for position tracking, always bundled with new line)
		case ",":
			t = Comma
		case "{":
			t = LCurly
		case "}":
			t = RCurly
		case "(":
			t = LParen
		case ")":
			t = RParen
		case ";":
			t = EOL
		case "=":
			t = Asmt
		case "==":
			t = Cmp
		case ">=":
			t = Cmp
		case "<=":
			t = Cmp
		}

		lex.Tokens <- Token{t, val, pos}
		pos[1] += charAdd
	}
}
