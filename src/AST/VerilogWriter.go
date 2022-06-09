package AST

func (mod Module) WriteVerilog(ident int) string {
	var str string
	str += "module " + mod.Name + "(\n"

	for i, param := range mod.Params {
		str += Ident(ident+1) + param.Dir.String() + " " + param.Name

		if i == len(mod.Params)-1 {
			str += "\n"
		} else {
			str += ",\n"
		}
	}

	str += ");\n"

	str += mod.Block.WriteVerilog(ident)

	str += "endmodule\n"

	return str
}

func (blk Block) WriteVerilog(ident int) string {
	var str string

	for _, elem := range blk.Elements {
		str += elem.WriteVerilog(ident+1) + "\n"
	}

	return str
}

func (expr ValueExpression) WriteVerilog(ident int) string {
	return expr.Value
}

func (expr AssignmentExpression) WriteVerilog(ident int) string {
	return Ident(ident) + "assign " + expr.Name + " = " + expr.RHS.WriteVerilog(0) + ";"
}

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
