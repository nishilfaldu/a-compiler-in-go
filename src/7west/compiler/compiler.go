package compiler

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/object"
	"fmt"
)

type Compiler struct {
	symbolTable *SymbolTable
}

// TODO: store types in symbol table...maybe?

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

	case *ast.ProgramHeader:
		// No compilation needed for the program header
		// as it typically contains metadata about the program.
		// You can optionally perform any necessary validation or processing here.

	case *ast.ProgramBody:
		for _, decl := range node.Declarations {
			err := c.Compile(decl)
			if err != nil {
				return err
			}
		}

		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}

	case *ast.VariableDeclaration:
		symbol := c.symbolTable.Define(node.Name.Value)
		// err := c.Compile(node.)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

	case *ast.ProcedureDeclaration:
		err := c.Compile(node.Header)
		if err != nil {
			return err
		}
		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

	case *ast.ProcedureHeader:
		c.enterScope()

		// define the function name and parameters in the symbol table
		if node.Name != nil {
			c.symbolTable.DefineFunctionName(node.Name.Value)
		}

		for _, param := range node.Parameters {
			c.symbolTable.Define(param.Name.Value)
		}

	case *ast.ProcedureBody:
		for _, decl := range node.Declarations {
			err := c.Compile(decl)
			if err != nil {
				return err
			}
		}
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

func (c *Compiler) enterScope() {
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() {
	c.symbolTable = c.symbolTable.Outer
}
