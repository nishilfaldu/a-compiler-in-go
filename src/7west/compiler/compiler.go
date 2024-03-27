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
	print(node.String(), " in Compile\n")
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
		print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case\n")

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		PrintSymbolTable(c.symbolTable)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		print(symbol.Name, symbol.Index, symbol.Scope, "in Identifier case\n")

	case *ast.LoopStatement:
		// Compile the initialization statement
		err := c.Compile(node.InitStatement)
		if err != nil {
			return err
		}
		// Compile the loop condition
		err_ := c.Compile(node.Condition)
		if err_ != nil {
			return err_
		}

		// Compile the loop body
		err = c.Compile(node.Body)
		if err != nil {
			return err
		}

	case *ast.ForBlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err_ := c.Compile(node.Left)
			if err_ != nil {
				return err_
			}
		}

		print("i executed lmaoo found it\n")

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err_ := c.Compile(node.Right)
		if err_ != nil {
			return err_
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
		print(node.Value.String(), " in AssignmentStatement case 2\n")
		err := c.Compile(node.Value)
		if err != nil {
			print("here's the error\n")
			return err
		}
		// Compile the identifier part of the destination
		symbol, ok := c.symbolTable.Resolve(node.Destination.Identifier.Value)

		print(symbol.Name, symbol.Index, symbol.Scope, "in AssignmentStatement case\n")
		if !ok {
			// If the identifier is not found in the local symbol table, check the outer scopes
			symbol = c.symbolTable.Define(node.Destination.Identifier.Value)
		}

		// If the assignment has an index expression - array indexing - compile it first
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
		print(str)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		print(integer)

	case *ast.Boolean:
		// TODO: not sure if this is required
		// boolean := &object.Boolean{Value: node.Value}

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
		print(node.Function.String(), " . ", node.String(), " in CallExpression case\n")
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
		print(numLocals)
		// c.leaveScope()

	case *ast.ProcedureHeader:
		c.enterScope()

		// define the function name and parameters in the symbol table
		if node.Name.Value != "" {
			print("i executed\n")
			c.symbolTable.DefineFunctionName(node.Name.Value)
			s, ok := c.symbolTable.Resolve(node.Name.Value)
			if !ok {
				print(s.Name, " i executed too yayy\n")
			}
		}

		PrintSymbolTable(c.symbolTable)

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
