package compiler

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/object"
	"fmt"
	"sort"
)

type Compiler struct {
	symbolTable *SymbolTable
}

type CompileResult struct {
	Type string
}

// TODO: check if the expression in the infix expression is a boolean or not.
// TODO: Point 7 and 14 on the doc: implicit conversions of boolean and integer - implemented but check still...

func New() *Compiler {
	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		// TODO: this will have to change to ReturnType - but for some reason Go typing is not working
		if v.Name == "putbool" || v.Name == "putinteger" || v.Name == "putfloat" || v.Name == "putstring" || v.Name == "getbool" {
			symbolTable.DefineBuiltin(i, v.Name, "bool")
		} else if v.Name == "getinteger" {
			symbolTable.DefineBuiltin(i, v.Name, "integer")
		} else if v.Name == "getfloat" {
			symbolTable.DefineBuiltin(i, v.Name, "float")
		} else if v.Name == "getstring" {
			symbolTable.DefineBuiltin(i, v.Name, "string")
		} else {
			symbolTable.DefineBuiltin(i, v.Name, "")
		}
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
			symbol := c.symbolTable.Define(node.Name.Value, node.Type.Name, false)
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
		print(symbol.Scope, symbol.Type, " in Identifier case i just need scope\n")
		if !ok {
			return CompileResult{}, fmt.Errorf("undefined variable %s", node.Value)
		}
		if symbol.Scope == FunctionScope {
			return CompileResult{Type: symbol.Type}, nil
		}
		if symbol.Scope == BuiltinScope {
			return CompileResult{Type: symbol.Type}, nil
		}
		// if symbol.Scope == GlobalScope {
		// 	return CompileResult{Type: symbol.Type}, nil
		// }
		print(symbol.Name, symbol.Type, symbol.Index, symbol.Scope, "in Identifier case\n")
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
		print(node.Operator, " - operator\n")
		fmt.Printf("Type of curr node in prefix: %T\n", node.Right)
		cr, err := c.Compile(node.Right)
		if err != nil {
			return CompileResult{}, err
		}

		return CompileResult{Type: cr.Type}, nil

	case *ast.InfixExpression:
		print(node.Operator, " - operator\n")
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
		fmt.Printf("Type of curr node in infix - left: %T\n", node.Left)
		cr, err := c.Compile(node.Left)
		if err != nil {
			return CompileResult{}, err
		}

		cr_, err_ := c.Compile(node.Right)
		if err_ != nil {
			return CompileResult{}, err_
		}

		// check if left and right types match
		if cr.Type != cr_.Type {
			return CompileResult{}, fmt.Errorf("type mismatch: cannot perform operation %s on %s and %s", node.Operator, cr.Type, cr_.Type)
		} else {
			return CompileResult{Type: cr.Type}, nil
		}

	case *ast.IfExpression:
		print(node.Condition, " - condition\n")
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
		print(len(node.Statements), " - len of statements\n")
		for _, stmt := range node.Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return CompileResult{}, err
			}
		}

	// AssignmentStatement node and Destination node are merged in one case
	case *ast.AssignmentStatement:
		fmt.Printf("Type of curr node in assignment statement: %T\n", node.Value)
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
		fmt.Printf("Type of curr node in return: %T\n", node.ReturnValue)
		cr, err := c.Compile(node.ReturnValue)
		if err != nil {
			return CompileResult{}, err
		}

		function, ok := c.symbolTable.getCurrentFunction()
		if !ok {
			return CompileResult{}, fmt.Errorf("return statement outside of function")
		}
		if cr.Type != function.ReturnType {
			if typesCompatible(cr.Type, function.ReturnType) {
				return CompileResult{Type: function.ReturnType}, nil
			}
			return CompileResult{}, fmt.Errorf("type mismatch for function %s: cannot return %s from function of type %s", function.Name, cr.Type, function.ReturnType)
		}

		return CompileResult{Type: cr.Type}, nil

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		print(str)
		return CompileResult{Type: string(str.Type())}, nil

	case *ast.IntegerLiteral:
		fmt.Printf("Type of curr node in integer literal: %T\n", node.Value)
		integer := &object.Integer{Value: node.Value}
		return CompileResult{Type: string(integer.Type())}, nil

	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		return CompileResult{Type: string(float.Type())}, nil

	case *ast.Boolean:
		boolean := &object.Boolean{Value: node.Value}
		return CompileResult{Type: string(boolean.Type())}, nil

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
		fmt.Printf("Type of curr node in call expression: %T\n", node.Function)
		cr, err := c.Compile(node.Function)
		if err != nil {
			return CompileResult{}, err
		}

		currentFuncName := node.Function.String()
		if builtinWithExists(currentFuncName) {
			return c.compileBuiltInFunction(node)
		}

		paramLocalSymbols := getParamLocalSymbols(c.symbolTable, node.Function.String())
		print(len(node.Arguments), " - len of arguments\n")
		// Check if there are enough local symbols for the arguments
		if len(node.Arguments) < len(paramLocalSymbols) {
			return CompileResult{}, fmt.Errorf("not enough arguments provided for function call")
		} else if len(node.Arguments) > len(paramLocalSymbols) {
			return CompileResult{}, fmt.Errorf("too many arguments provided for function call")
		}

		for i, a := range node.Arguments {
			fmt.Printf("Type of curr node in call expression - arguments: %T\n", a)
			cr, err := c.Compile(a)
			if err != nil {
				return CompileResult{}, err
			}

			if cr.Type != paramLocalSymbols[i].Type {
				return CompileResult{}, fmt.Errorf("type mismatch: cannot pass %s as argument %d of type %s", cr.Type, i, paramLocalSymbols[i].Type)
			}
		}

		return CompileResult{Type: cr.Type}, nil

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
		// popping a function from the stack after return and its compilation
		c.symbolTable.popFunction()
		// c.leaveScope()

	case *ast.ProcedureHeader:
		c.enterScope()

		// define the function name and parameters in the symbol table
		if node.Name.Value != "" {
			c.symbolTable.DefineFunctionName(node.Name.Value, node.TypeMark.Name)
			// push function onto the stack for tracking return type
			c.symbolTable.pushFunction(node.Name.Value, node.TypeMark.Name)
		}

		for _, param := range node.Parameters {
			c.symbolTable.Define(param.Name.Value, param.Type.Name, true)
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

func sortParamLocalSymbols(localSymbols []Symbol) {
	sort.Slice(localSymbols, func(i, j int) bool {
		return localSymbols[i].Index < localSymbols[j].Index
	})
}

func getParamLocalSymbols(symbolTable *SymbolTable, functionName string) []Symbol {
	functionScope := findFunctionScope(symbolTable, functionName)
	if functionScope == nil {
		// Function not found, return empty slice
		return []Symbol{}
	}

	localSymbols := make([]Symbol, 0)
	for _, sym := range functionScope.store {
		if sym.Scope == ParamLocalScope {
			print(sym.Name + " haha\n")
			localSymbols = append(localSymbols, sym)
		}
	}
	sortParamLocalSymbols(localSymbols)

	return localSymbols
}

// Find the symbol table containing the function definition
func findFunctionScope(symbolTable *SymbolTable, functionName string) *SymbolTable {
	current := symbolTable
	for current != nil {
		// Check if the current symbol table contains the function definition
		if _, ok := current.store[functionName]; ok {
			return current
		}
		// Move to the outer symbol table
		current = current.Outer
	}
	// Function scope not found
	return nil
}

func builtinWithExists(name string) bool {
	for _, builtin := range object.Builtins {
		if builtin.Name == name {
			return true
		}
	}
	return false
}

func (c *Compiler) compileBuiltInFunction(node *ast.CallExpression) (CompileResult, error) {
	switch currentFuncName := node.Function.String(); currentFuncName {
	case "putinteger":
		return c.checkEnoughArgumentsAndCompile(node, 1, "bool", "integer")
	case "putfloat":
		return c.checkEnoughArgumentsAndCompile(node, 1, "bool", "float")
	case "putstring":
		return c.checkEnoughArgumentsAndCompile(node, 1, "bool", "string")
	case "putbool":
		return c.checkEnoughArgumentsAndCompile(node, 1, "bool", "bool")
	case "sqrt":
		return c.checkEnoughArgumentsAndCompile(node, 1, "float", "integer")
	case "getinteger":
		return c.checkEnoughArgumentsAndCompile(node, 0, "integer", "")
	case "getfloat":
		return c.checkEnoughArgumentsAndCompile(node, 0, "float", "")
	case "getstring":
		return c.checkEnoughArgumentsAndCompile(node, 0, "string", "")
	case "getbool":
		return c.checkEnoughArgumentsAndCompile(node, 0, "bool", "")
	default:
		return CompileResult{}, fmt.Errorf("unknown built-in function: %s", currentFuncName)
	}
}

func (c *Compiler) checkEnoughArgumentsAndCompile(node *ast.CallExpression, expectedArgs int, returnType string, argType string) (CompileResult, error) {
	if len(node.Arguments) != expectedArgs {
		return CompileResult{}, fmt.Errorf("wrong number of arguments for %s: got %d, want %d", node.Function.String(), len(node.Arguments), expectedArgs)
	} else {
		// Check if the argument is an integer
		if expectedArgs == 1 {
			cr, err := c.Compile(node.Arguments[0])
			if err != nil {
				return CompileResult{}, err
			}
			if cr.Type != argType {
				return CompileResult{}, fmt.Errorf("wrong type of argument for %s: got %s, want %s", node.Function.String(), cr.Type, argType)
			}
			return CompileResult{Type: returnType}, nil
		} else if expectedArgs == 0 {
			return CompileResult{Type: returnType}, nil
		}

		return CompileResult{Type: returnType}, nil
	}
}

func typesCompatible(type1 string, type2 string) bool {
	// Check for compatibility between bool and integer types
	if (type1 == "bool" && type2 == "integer") ||
		(type1 == "integer" && type2 == "bool") {
		return true
	} else if (type1 == "integer" && type2 == "float") ||
		(type1 == "float" && type2 == "integer") {
		return true
	}
	return false
}
