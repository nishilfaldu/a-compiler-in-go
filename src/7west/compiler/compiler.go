package compiler

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/llirgen"
	"a-compiler-in-go/src/7west/src/7west/object"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type FuncBlock struct {
	block    *ir.Block
	function *ir.Func
}

// context definition starts here
type Context struct {
	*ir.Block
	parent     *Context
	vars       map[string]value.Value
	leaveBlock *ir.Block
}

// this can be used to create initial context
func NewContext(b *ir.Block) *Context {
	return &Context{
		Block:      b,
		parent:     nil,
		vars:       make(map[string]value.Value),
		leaveBlock: nil,
	}
}

// this can used to create new context with parent context
func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b)
	ctx.parent = c
	return ctx
}

func (c *Context) lookupVariable(name string) value.Value {
	if v, ok := c.vars[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.lookupVariable(name)
	} else {
		fmt.Printf("variable: `%s`\n", name)
		panic("no such variable")
	}
}

type Compiler struct {
	symbolTable *SymbolTable
	LLVMModule  *ir.Module
	funcStack   []*FuncBlock
	ctx         *Context
}

type CompileResult struct {
	Type string
	Val  value.Value
}

type Parameter struct {
	Name string
	Type string
}

type Argument struct {
	Value string
	Type  string
}

// TODO:  Fib(Sub(val)) - this is an interesting case and not sure if its working - depends on how many values are returned
// TODO: Returning boolean expressions for functions with return type bool

// TODO: avoided array as a function parameter for now
// for function calls
var funcMap = make(map[string]*FuncBlock)

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

	LLVMModule := llirgen.LLVMIRModule()
	funcMain := llirgen.LLVMIRFuncMain(LLVMModule)
	mainBlock := llirgen.LLVMIRFunctionBlock(funcMain, "entry")
	ctx := NewContext(mainBlock)

	funcMap["entry"] = &FuncBlock{block: ctx.Block, function: funcMain}

	// c.pushFunctionScope(&FuncBlock{block: mainBlock, func_: funcMain})

	return &Compiler{
		symbolTable: symbolTable,
		// STEP 1: create a module in program header
		LLVMModule: LLVMModule,
		ctx:        ctx,
	}
}

