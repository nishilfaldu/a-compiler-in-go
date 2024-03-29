package compiler

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/object"
	"fmt"
)

type Compiler struct {
	symbolTable *SymbolTable
}

type CompileResult struct {
	Type string
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

func NewWithState() *Compiler {
	compiler := New()
	// compiler.symbolTable = s
	return compiler
}

func (c *Compiler) Compile(node ast.Node) (CompileResult, error) {
	switch node := node.(type) {
	case *ast.Program:
		_, err := c.Compile(node.Header)
		if err != nil {
			return CompileResult{}, err
		}

		_, err = c.Compile(node.Body)
		if err != nil {
			return CompileResult{}, err
		}

	case *ast.ProgramHeader:
		// No compilation needed for the program header
		// as it typically contains metadata about the program.
		// You can optionally perform any necessary validation or processing here.

	case *ast.ProgramBody:
		for _, decl := range node.Declarations {
			_, err := c.Compile(decl)
			if err != nil {
				return CompileResult{}, err
			}
		}

		for _, stmt := range node.Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.VariableDeclaration:
		if node.Type.Array != nil {
			// Handle array declaration
			// First, compile the inner variable declaration
			// Then, define the symbol in the symbol table as an array
			symbol := c.symbolTable.DefineArray(node.Name.Value, node.Type.Name+"[]", node.Type.Array.Value)
			print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case\n")
		} else {
			// Handle variable declaration
			// First, compile the inner variable declaration
			// Then, define the symbol in the symbol table as a variable
			symbol := c.symbolTable.Define(node.Name.Value, node.Type.Name)
			print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case\n")
		}
		// symbol := c.symbolTable.Define(node.Name.Value, node.Type.Name)
		// err := c.Compile(node.)
		// print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case\n")

	case *ast.GlobalVariableDeclaration:
		// Handle global variable declaration
		// First, compile the inner variable declaration
		// Then, define the symbol in the symbol table as a global variable
		symbol := c.symbolTable.DefineGlobal(node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name)
		print(symbol.Name, symbol.Index, symbol.Scope, "in Global Variable Declaration case\n")

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return CompileResult{}, fmt.Errorf("undefined variable %s", node.Value)
		}
		print(symbol.Name, symbol.Index, symbol.Scope, "in Identifier case\n")
		return CompileResult{Type: symbol.Type}, nil

	case *ast.LoopStatement:
		// Compile the initialization statement
		_, err := c.Compile(node.InitStatement)
		if err != nil {
			return CompileResult{}, err
		}
		// Compile the loop condition
		_, err_ := c.Compile(node.Condition)
		if err_ != nil {
			return CompileResult{}, err_
		}

		// Compile the loop body
		_, err = c.Compile(node.Body)
		if err != nil {
			return CompileResult{}, err
		}

	case *ast.ForBlockStatement:
		for _, stmt := range node.Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.PrefixExpression:
		_, err := c.Compile(node.Right)
		if err != nil {
			return CompileResult{}, err
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			_, err := c.Compile(node.Right)
			if err != nil {
				return CompileResult{}, err
			}

			_, err_ := c.Compile(node.Left)
			if err_ != nil {
				return CompileResult{}, err_
			}
		}

		_, err := c.Compile(node.Left)
		if err != nil {
			return CompileResult{}, err
		}

		_, err_ := c.Compile(node.Right)
		if err_ != nil {
			return CompileResult{}, err_
		}

	case *ast.IfExpression:
		_, err := c.Compile(node.Condition)
		if err != nil {
			return CompileResult{}, err
		}
		_, err_ := c.Compile(node.Consequence)
		if err_ != nil {
			return CompileResult{}, err
		}

		if node.Alternative != nil {
			_, err := c.Compile(node.Alternative)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.IfBlockStatement:
		for _, stmt := range node.Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return CompileResult{}, err
			}
		}

	// AssignmentStatement node and Destination node are merged in one case
	case *ast.AssignmentStatement:
		fmt.Printf("Type of curr node: %T\n", node.Value)
		cr, err := c.Compile(node.Value)
		if err != nil {
			return CompileResult{}, err
		}
		// Compile the identifier part of the destination
		symbol, ok := c.symbolTable.Resolve(node.Destination.Identifier.Value)
		print(symbol.Name, symbol.Index, symbol.Scope, "in AssignmentStatement case - print for usage\n")
		if !ok {
			return CompileResult{}, fmt.Errorf("variable %s not defined", node.Destination.Identifier.Value)
		}

		print(symbol.Type, cr.Type, "hello symbol type here\n")
		// check if types match
		if cr.Type != symbol.Type {
			return CompileResult{}, fmt.Errorf("type mismatch: cannot assign %s to %s", cr.Type, symbol.Type)
		}

		// If the assignment has an index expression - array indexing - compile it first
		if node.Destination.Expression != nil {
			_, err := c.Compile(node.Destination.Expression)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.ExpressionStatement:
		_, err := c.Compile(node.Expression)
		if err != nil {
			return CompileResult{}, err
		}

	case *ast.ReturnStatement:
		_, err := c.Compile(node.ReturnValue)
		if err != nil {
			return CompileResult{}, err
		}

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		print(str)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		return CompileResult{Type: string(integer.Type())}, nil

	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		return CompileResult{Type: string(float.Type())}, nil

	case *ast.Boolean:
		// TODO: not sure if this is required
		// boolean := &object.Boolean{Value: node.Value}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			_, err := c.Compile(el)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.IndexExpression:
		_, err := c.Compile(node.Left)
		if err != nil {
			return CompileResult{}, err
		}
		_, err_ := c.Compile(node.Index)
		if err_ != nil {
			return CompileResult{}, err
		}

	case *ast.CallExpression:
		_, err := c.Compile(node.Function)
		if err != nil {
			return CompileResult{}, err
		}

		for _, a := range node.Arguments {
			_, err := c.Compile(a)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.ProcedureDeclaration:
		_, err := c.Compile(node.Header)
		if err != nil {
			return CompileResult{}, err
		}
		_, err = c.Compile(node.Body)
		if err != nil {
			return CompileResult{}, err
		}

		numLocals := c.symbolTable.numDefinitions
		print(numLocals)
		// c.leaveScope()

	case *ast.ProcedureHeader:
		c.enterScope()

		// define the function name and parameters in the symbol table
		if node.Name.Value != "" {
			c.symbolTable.DefineFunctionName(node.Name.Value)
		}

		for _, param := range node.Parameters {
			c.symbolTable.Define(param.Name.Value, param.Type.Name)
		}

	case *ast.ProcedureBody:
		for _, decl := range node.Declarations {
			_, err := c.Compile(decl)
			if err != nil {
				return CompileResult{}, err
			}
		}

		for _, stmt := range node.Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return CompileResult{}, err
			}
		}

		PrintSymbolTable(c.symbolTable)
	}

	return CompileResult{}, nil

}

func (c *Compiler) enterScope() {
	// 	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
	// c.symbolTable = c.symbolTable.NewChildSymbolTable()
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() {
	c.symbolTable = c.symbolTable.Outer
}

// getTypeOfNode returns the type of the AST node.
// This function assumes that the AST node has a Type field representing its type.
func getTypeOfNode(node ast.Node) string {
	switch node.(type) {
	case *ast.IntegerLiteral:
		return "integer"
	case *ast.Boolean:
		return "boolean"
	case *ast.StringLiteral:
		return "string"
	case *ast.FloatLiteral:
		return "float"
	// Add cases for other types as needed...
	default:
		return "" // Return empty string if type is unknown
	}
}
