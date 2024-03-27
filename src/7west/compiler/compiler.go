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

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		err_ := c.Compile(node.Consequence)
		if err_ != nil {
			return err_
		}

		if node.Alternative != nil {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}
		}

	case *ast.IfBlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}

	// AssignmentStatement node and Destination node are merged in one case
	case *ast.AssignmentStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// Compile the identifier part of the destination
		symbol, ok := c.symbolTable.Resolve(node.Destination.Identifier.Value)
		if !ok {
			// If the identifier is not found in the local symbol table, check the outer scopes
			symbol = c.symbolTable.Define(node.Destination.Identifier.Value)
		}

		// If the assignment has an index expression, compile it first
		if node.Destination.Expression != nil {
			err := c.Compile(node.Destination.Expression)
			if err != nil {
				return err
			}
		}

		// Compile the identifier part of the assignment
		// err = c.Compile(node.Destination.Identifier)
		// if err != nil {
		// 	return err
		// }

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

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err_ := c.Compile(node.Index)
		if err_ != nil {
			return err_
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

		numLocals := c.symbolTable.numDefinitions
		c.leaveScope()

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
