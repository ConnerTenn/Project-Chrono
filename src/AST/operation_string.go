// Code generated by "stringer -type=Operation"; DO NOT EDIT.

package AST

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Asmt-0]
	_ = x[AsmtReg-1]
	_ = x[LShift-2]
	_ = x[RShift-3]
	_ = x[Add-4]
	_ = x[Sub-5]
	_ = x[Multi-6]
	_ = x[Div-7]
	_ = x[Bracket-8]
	_ = x[Equals-9]
}

const _Operation_name = "AsmtAsmtRegLShiftRShiftAddSubMultiDivBracketEquals"

var _Operation_index = [...]uint8{0, 4, 11, 17, 23, 26, 29, 34, 37, 44, 50}

func (i Operation) String() string {
	if i < 0 || i >= Operation(len(_Operation_index)-1) {
		return "Operation(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Operation_name[_Operation_index[i]:_Operation_index[i+1]]
}
