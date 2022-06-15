// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package Lexer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOL-0]
	_ = x[Iden-1]
	_ = x[Literal-2]
	_ = x[Direction-3]
	_ = x[Spec-4]
	_ = x[Default-5]
	_ = x[If-6]
	_ = x[Switch-7]
	_ = x[LParen-8]
	_ = x[RParen-9]
	_ = x[LCurly-10]
	_ = x[RCurly-11]
	_ = x[LBrace-12]
	_ = x[RBrace-13]
	_ = x[Atmark-14]
	_ = x[Math-15]
	_ = x[Comma-16]
	_ = x[Colon-17]
	_ = x[Asmt-18]
	_ = x[Cmp-19]
	_ = x[Unknown-20]
}

const _TokenType_name = "EOLIdenLiteralDirectionSpecDefaultIfSwitchLParenRParenLCurlyRCurlyLBraceRBraceAtmarkMathCommaColonAsmtCmpUnknown"

var _TokenType_index = [...]uint8{0, 3, 7, 14, 23, 27, 34, 36, 42, 48, 54, 60, 66, 72, 78, 84, 88, 93, 98, 102, 105, 112}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}