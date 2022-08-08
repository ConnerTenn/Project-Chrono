package verilog

import AST "github.com/ConnerTenn/Project-Chrono/AST"

type VariableScope struct {
	Variables [][]AST.SignalDecl
}

func (scope *VariableScope) DeclVariable(variable AST.SignalDecl) {
	currScope := &(scope.Variables[len(scope.Variables)-1])
	*currScope = append(*currScope, variable)
}

func (scope *VariableScope) EnterScope() {
	scope.Variables = append(scope.Variables, []AST.SignalDecl{})
}

func (scope *VariableScope) ExitScope() {
	scope.Variables = scope.Variables[0 : len(scope.Variables)-1]
}
