// Code generated by "stringer -type=SignalDir"; DO NOT EDIT.

package AST

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[In-0]
	_ = x[Out-1]
	_ = x[Inout-2]
}

const _SignalDir_name = "InOutInout"

var _SignalDir_index = [...]uint8{0, 2, 5, 10}

func (i SignalDir) String() string {
	if i < 0 || i >= SignalDir(len(_SignalDir_index)-1) {
		return "SignalDir(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SignalDir_name[_SignalDir_index[i]:_SignalDir_index[i+1]]
}
