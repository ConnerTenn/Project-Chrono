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
	Spec                // Wire / Reg / Param, specifier for 'variables'
	Default             // Default case
	If
	Switch
	LParen
	RParen
	LCurly
	RCurly
	LBrace
	RBrace
	Atmark
	Math
	Comma
	Colon
	Asmt
	Cmp
	Unknown
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
var tokenMap = map[string]TokenType{
	"in":      Direction,
	"out":     Direction,
	"inout":   Direction,
	"wire":    Spec,
	"reg":     Spec,
	"param":   Spec,
	"if":      If,
	"switch":  Switch,
	"default": Default,
	",":       Comma,
	"{":       LCurly,
	"}":       RCurly,
	"(":       LParen,
	")":       RParen,
	"[":       LBrace,
	"]":       RBrace,
	";":       EOL,
	"@":       Atmark,
	":":       Colon,
	"*":       Math,
	"/":       Math,
	"-":       Math,
	"+":       Math,
	"=":       Asmt,
	"==":      Cmp,
	">=":      Cmp,
	"<=":      Cmp,
}

// could be done functionally :)
// lexer consumes file line by line and sends tokens to a channel to be consumed by the parser
type Lexer struct {
	file   *os.File
	tokens PeekableQueue
}

func NewLexer(fileName string) (Lexer, error) {
	var (
		lex Lexer
		err error
	)
	lex.tokens = NewQueue()
	lex.file, err = os.Open(fileName)

	return lex, err // err will be nil if nothing was thrown, no need to check here
}

// provide an interface over the PeekableQueue so it doesn't have to be directly exported
func (lex Lexer) GetNext() (Token, bool) {
	return lex.tokens.GetNext()
}

func (lex Lexer) PeekNext() Token {
	return lex.tokens.PeekNext()
}

func (lex Lexer) ExpectNext(t TokenType) bool {
	nextToken := lex.tokens.PeekNext()
	return nextToken.Type == t
}

func multiToken(first rune, next rune) bool {
	// check if rune is part of a Name or Value
	checkNameVal := func(val rune) bool {
		return (val >= '0' && val <= '9') || // test if number
			(val >= 'A' && val <= 'Z') || // test if uppercase character
			(val >= 'a' && val <= 'z') || // test if lowercase character
			val == '_' // test if underscore
	}

	//Check if rune is part of a comparison
	checkCmp := func(val rune) bool {
		return (val == '<' || val == '=' || val == '>')
	}

	//Check if this is a valid Name/Value multi token
	if checkNameVal(first) && checkNameVal(next) {
		return true
	}

	//Check if this is a valid Comparison multi token
	if checkCmp(first) && checkCmp(next) {
		return true
	}

	return false
}

func (lex Lexer) Tokenizer() {
	reader := bufio.NewReader(lex.file)

	defer lex.file.Close()
	defer lex.tokens.Close()

	pos := [2]int{1, 1} // position / head tracker for error reporting

	for true {
		var (
			charAdd int       = 1
			val     string    = ""
			t       TokenType = Unknown

			nextRune  rune
			firstRune rune
		)

		//Do While type loop. Is guaranteed to execute at least once.
		//Will continue to loop based on the MultiToken condition
		for buildVal := true; buildVal; buildVal = multiToken(firstRune, nextRune) {
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

		// keywords tokenizing
		if keyType, ok := tokenMap[val]; ok { // if val is in keyword map
			t = keyType
		}

		// iden tokenizing
		if ((val[0] >= 'A' && val[0] <= 'Z') || (val[0] >= 'a' && val[0] <= 'z')) && t == Unknown {
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
		}

		lex.tokens.PushBack(Token{t, val, pos})
		pos[1] += charAdd
	}
}
