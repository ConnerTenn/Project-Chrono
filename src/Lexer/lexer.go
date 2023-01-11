package Lexer

import (
	"bufio"
	"fmt"
	"os"
)

// token types split into enum per type for easier parsing (on switch)
//go:generate stringer -type=TokenType
type TokenType int

const (
	EOL       TokenType = iota
	Iden                // any named identifier
	Literal             // can be split into the different literal types during lexing
	Direction           // signal direction, input or output
	Spec                // Wire / Reg / Param, specifier for 'variables'
	Default             // Default case
	If
	Else
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
	Asmt //Assignment
	Cmp
	Unknown
)

type Token struct {
	// Type  TokenType
	Value string // store the direct value
	Pos   [2]int // store the position for error reporting
}

// stringer interface
func (t Token) String() string {
	return fmt.Sprintf("%v at %d:%d", t.Value, t.Pos[0], t.Pos[1])
}

func (t Token) Is(str string) bool {
	return t.Value == str
}

//Not including keywords
func (t Token) IsIden() bool {
	if t.IsKeyword() {
		return false
	}
	if (t.Value[0] >= 'A' && t.Value[0] <= 'Z') || (t.Value[0] >= 'a' && t.Value[0] <= 'z') {
		return true
	}
	return false
}

func (t Token) IsLiteral() bool {
	if t.Value[0] >= '0' && t.Value[0] <= '9' {
		return true
	}
	return false
}

func (t Token) IsMath() bool {
	switch t.Value {
	case "+", "-", "*", "/":
		return true
	}
	return false
}

func (t Token) IsLogical() bool {
	switch t.Value {
	case "!", "&&", "||":
		return true
	}
	return false
}

func (t Token) IsBitwise() bool {
	switch t.Value {
	case "~", "&", "|", ">>", "<<":
		return true
	}
	return false
}

func (t Token) IsComparison() bool {
	switch t.Value {
	case "==", "<", ">", "<=":
		return true
	}
	return false
}

func (t Token) IsOperator() bool {
	return t.IsMath() || t.IsLogical() || t.IsBitwise() || t.IsComparison()
}

func (t Token) IsLParen() bool {
	return t.Value == "("
}
func (t Token) IsRParen() bool {
	return t.Value == ")"
}
func (t Token) IsParen() bool {
	return t.IsLParen() || t.IsRParen()
}

func (t Token) IsLBracket() bool {
	return t.Value == "["
}
func (t Token) IsRBracket() bool {
	return t.Value == "]"
}
func (t Token) IsBrace() bool {
	return t.IsLBracket() || t.IsRBracket()
}

func (t Token) IsLCurly() bool {
	return t.Value == "{"
}
func (t Token) IsRCurly() bool {
	return t.Value == "}"
}
func (t Token) IsCurly() bool {
	return t.IsLCurly() || t.IsRCurly()
}

func (t Token) IsAssignment() bool {
	switch t.Value {
	case "=", "<-":
		return true
	}
	return false
}

func (t Token) IsKeyword() bool {
	switch t.Value {
	case "sig", "in", "out", "inout", "if", "else":
		return true
	}
	return false
}

func (t Token) IsSpec() bool {
	switch t.Value {
	case "in", "out", "inout":
		return true
	}
	return false
}

func (t Token) IsEOL() bool {
	switch t.Value {
	case ";", "\n":
		return true
	}
	return false
}

func (t Token) GetType() TokenType {
	if t.IsIden() {
		return Iden
	}
	if t.IsLiteral() && !t.IsKeyword() {
		return Literal
	}
	return tokenMap[t.Value]
}

// keyword list
var tokenMap = map[string]TokenType{
	"in":      Direction,
	"out":     Direction,
	"inout":   Direction,
	"wire":    Spec,
	"reg":     Spec,
	"var":     Spec,
	"if":      If,
	"else":    Else,
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
	"<<":      Math,
	">>":      Math,
	"!":       Math,
	"=":       Asmt,
	"<-":      Asmt,
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
func (lex *Lexer) GetNext() Token {
	return lex.tokens.GetNext()
}

func (lex *Lexer) PeekNext() Token {
	return lex.tokens.PeekNext()
}

// func (lex *Lexer) ExpectNext(t TokenType) bool {
// 	nextToken := lex.PeekNext()
// 	return nextToken.Type == t
// }
func (lex *Lexer) ExpectNextType(test func(Token) bool) bool {
	nextToken := lex.PeekNext()
	return test(nextToken)
}

func (lex *Lexer) ExpectNext(str string) bool {
	nextToken := lex.PeekNext()
	return nextToken.Is(str)
}

func (lex *Lexer) NextExists() bool {
	return !lex.tokens.IsEmpty()
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

	//Check reg assigment
	if first == '<' && next == '-' {
		return true
	}

	return false
}

func (lex Lexer) Tokenizer() {
	reader := bufio.NewReader(lex.file)

	defer lex.file.Close()
	defer lex.tokens.Close()

	pos := [2]int{1, 1} // position / head tracker for error reporting

	for {
		var (
			charAdd int    = 1
			val     string = ""

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

		lex.tokens.PushBack(Token{val, pos})
		pos[1] += charAdd
	}
}
