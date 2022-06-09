package AST

func (mod Module) WriteVerilog(ident int) string {
	var str string
	str += "module " + mod.Name + "(\n"

	for i, param := range mod.Params {
		str += Ident(ident+1) + param.Dir.String() + " " + param.Name

		if i == len(mod.Params)-1 {
			str += ",\n"
		} else {
			str += "\n"
		}
	}

	str += ");\n"

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
	return ""
}

func (expr AssignmentExpression) WriteVerilog(ident int) string {
	return ""
}

func (expr MathExpression) WriteVerilog(ident int) string {
	return ""
}