func NewWithState() *Compiler {
	compiler := New()
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

		// STEP 1: create a module in program header
		// var LLVMModule = llirgen.LLVMIRModule()
		// Code gen: Set main entrypoint

	case *ast.ProgramBody:
		for _, decl := range node.Declarations {
			_, err := c.Compile(decl)
			if err != nil {
				return CompileResult{}, err
			}
		}

		print("before program body declarations\n")

		print("after program body declarations\n")

		for _, stmt := range node.Statements {
			_, err := c.Compile(stmt)
			if err != nil {
				return CompileResult{}, err
			}
		}

		// Code gen: end main function, return 0 from main
		function, ok := c.symbolTable.getCurrentFunction()
		if !ok {
			return CompileResult{}, fmt.Errorf("could not retrieve current function")
		}
		if function.Name != "entry" {
			return CompileResult{}, fmt.Errorf("entry function not found")
		}
		currFuncBlock := funcMap[function.Name]
		llirgen.LLVMIRReturn(currFuncBlock.block, constant.NewInt(types.I64, 0))

	case *ast.VariableDeclaration:
		function, ok := c.symbolTable.getCurrentFunction()
		if !ok {
			return CompileResult{}, fmt.Errorf("could not retrieve current function")
		}
		print("function.name: ", function.Name, "\n")
		currFuncBlock := funcMap[function.Name]
		// entryBlock := funcMap["entry"]

		if node.Type.Array != nil {
			print("did you run array - in variable\n")
			// Handle array declaration
			// First, compile the inner variable declaration
			// Then, define the symbol in the symbol table as an array
			if c.symbolTable.IsGlobalScope() {
				c.symbolTable.DefineArray(node.Name.Value, node.Type.Name+"[]", node.Type.Array.Value, GlobalScope)
				global := llirgen.LLVMIRGlobalVariable(c.LLVMModule, node.Name.Value, node.Type.Name+"[]")
				// loadGlobal := entryBlock.block.NewLoad(llirgen.GetLLVMIRType(node.Type.Name+"[]"), global)
				c.ctx.vars[node.Name.Value] = global
			} else {
				symbol := c.symbolTable.DefineArray(node.Name.Value, node.Type.Name+"[]", node.Type.Array.Value, LocalScope)
				alloca := llirgen.LLVMIRAlloca(currFuncBlock.block, node.Name.Value, node.Type.Name+"[]")
				// currFuncBlock.block.NewStore(llirgen.GetLLVMIRConstant(node.Type.Name+"[]"), alloca)
				c.ctx.vars[node.Name.Value] = alloca
				print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case\n")
			}
		} else {
			// Handle variable declaration
			// First, compile the inner variable declaration
			// Then, define the symbol in the symbol table as a variable
			if c.symbolTable.IsGlobalScope() {
				print("you only you\n")
				symbol := c.symbolTable.Define(node.Name.Value, node.Type.Name, false)
				global := llirgen.LLVMIRGlobalVariable(c.LLVMModule, node.Name.Value, node.Type.Name)
				// loadGlobal := entryBlock.block.NewLoad(llirgen.GetLLVMIRType(node.Type.Name), global)
				c.ctx.vars[node.Name.Value] = global

				print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case - 1\n")
			} else {
				symbol := c.symbolTable.Define(node.Name.Value, node.Type.Name, false)
				alloca := llirgen.LLVMIRAlloca(currFuncBlock.block, node.Name.Value, node.Type.Name)
				// currFuncBlock.block.NewStore(llirgen.GetLLVMIRConstant(node.Type.Name), alloca)
				c.ctx.vars[node.Name.Value] = alloca

				print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case - 2\n")
			}
		}

	case *ast.GlobalVariableDeclaration:
		// currFuncBlock := funcMap["entry"]
		// Handle global variable declaration
		// First, compile the inner variable declaration
		// Then, define the symbol in the symbol table as a global variable
		if node.VariableDeclaration.Type.Array != nil {
			print("did you run array - in global\n")
			// check if the array index value is less than or equal to 0
			if node.VariableDeclaration.Type.Array.Value <= 0 {
				return CompileResult{}, fmt.Errorf("array size must be greater than 0")
			}
			if typeofObject(node.VariableDeclaration.Type.Array.Value) != "int" {
				return CompileResult{}, fmt.Errorf("array size must be an integer")
			}
			symbol := c.symbolTable.DefineArray(node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name+"[]", node.VariableDeclaration.Type.Array.Value, GlobalScope)
			global := llirgen.LLVMIRGlobalVariable(c.LLVMModule, node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name+"[]")
			// loadGlobal := currFuncBlock.block.NewLoad(llirgen.GetLLVMIRType(node.VariableDeclaration.Type.Name+"[]"), global)
			c.ctx.vars[node.VariableDeclaration.Name.Value] = global

			print(symbol.Name, symbol.Index, symbol.Scope, "1 - in Global Variable Declaration case\n")
		} else {
			symbol := c.symbolTable.DefineGlobal(node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name)
			global := llirgen.LLVMIRGlobalVariable(c.LLVMModule, node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name)
			// loadGlobal := currFuncBlock.block.NewLoad(llirgen.GetLLVMIRType(node.VariableDeclaration.Type.Name), global)
			c.ctx.vars[node.VariableDeclaration.Name.Value] = global
			print(symbol.Name, symbol.Index, symbol.Scope, "2 - in Global Variable Declaration case\n")
		}

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
		if symbol.Scope == ParamLocalScope {
			// get the local symbols for the current function
			currentFunction := c.ctx.Parent
			params := currentFunction.Params
			for _, param := range params {
				if param.Name() == symbol.Name {
					print(param.String(), " : param name\n")
					return CompileResult{Type: symbol.Type, Val: param}, nil
				}
			}
			print(currentFunction.String(), " : current function param local in Identifier\n")
			// currFuncBlock := funcMap[""]
		}
		print(symbol.Name, symbol.Type, symbol.Index, symbol.Scope, "in Identifier case\n")
		// TODO: might have to change this for function call because they are also identifiers
		val := c.ctx.lookupVariable(symbol.Name)
		print("value vroo: ", val.String(), "\n")
		return CompileResult{Type: symbol.Type, Val: val}, nil

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
		conditionExprString := node.Condition.String()
		if !c.isBooleanExpression(conditionExprString) {
			return CompileResult{}, fmt.Errorf("loop condition must be a boolean expression")
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
		fmt.Printf("Type of curr node in infix - right: %T\n", node.Right)
		cr_, err_ := c.Compile(node.Right)
		if err_ != nil {
			return CompileResult{}, err_
		}

		// check if left and right types match
		if cr.Type != cr_.Type {
			return CompileResult{}, fmt.Errorf("type mismatch: cannot perform operation %s on %s and %s", node.Operator, cr.Type, cr_.Type)
		} else {
			function, ok := c.symbolTable.getCurrentFunction()
			if !ok {
				return CompileResult{}, fmt.Errorf("could not retrieve current function")
			}
			currFuncBlock := funcMap[function.Name]
			print(function.Name, " : function name in Infix Expr\n")
			exprValue := c.CompileInfixExpression(currFuncBlock, node, cr, cr_)
			return CompileResult{Type: cr.Type, Val: exprValue}, nil
		}

	case *ast.IfExpression:
		print(node.Condition, " - condition\n")
		_, err := c.Compile(node.Condition)
		if err != nil {
			return CompileResult{}, err
		}
		conditionExprString := node.Condition.String()
		if !c.isBooleanExpression(conditionExprString) {
			return CompileResult{}, fmt.Errorf("if condition must be a boolean expression")
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
		// currFuncBlock := c.currentFunction()

		cr, err := c.Compile(node.Value)
		if err != nil {
			return CompileResult{}, err
		}

		// If the assignment has an index expression - array indexing - compile it first
		if node.Destination.Expression != nil {
			cr_, err_ := c.Compile(node.Destination.Expression)
			if err != nil {
				return CompileResult{}, err_
			}
			if cr_.Type != cr.Type {
				return CompileResult{}, fmt.Errorf("type mismatch in array assignment: cannot assign %s to %s", cr.Type, cr_.Type)
			}
		} else {
			// Compile the identifier part of the destination
			symbol, ok := c.symbolTable.Resolve(node.Destination.Identifier.Value)
			print(symbol.Name, symbol.Index, symbol.Scope, "in AssignmentStatement case - print for usage\n")
			if !ok {
				return CompileResult{}, fmt.Errorf("variable %s not defined", node.Destination.Identifier.Value)
			}
			if cr.Type != symbol.Type {
				return CompileResult{}, fmt.Errorf("type mismatch: cannot assign %s to %s", cr.Type, symbol.Type)
			}

			// assignment statement codegen
			if _, ok := node.Value.(*ast.Identifier); ok {
				print("yes?\n")
				// alloca := c.ctx.NewAlloca(llirgen.GetLLVMIRType(symbol.Type))
				// print(alloca.Type().String(), " : alloca type\n")
				print("yes-1?\n")
				// store := llirgen.LLVMIRStore(c.ctx.Block, cr.Val, alloca)
				// crVal := c.ctx.NewBitCast(cr.Val, llirgen.GetLLVMIRType(symbol.Type))
				// alloca := c.ctx.vars[node.Destination.Identifier.Value]
				print(cr.Val.String(), symbol.Name, " : cr val-1\n")
				lhsAlloca := c.ctx.vars[symbol.Name]
				rhsVal := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
				c.ctx.NewStore(rhsVal, lhsAlloca)
			} else if _, ok := node.Value.(*ast.CallExpression); ok {
				lhsAlloca := c.ctx.vars[symbol.Name]
				print(lhsAlloca.Type().String(), " herr question?\n")
				print(cr.Val.String(), " : cr val in call expr in assign statement\n")
				// print(cr.Val.Type().String(), " : cr val type haha\n")
				// rhsVal := c.ctx.NewLoad(types.NewPointer(cr.Val.Type()), cr.Val)
				c.ctx.NewStore(cr.Val, lhsAlloca)
				print("yes-2?\n")
			} else {
				print(symbol.Name, " - symbol name\n")
				// this is nil
				lhsAlloca := c.ctx.vars[symbol.Name]
				print(lhsAlloca.String(), " : lhs alloca\n")
				print(cr.Val.String(), " : cr val\n")
				print(cr.Val.Type().String(), " : cr val type\n")
				ptrType := types.NewPointer(cr.Val.Type())
				print(ptrType.String(), " : ptr type\n")
				// rhsVal := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
				c.ctx.NewStore(cr.Val, lhsAlloca)
			}
			print(symbol.Type, cr.Type, "hello symbol type here\n")
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

		print("now: \n")
		function, ok := c.symbolTable.getCurrentFunction()
		print(function.Name, function.ReturnType, "in return statement printing function stuff\n")
		if !ok {
			return CompileResult{}, fmt.Errorf("return statement outside of function")
		}
		if cr.Type != function.ReturnType {
			if typesCompatible(cr.Type, function.ReturnType) {
				return CompileResult{Type: function.ReturnType}, nil
			}
			return CompileResult{}, fmt.Errorf("type mismatch for function %s: cannot return %s from function of type %s", function.Name, cr.Type, function.ReturnType)
		}

		// Code gen: Return statement
		// currFuncBlock := c.currentFunction()
		c.ctx.NewRet(cr.Val)
		// currFuncBlock.block.NewRet(llirgen.LLVMIRGlobalVariable(c.LLVMModule, node.ReturnValue.String(), 0))

		return CompileResult{Type: cr.Type}, nil

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		print(str, " : haha i in string literal\n")
		return CompileResult{Type: string(str.Type())}, nil

	case *ast.IntegerLiteral:
		fmt.Printf("Type of curr node in integer literal: %T\n", node.Value)
		integer := &object.Integer{Value: node.Value}
		print(constant.NewInt(types.I64, node.Value).String(), " : integer literal\n")
		return CompileResult{Type: string(integer.Type()), Val: constant.NewInt(types.I64, node.Value)}, nil

	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		return CompileResult{Type: string(float.Type()), Val: constant.NewFloat(types.Float, node.Value)}, nil

	case *ast.Boolean:
		boolean := &object.Boolean{Value: node.Value}
		if node.Value {
			return CompileResult{Type: string(boolean.Type()), Val: constant.NewInt(types.I1, 1)}, nil
		} else {
			return CompileResult{Type: string(boolean.Type()), Val: constant.NewInt(types.I1, 0)}, nil
		}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			_, err := c.Compile(el)
			if err != nil {
				return CompileResult{}, err
			}
		}

	case *ast.IndexExpression:
		print("ran index\n")
		// left is the identifier
		cr, err := c.Compile(node.Left)
		if err != nil {
			return CompileResult{}, err
		}
		// right is the index ryan[3] - where right is 3
		cr_, err_ := c.Compile(node.Index)
		if err_ != nil {
			return CompileResult{}, err
		}
		if cr_.Type != "integer" {
			return CompileResult{}, fmt.Errorf("index must be an integer")
		}
		// also check if the index passed is less than the size of the array
		symbol, ok := c.symbolTable.Resolve(node.Left.String())
		if !ok {
			return CompileResult{}, fmt.Errorf("undefined variable %s", node.Left.String())
		} else {
			if symbol.Type != "integer[]" {
				return CompileResult{}, fmt.Errorf("variable %s is not an array", node.Left.String())
			} else {
				// check if the index is within the bounds of the array
				int1, err := strconv.ParseInt(node.Index.String(), 6, 12)
				if err != nil {
					return CompileResult{}, fmt.Errorf("there was an error parsing the index to int: %w", err)
				}
				if int1 >= symbol.ArraySize {
					return CompileResult{}, fmt.Errorf("index out of bounds")
				}
			}
		}
		sub := cr.Type[0 : len(cr.Type)-2] // remove the [] from the type string
		return CompileResult{Type: sub}, nil

	case *ast.CallExpression:
		fmt.Printf("Type of curr node in call expression: %T\n", node.Function)

		if _, ok := node.Function.(*ast.Identifier); ok {
			// node.Function is of type *ast.Identifier
			currentFuncName := node.Function.String()
			if builtinWithExists(currentFuncName) {
				c.compileBuiltInFunction(node)
				// TODO: uncomment the below after installing llvm
				return c.insertRuntimeFunctions(node)
			} else {
				// Try to resolve the function name in inner scopes
				symbol, ok := c.symbolTable.ResolveInner(currentFuncName)
				print("haha here i am\n")
				if ok {
					// functon found in inner scope - might be nested procedures
					_, err := c.checkArguments(node)
					if err != nil {
						return CompileResult{}, err
					}
					// funcBlock := funcMap["entry"]
					funcBlock := c.ctx.Block
					funcBlockCaller := funcMap[currentFuncName]
					fnArgs := make([]value.Value, len(node.Arguments))
					for i, arg := range node.Arguments {
						crExpVal, err := c.Compile(arg)
						if err != nil {
							return CompileResult{}, fmt.Errorf("error compiling argument %d: %w", i, err)
						}
						fnArgs[i] = crExpVal.Val
					}
					// callInst := funcBlock.block.NewCall(funcBlockCaller.function, fnArgs...)
					callInst := funcBlock.NewCall(funcBlockCaller.function, fnArgs...)
					print(callInst.String(), " : callInst\n")
					// retValue := funcBlock.block.NewExtractValue(callInst, 0)
					return CompileResult{Type: symbol.Type, Val: callInst}, nil
				} else {
					// function not found in inner scopes - might be a call expression in program body
					symbol, ok := c.symbolTable.Resolve(currentFuncName)
					if !ok {
						return CompileResult{}, fmt.Errorf("undefined function %s", currentFuncName)
					} else {
						// function found in outer scope
						_, err := c.checkArguments(node)
						if err != nil {
							return CompileResult{}, err
						}
						// funcBlock := funcMap["entry"]
						// funcBlockCaller := funcMap[currentFuncName]
						// fnArgs := make([]value.Value, len(node.Arguments))
						// for i, arg := range node.Arguments {
						// 	crExpVal, err := c.Compile(arg)
						// 	if err != nil {
						// 		return CompileResult{}, fmt.Errorf("error compiling argument %d: %w", i, err)
						// 	}
						// 	fnArgs[i] = crExpVal.Val
						// }
						// callInst := funcBlock.block.NewCall(funcBlockCaller.function, fnArgs...)
						// print(callInst.String(), " : callInst\n")
						// retValue := funcBlock.block.NewExtractValue(callInst, 0)
						return CompileResult{Type: symbol.Type}, nil
					}
				}
			}
		} else {
			// node.Function is not of type *ast.Identifier
			_, err := c.Compile(node.Function)
			if err != nil {
				return CompileResult{}, err
			}
			_, err_ := c.checkArguments(node)
			if err_ != nil {
				return CompileResult{}, err_
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

	case *ast.ProcedureHeader:
		c.enterScope()
		print("XXX")
		// define the function name and parameters in the symbol table
		if node.Name.Value != "" {
			c.symbolTable.DefineFunctionName(node.Name.Value, node.TypeMark.Name)
		}

		for _, param := range node.Parameters {
			c.symbolTable.Define(param.Name.Value, param.Type.Name, true)
		}

		// Code gen: function
		funcDef := llirgen.LLVMIRFunctionDefinition(c.LLVMModule, node.Name.Value, node.TypeMark.Name, node.Parameters)
		// f := c.ctx.Parent
		newCtx := c.ctx.NewContext(funcDef.NewBlock(node.Name.Value))
		c.ctx = newCtx
		funcMap[node.Name.Value] = &FuncBlock{block: newCtx.Block, function: funcDef}

		// print(funcDef.String() + " :funcDef\n")
		print("in procedure header after enter scope\n")
		PrintSymbolTable(c.symbolTable)

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
		// leave scope after the body of the procedure
		// f := c.ctx.Parent
		// leaveB := f.NewBlock("leave.func." + f.Name())
		// TODO: very iffy
		// c.ctx.leaveBlock = leaveB
		// c.ctx.NewBr(c.ctx.parent.Block)
		c.ctx = c.ctx.parent
		c.leaveScope()
		print("after procedure body leave scope\n")
		PrintSymbolTable(c.symbolTable)
	}

	return CompileResult{}, nil

}

func (c *Compiler) enterScope() {
	if c.symbolTable.Inner == nil {
		c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
	} else {
		c.symbolTable = c.symbolTable.Inner
	}
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
	print("\nafter getParams\n")
	functionScope, funcIndex := findFunctionScope(symbolTable, functionName)
	if functionScope == nil {
		// Function not found, return empty slice
		return []Symbol{}
	}

	localSymbols := make([]Symbol, 0)
	for _, sym := range functionScope.store {
		if sym.Scope == ParamLocalScope && sym.Index == funcIndex {
			print(sym.Name + " haha\n")
			localSymbols = append(localSymbols, sym)
		}
	}
	sortParamLocalSymbols(localSymbols)

	return localSymbols
}

// Find the symbol table containing the function definition
func findFunctionScope(symbolTable *SymbolTable, functionName string) (*SymbolTable, int) {
	current := symbolTable
	print(functionName)
	print("\nprinting current symbol table to detect loop\n")
	PrintSymbolTable(current)
	for current != nil {
		// Check if the current symbol table contains the function definition
		if _, ok := current.store[functionName]; ok {
			print(current.store[functionName].Name, " - function name\n")
			return current, current.store[functionName].Index
		}
		// Move to the inner symbol table
		current = current.Inner
	}
	// Function scope not found
	return nil, -1
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

func (c *Compiler) insertRuntimeFunctions(node *ast.CallExpression) (CompileResult, error) {
	// Insert runtime functions
	// putinteger
	// putfloat
	// putstring
	// putbool
	// sqrt
	// getinteger
	// getfloat
	// getstring
	// getbool
	m := c.LLVMModule
	switch node.Function.String() {
	case "putinteger":
		putinteger := m.NewFunc("putinteger", types.I1)
		putinteger.Params = append(putinteger.Params, ir.NewParam("paramValue", types.I64Ptr))
		putintegerEntry := putinteger.NewBlock("putinteger.entry")
		loaded := putintegerEntry.NewLoad(types.I64, putinteger.Params[0])
		formatStrGlobal := m.NewGlobalDef(".textstr", constant.NewCharArrayFromString("%d\n\x00"))
		// Define indices for the getelementptr instruction.
		indices := []value.Value{
			constant.NewInt(types.I64, 0), // Index 0 to access the first element.
			constant.NewInt(types.I64, 0), // Index 0 again, since it's a flat array.
		}
		formatStrPtr := putintegerEntry.NewGetElementPtr(formatStrGlobal.Typ.ElemType, formatStrGlobal, indices...)
		cr, err := c.Compile(node.Arguments[0])
		if err != nil {
			return CompileResult{}, fmt.Errorf("error compiling argument for putinteger: %w", err)
		}
		putintegerEntry.NewCall(
			m.NewFunc("printf", types.I32, ir.NewParam("format", types.NewPointer(types.I8))),
			formatStrPtr,
			loaded, // Pass the integer argument to printf.
		)
		putintegerEntry.NewRet(constant.NewInt(types.I1, 1))

		callInstFromMain := funcMap["entry"].block.NewCall(putinteger, cr.Val)
		return CompileResult{Type: "bool", Val: callInstFromMain}, nil
	case "getbool":
		getbool := m.NewFunc("getbool", llirgen.GetLLVMIRType("bool"))
		getboolEntry := getbool.NewBlock("getbool.entry")
		boolVal := constant.NewInt(types.I1, 1)
		getboolEntry.NewRet(constant.NewInt(types.I1, 1))
		return CompileResult{Type: "bool", Val: boolVal}, nil
	default:
		return CompileResult{}, nil
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

func (c *Compiler) checkArguments(node *ast.CallExpression) (CompileResult, error) {
	// Get local symbols for the function being called
	paramLocalSymbols := getParamLocalSymbols(c.symbolTable, node.Function.String())
	print(len(paramLocalSymbols), " - paramLocalSymbols\n")
	print(node.Function.String(), "here in checkArgs\n")
	print(len(node.Arguments), " - len of arguments\n")
	// Check if the number of arguments matches the number of local symbols
	if len(node.Arguments) < len(paramLocalSymbols) {
		return CompileResult{}, fmt.Errorf("not enough arguments provided for function call")
	} else if len(node.Arguments) > len(paramLocalSymbols) {
		return CompileResult{}, fmt.Errorf("too many arguments provided for function call")
	}

	// Check the type of each argument against the corresponding local symbol
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

	print(node.Arguments[0].String(), " - arguments lmao\n")

	// All argument types match, return success
	return CompileResult{Type: "success"}, nil
}

func typeofObject(variable interface{}) string {
	switch variable.(type) {
	case int, int64:
		return "int"
	case float32:
		return "float32"
	case bool:
		return "boolean"
	case string:
		return "string"
	default:
		return "unknown"
	}
}

// isBooleanExpression checks if the given expression is a boolean expression.
// It returns true if the expression is boolean, otherwise false.
func (c *Compiler) isBooleanExpression(expr string) bool {
	print(expr, " - expr\n")
	if expr == "true" || expr == "false" {
		return true
	}
	// List of relational operators
	relationalOperators := map[string]bool{
		"<":  true,
		">":  true,
		"<=": true,
		">=": true,
		"==": true,
		"!=": true,
	}

	// Check if expr is an identifier representing a boolean variable
	if sym, ok := c.symbolTable.store[expr]; ok && sym.Type == "bool" {
		return true
	}

	// tokenize the string
	tokens := strings.Split(expr, " ")

	// Check if a character is a relational operator in the expr string
	for _, token := range tokens {
		if relationalOperators[token] {
			return true
		}
	}

	// If no relational operators found, it's not a boolean expression
	return false
}

func (c *Compiler) CompileInfixExpression(funcBlock *FuncBlock, node *ast.InfixExpression, cr CompileResult, cr_ CompileResult) value.Value {
	switch node.Operator {
	case "+":
		return funcBlock.block.NewAdd(cr.Val, cr_.Val)
	case "-":
		return funcBlock.block.NewSub(cr.Val, cr_.Val)
	case "*":
		if _, ok := cr.Val.(*ir.InstAlloca); ok {
			// Load the value if cr.Val is an *ir.InstAlloca
			val := funcBlock.block.NewLoad(types.I64, cr.Val)
			return funcBlock.block.NewMul(val, cr_.Val)
		} else if _, ok := cr_.Val.(*ir.InstAlloca); ok {
			// Load the value if cr_.Val is an *ir.InstAlloca
			val := funcBlock.block.NewLoad(types.I64, cr_.Val)
			return funcBlock.block.NewMul(cr.Val, val)
		}
		return funcBlock.block.NewMul(cr.Val, cr_.Val)
	case "/":
		return funcBlock.block.NewSDiv(cr.Val, cr_.Val)
	case "<":
		return funcBlock.block.NewICmp(enum.IPredSLT, cr.Val, cr_.Val)
	case ">":
		return funcBlock.block.NewICmp(enum.IPredSGT, cr.Val, cr_.Val)
	case "<=":
		return funcBlock.block.NewICmp(enum.IPredSLE, cr.Val, cr_.Val)
	case ">=":
		return funcBlock.block.NewICmp(enum.IPredSGE, cr.Val, cr_.Val)
	case "==":
		return funcBlock.block.NewICmp(enum.IPredEQ, cr.Val, cr_.Val)
	case "!=":
		return funcBlock.block.NewICmp(enum.IPredNE, cr.Val, cr_.Val)
	default:
		panic("Unimplemented infix expression")
	}
}
