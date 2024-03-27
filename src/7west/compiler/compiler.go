package compiler

import "a-compiler-in-go/src/7west/src/7west/object"

type Compiler struct {
	symbolTable *SymbolTable
}

func New() *Compiler {
	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		symbolTable: symbolTable,
	}
}
