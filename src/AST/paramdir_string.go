// Code generated by "stringer -type=ParamDir"; DO NOT EDIT.

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

const _ParamDir_name = "InOutInout"

var _ParamDir_index = [...]uint8{0, 2, 5, 10}

func (i ParamDir) String() string {
	if i < 0 || i >= ParamDir(len(_ParamDir_index)-1) {
		return "ParamDir(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ParamDir_name[_ParamDir_index[i]:_ParamDir_index[i+1]]
}
