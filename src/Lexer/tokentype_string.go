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
	_ = x[Else-7]
	_ = x[Switch-8]
	_ = x[LParen-9]
	_ = x[RParen-10]
	_ = x[LCurly-11]
	_ = x[RCurly-12]
	_ = x[LBrace-13]
	_ = x[RBrace-14]
	_ = x[Atmark-15]
	_ = x[Math-16]
	_ = x[Comma-17]
	_ = x[Colon-18]
	_ = x[Asmt-19]
	_ = x[Cmp-20]
	_ = x[Unknown-21]
}

const _TokenType_name = "EOLIdenLiteralDirectionSpecDefaultIfElseSwitchLParenRParenLCurlyRCurlyLBraceRBraceAtmarkMathCommaColonAsmtCmpUnknown"

var _TokenType_index = [...]uint8{0, 3, 7, 14, 23, 27, 34, 36, 40, 46, 52, 58, 64, 70, 76, 82, 88, 92, 97, 102, 106, 109, 116}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
