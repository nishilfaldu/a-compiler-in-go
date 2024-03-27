package compiler

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/object"
)

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

func NewWithState(s *SymbolTable) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		err := c.Compile(node.Header)
		if err != nil {
			return err
		}

		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

		// case *ast.ProgramHeader:
	}

	return nil

}
