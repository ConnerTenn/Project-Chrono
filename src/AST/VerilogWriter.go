package AST

import "fmt"

//== Module ==
func (mod Module) WriteVerilog(ident int) string {
	var str string
	str += "module " + mod.Name + "(\n"

	//Generate parameter list
	for i, param := range mod.Params {
		str += Ident(ident+1) + param.Dir.String() + " "
		//Only include signal width if the width is > 1
		if param.Width > 1 {
			str += "[" + fmt.Sprint(param.Width-1) + ":0] "
		}
		str += param.Name

		//The last parameter must not have a comma
		if i == len(mod.Params)-1 {
			str += "\n"
		} else {
			str += ",\n"
		}
	}

	str += ");\n"

	//Module contents
	str += mod.Block.WriteVerilog(ident)

	str += "endmodule\n"

	return str
}

//== Block ==
func (blk Block) WriteVerilog(ident int) string {
	var str string

	//Generate each element within the block
	for _, elem := range blk.Elements {
		str += elem.WriteVerilog(ident+1) + "\n"
	}

	return str
}

//== ValueExpression ==
func (expr ValueExpression) WriteVerilog(ident int) string {
	return expr.Value
}

//== AssignmentExpression ==
func (expr AssignmentExpression) WriteVerilog(ident int) string {
	return Ident(ident) + "assign " + expr.Name + " = " + expr.RHS.WriteVerilog(0) + ";"
}

//== MathExpression ==
func (expr MathExpression) WriteVerilog(ident int) string {
	var op string
	switch expr.Op {
	case Add:
		op = " + "
	case Sub:
		op = " - "
	case Multi:
		op = " * "
	case Div:
		op = " / "
	case LShift:
		op = " << "
	case RShift:
		op = " >> "
	}

	return expr.LHS.WriteVerilog(0) + op + expr.RHS.WriteVerilog(0)
}
