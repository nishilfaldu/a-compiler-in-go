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

func (ctx *Context) HasTerminator() bool {
	// Get the last instruction in the current basic block
	// lastInstr := ctx.Block.Insts[len(ctx.Block.Insts)-1]

	// Check if the last instruction is a terminator instruction
	// print(lastInstr.LLString(), " : last instruction\n")
	// print(ctx.Block.Term.LLString(), " : last instruction\n")
	print(ctx.Block.Term != nil, " woppsie\n")
	return ctx.Block.Term != nil
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
	ctx         *Context
	ifCounter   int
	argCounter  int
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
		ifCounter:  0,
		argCounter: 0,
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
		c.InsertPutStringAtRuntime() // Insert "putstring" string at runtime
		c.DeclareStrCmp()

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
		if !c.ctx.HasTerminator() {
			c.ctx.NewRet(constant.NewInt(types.I64, 0))
		} else {
			c.ctx.leaveBlock.NewRet(constant.NewInt(types.I64, 0))
		}

		// c.ctx.NewRet(constant.NewInt(types.I64, 0))
		// print(c.ctx.Parent.LLString(), " : block name\n")

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
				global := c.LLVMModule.NewGlobalDef(node.Name.Value, constant.NewArray(types.NewArray(uint64(node.Type.Array.Value), llirgen.GetLLVMIRType(node.Type.Name))))
				if node.Type.Name == "integer" {
					global.Init = constant.NewZeroInitializer(types.NewArray(uint64(node.Type.Array.Value), llirgen.GetLLVMIRType(node.Type.Name)))
				}
				c.ctx.vars[node.Name.Value] = global
			} else {
				symbol := c.symbolTable.DefineArray(node.Name.Value, node.Type.Name+"[]", node.Type.Array.Value, LocalScope)
				alloca := currFuncBlock.block.NewAlloca(types.NewArray(uint64(node.Type.Array.Value), llirgen.GetLLVMIRType(node.Type.Name)))
				alloca.SetName(node.Name.Value)
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
				c.ctx.vars[node.Name.Value] = global

				print(symbol.Name, symbol.Index, symbol.Scope, "in Variable Declaration case - 1\n")
			} else {
				symbol := c.symbolTable.Define(node.Name.Value, node.Type.Name, false)
				print(c.ctx.Block.Name(), " : block name in var dac\n")
				print(currFuncBlock.block.Name(), " : block in var dac\n")
				// changed currFuncBlock to c.ctx.Block
				alloca := llirgen.LLVMIRAlloca(c.ctx.Block, node.Name.Value, node.Type.Name)
				c.ctx.vars[node.Name.Value] = alloca
				print(alloca.LLString(), " : alloca in var decl\n")
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
			global := c.LLVMModule.NewGlobalDef(node.VariableDeclaration.Name.Value, constant.NewArray(types.NewArray(uint64(node.VariableDeclaration.Type.Array.Value), llirgen.GetLLVMIRType(node.VariableDeclaration.Type.Name))))
			if node.VariableDeclaration.Type.Name == "integer" {
				global.Init = constant.NewZeroInitializer(types.NewArray(uint64(node.VariableDeclaration.Type.Array.Value), llirgen.GetLLVMIRType(node.VariableDeclaration.Type.Name)))
			}
			c.ctx.vars[node.VariableDeclaration.Name.Value] = global
			print(symbol.Name, symbol.Index, symbol.Scope, "1 - in Global Variable Declaration case\n")
		} else {
			symbol := c.symbolTable.DefineGlobal(node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name)
			global := llirgen.LLVMIRGlobalVariable(c.LLVMModule, node.VariableDeclaration.Name.Value, node.VariableDeclaration.Type.Name)
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
			// currentFunction := c.ctx.Parent
			// params := currentFunction.Params
			// for _, param := range params {
			// 	if param.Name() == symbol.Name {
			// 		print(param.String(), " : param name\n")
			// 		return CompileResult{Type: symbol.Type, Val: param}, nil
			// 	}
			// }
			// print(currentFunction.String(), " : current function param local in Identifier\n")
			return c.FindIdentifierInScopes(symbol)
		}
		print(symbol.Name, symbol.Type, symbol.Index, symbol.Scope, "in Identifier case\n")
		// TODO: might have to change this for function call because they are also identifiers
		val := c.ctx.lookupVariable(symbol.Name)
		print("value vroo: ", val.String(), "\n")
		return CompileResult{Type: symbol.Type, Val: val}, nil

	case *ast.LoopStatement:
		// Compile the initialization statement
		currentFunction := c.ctx.Parent
		print(currentFunction.Name() + " is parent of loop stmt \n")
		loopCondCtx := c.ctx.NewContext(currentFunction.NewBlock("for.cond"))
		loopBodyCtx := c.ctx.NewContext(currentFunction.NewBlock("for.body"))
		// endBodyCtx := c.ctx.NewContext(currentFunction.NewBlock("for.end"))
		leaveForBlock := currentFunction.NewBlock("leave.for.loop")

		// initName := node.InitStatement.Destination.Identifier.Value
		_, err := c.Compile(node.InitStatement)
		if err != nil {
			return CompileResult{}, fmt.Errorf("error compiling loop initialization expression: %w", err)
		}

		// add a break in the entry body
		x := c.ctx.NewBr(loopCondCtx.Block)
		print(x.LLString(), " : new br\n")

		// Compile the loop condition
		c.ctx = loopCondCtx
		crCond, err_ := c.Compile(node.Condition)
		if err_ != nil {
			return CompileResult{}, fmt.Errorf("error compiling loop condition: %w", err_)
		}
		conditionExprString := node.Condition.String()
		if !c.isBooleanExpression(conditionExprString) {
			return CompileResult{}, fmt.Errorf("loop condition must be a boolean expression")
		}

		loopCondCtx.NewCondBr(crCond.Val, loopBodyCtx.Block, leaveForBlock)
		c.ctx = c.ctx.parent
		// endBodyCtx.NewRet(nil)

		// Compile the loop body
		c.ctx = loopBodyCtx
		_, err = c.Compile(node.Body)
		if err != nil {
			return CompileResult{}, err
		}
		loopBodyCtx.NewBr(loopCondCtx.Block)
		c.ctx = c.ctx.parent
		c.ctx.leaveBlock = leaveForBlock // set this to allow breaks from within the loop body

	case *ast.ForBlockStatement:
		print("i did run okay in for block?\n")
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
		// function, ok := c.symbolTable.getCurrentFunction()
		// if !ok {
		// 	return CompileResult{}, fmt.Errorf("could not retrieve current function")
		// }
		// currFuncBlock := funcMap[function.Name]
		exprValue := c.CompilePrefixExpression(c.ctx.Block, node, cr)
		// print(c.ctx.LLString(), " : ctx in prefix expr\n")
		return CompileResult{Type: cr.Type, Val: exprValue}, nil

	case *ast.InfixExpression:
		print(node.Operator, " - operator\n")
		// if node.Operator == "<" {
		// 	_, err := c.Compile(node.Right)
		// 	if err != nil {
		// 		return CompileResult{}, err
		// 	}

		// 	_, err_ := c.Compile(node.Left)
		// 	if err_ != nil {
		// 		return CompileResult{}, err_
		// 	}
		// }
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
			// currFuncBlock := funcMap[function.Name]
			print(function.Name, " : function name in Infix Expr\n")
			print(cr.Val.String(), cr_.Val.String(), " : cr val in Infix Expr\n")
			var exprValue value.Value
			if c.ctx.leaveBlock != nil {
				exprValue = c.CompileInfixExpression(c.ctx.leaveBlock, node, cr, cr_)
			} else {
				exprValue = c.CompileInfixExpression(c.ctx.Block, node, cr, cr_)
			}
			return CompileResult{Type: cr.Type, Val: exprValue}, nil
		}

	case *ast.IfExpression:
		currentFunction := c.ctx.Parent
		print(currentFunction.Name() + " is parent of if expr \n")
		thenCtx := c.ctx.NewContext(currentFunction.NewBlock("if.then" + strconv.Itoa(c.ifCounter)))
		elseCtx := c.ctx.NewContext(currentFunction.NewBlock("if.else" + strconv.Itoa(c.ifCounter)))
		leaveIfBlock := currentFunction.NewBlock("leave.if" + strconv.Itoa(c.ifCounter))
		// leaveCtx := c.ctx.NewContext(leaveIfBlock)

		crCond, err := c.Compile(node.Condition)
		if err != nil {
			return CompileResult{}, err
		}
		conditionExprString := node.Condition.String()
		if !c.isBooleanExpression(conditionExprString) {
			return CompileResult{}, fmt.Errorf("if condition must be a boolean expression")
		}

		if crCondVal, ok := crCond.Val.(*ir.Global); ok {
			loaded := c.ctx.NewLoad(llirgen.GetLLVMIRType(crCond.Type), crCondVal)
			crCond.Val = loaded
		}

		if c.ctx.leaveBlock != nil {
			c.ctx.leaveBlock.NewCondBr(crCond.Val, thenCtx.Block, elseCtx.Block)
		} else {
			c.ctx.NewCondBr(crCond.Val, thenCtx.Block, elseCtx.Block)
		}
		print(c.ctx.LLString(), " : ctx in if expr\n")

		c.ctx = thenCtx
		_, err_ := c.Compile(node.Consequence)
		if err_ != nil {
			return CompileResult{}, err
		}
		if !thenCtx.HasTerminator() {
			thenCtx.NewBr(leaveIfBlock)
		}
		// thenCtx.NewBr(leaveCtx.Block)
		c.ctx = c.ctx.parent
		c.ctx.leaveBlock = leaveIfBlock

		// thenCtx.HasTerminator()
		// else block
		if node.Alternative != nil {
			c.ctx = elseCtx
			_, err := c.Compile(node.Alternative)
			if err != nil {
				return CompileResult{}, err
			}
			elseCtx.NewBr(leaveIfBlock)
			c.ctx = c.ctx.parent
		}

		if node.Alternative == nil {
			if !elseCtx.HasTerminator() {
				elseCtx.NewBr(leaveIfBlock)
			}
		}
		// c.ctx = leaveCtx
		c.ifCounter += 1
		print(c.ctx.Parent.Name(), " : parent name\n")
		return CompileResult{}, nil

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
		print(cr.Type, " : cr type in AssignmentStatement case\n")
		print(node.Value.String(), " : node value in AssignmentStatement case\n")
		if err != nil {
			return CompileResult{}, err
		}

		// If the assignment has an index expression - array indexing - compile it first
		if node.Destination.Expression != nil {
			cr_, err_ := c.Compile(node.Destination)
			// print(node.Destination.String(), "\n")
			if err != nil {
				return CompileResult{}, err_
			}
			if cr_.Type != cr.Type {
				return CompileResult{}, fmt.Errorf("type mismatch in array assignment: cannot assign %s to %s", cr.Type, cr_.Type)
			}
			if _, ok := node.Value.(*ast.Identifier); ok {
				// where right is an identifier - not sure if it will work for all variables
				loaded := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
				c.ctx.NewStore(loaded, cr_.Val)
			} else {
				// array assignment where right is also an array indexing operation
				c.ctx.Block.NewStore(cr.Val, cr_.Val)
			}
		} else {
			// // Compile the identifier part of the destination
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
				print("yes-1?\n")
				print(cr.Val.String(), symbol.Name, " : cr val-1\n")
				lhsAlloca := c.ctx.lookupVariable(symbol.Name)
				rhsVal := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
				c.ctx.NewStore(rhsVal, lhsAlloca)
			} else if _, ok := node.Value.(*ast.IndexExpression); ok {
				print("how right i am?\n")
			} else if _, ok := node.Value.(*ast.CallExpression); ok {
				if symbol.Scope == ParamLocalScope {
					// get the local symbols for the current function
					currentFunction := c.ctx.Parent
					params := currentFunction.Params
					for _, param := range params {
						if param.Name() == symbol.Name {
							// might
							print(param.String(), " : param string\n")
							load := c.ctx.NewLoad(param.Type(), param)
							print(load.LLString(), "\n")
							print(cr.Val.String(), " : cr val in Identifier\n")
							// loadedCrVal := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
							c.ctx.NewStore(cr.Val, param)
						}
					}
					print(currentFunction.String(), " : current function param local in Identifier\n")
				} else {
					lhsAlloca := c.ctx.lookupVariable(symbol.Name)
					print(lhsAlloca.Type().String(), " lhs alloca herr question?\n")
					print(cr.Val.String(), " : cr val in call expr in assign statement\n")
					if c.ctx.leaveBlock != nil {
						c.ctx.leaveBlock.NewStore(cr.Val, lhsAlloca)
					} else {
						c.ctx.NewStore(cr.Val, lhsAlloca)
					}
				}
			} else if _, ok := node.Value.(*ast.IntegerLiteral); ok {
				lhsAlloca := c.ctx.lookupVariable(symbol.Name)
				c.ctx.NewStore(cr.Val, lhsAlloca)
			} else {
				print(symbol.Name, " - symbol name\n")
				if symbol.Scope == ParamLocalScope {
					if symbol.Scope == ParamLocalScope {
						// get the local symbols for the current function
						currentFunction := c.ctx.Parent
						params := currentFunction.Params
						for _, param := range params {
							if param.Name() == symbol.Name {
								// might
								print(param.String(), " : param string\n")
								load := c.ctx.NewLoad(param.Type(), param)
								print(load.LLString(), "\n")
								print(cr.Val.String(), " : cr val in Identifier\n")
								// loadedCrVal := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
								c.ctx.NewStore(cr.Val, param)
							}
						}
						print(currentFunction.String(), " : current function param local in Identifier\n")
					}
				} else {
					if symbol.Scope == GlobalScope && symbol.Type == "string" {
						// TODO: fix this
						global := c.ctx.lookupVariable(symbol.Name)
						if global_, ok := global.(*ir.Global); ok {
							print(node.Value.String(), " : node value\n")
							print(cr.Val.Type().String(), " : cr val type - haah\n")
							ptrConst := constant.NewGetElementPtr(cr.Val.Type(), c.LLVMModule.NewGlobalDef(symbol.Name+"_0", constant.NewCharArrayFromString(node.Value.String())), constant.NewInt(types.I64, 0))
							global_.Init = ptrConst
						}
					} else {
						lhsAlloca := c.ctx.lookupVariable(symbol.Name)
						if c.ctx.leaveBlock != nil {
							c.ctx.leaveBlock.NewStore(cr.Val, lhsAlloca)
						} else {
							c.ctx.NewStore(cr.Val, lhsAlloca)
						}
					}
				}
			}
			print(symbol.Type, cr.Type, "hello symbol type here\n")
		}

	case *ast.Destination:
		// Compile the identifier part of the destination
		print(node.String(), "here in destination")
		// here arr[idx] - idx is an expression
		cr_, err_ := c.Compile(node.Expression)
		if err_ != nil {
			return CompileResult{}, err_
		}
		if cr_.Type != "integer" {
			return CompileResult{}, fmt.Errorf("index must be an integer")
		}

		print(cr_.Val.String(), " : cr val in destination\n")

		symbol, ok := c.symbolTable.Resolve(node.Identifier.String())
		subTyp := symbol.Type[0 : len(symbol.Type)-2]
		if !ok {
			return CompileResult{}, fmt.Errorf("undefined array variable %s", node.Identifier.String())
		} else {
			if symbol.Type != "integer[]" {
				return CompileResult{}, fmt.Errorf("variable %s is not an array", node.Identifier.String())
			} else {
				// check if the index is within the bounds of the array
				int1, err := strconv.ParseInt(node.Expression.String(), 6, 12)
				if err != nil {
					return CompileResult{}, fmt.Errorf("there was an error parsing the index to int: %w", err)
				}
				if int1 >= symbol.ArraySize {
					return CompileResult{}, fmt.Errorf("index out of bounds")
				}
			}
		}
		block := c.ctx.Block
		arraySize := symbol.ArraySize
		arrayIdx, err := strconv.Atoi(node.Expression.String())
		if err != nil {
			return CompileResult{}, fmt.Errorf("error converting index to int: %w", err)
		}
		arrTyp := types.NewArray(uint64(arraySize), llirgen.GetLLVMIRType(subTyp))
		definedArray := c.ctx.lookupVariable(node.Identifier.String())
		pToElem := block.NewGetElementPtr(arrTyp, definedArray, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, int64(arrayIdx)))

		return CompileResult{Type: subTyp, Val: pToElem}, nil

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
				if c.ctx.leaveBlock != nil {
					// If there's a terminator, return from the leave block
					c.ctx.leaveBlock.NewRet(cr.Val)
					return CompileResult{Type: function.ReturnType}, nil
				}
				// return CompileResult{Type: function.ReturnType}, nil
			}
			return CompileResult{}, fmt.Errorf("type mismatch for function %s: cannot return %s from function of type %s", function.Name, cr.Type, function.ReturnType)
		}
		// Code gen: Return statement
		// Check if there's already a terminator
		if c.ctx.leaveBlock != nil {
			// If there's a terminator, return from the leave block
			if cr.Val.Type().Equal(types.NewPointer(llirgen.GetLLVMIRType(cr.Type))) {
				crVal := c.ctx.leaveBlock.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
				c.ctx.leaveBlock.NewRet(crVal)
			} else {
				c.ctx.leaveBlock.NewRet(cr.Val)
			}
			return CompileResult{Type: cr.Type}, nil
		}

		// if !c.ctx.HasTerminator() {
		// 	exitBlock := c.ctx.Parent.NewBlock(c.ctx.Parent.Name() + ".exit")
		// 	c.ctx.NewBr(exitBlock) // Branch to the exit block if no terminator exists
		// 	exitBlock.NewRet(cr.Val)
		// 	return CompileResult{Type: cr.Type}, nil
		// }

		// TODO: consider deleting later
		if cr.Val.Type().Equal(types.NewPointer(llirgen.GetLLVMIRType(cr.Type))) {
			crVal := c.ctx.NewLoad(llirgen.GetLLVMIRType(cr.Type), cr.Val)
			c.ctx.NewRet(crVal)
		} else {
			c.ctx.NewRet(cr.Val)
		}
		return CompileResult{Type: cr.Type}, nil

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		print(node.Value, " : haha i in string literal\n")
		m := c.LLVMModule

		// Create a global definition for the string literal with null termination
		strLit := m.NewGlobalDef("", constant.NewCharArrayFromString(node.Value+"\x00"))
		strPtr := c.ctx.Block.NewGetElementPtr(
			types.NewArray(uint64(len(node.Value)+1), types.I8), strLit, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))

		return CompileResult{Type: string(str.Type()), Val: strPtr}, nil
		// return CompileResult{Type: string(str.Type()), Val: constant.NewCharArrayFromString(node.Value)}, nil

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
		subTyp := cr.Type[0 : len(cr.Type)-2]
		print(cr.Val.String(), cr_.Val.String(), " : cr val in index expression\n")
		block := c.ctx.Block
		arrayIdx, err := strconv.Atoi(node.Index.String())
		if err != nil {
			return CompileResult{}, fmt.Errorf("error converting index to int: %w", err)
		}
		arraySize := symbol.ArraySize
		arrTyp := types.NewArray(uint64(arraySize), llirgen.GetLLVMIRType(subTyp))
		definedArray := c.ctx.lookupVariable(node.Left.String())
		print(definedArray.String(), " : defined array\n")
		print(arrTyp.LLString(), " : arrTyp\n")
		pToElem := block.NewGetElementPtr(arrTyp, definedArray, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, int64(arrayIdx)))
		print(pToElem.LLString(), " : gep\n")
		// loadedArray := block.NewLoad(arrTyp, definedArray)
		// print(loadedArray.String(), " : loaded array\n")
		// val := block.NewExtractValue(loadedArray, uint64(arrayIdx))
		// print(val.String(), " : extracted val in index expression\n")
		// block.NewStore(constant.NewInt(types.I64, 5), pToElem)
		// remove the [] from the type string
		loadedValue := block.NewLoad(llirgen.GetLLVMIRType(subTyp), pToElem)
		print(loadedValue.String(), " : loaded value\n")
		return CompileResult{Type: subTyp, Val: loadedValue}, nil

	case *ast.CallExpression:
		fmt.Printf("Type of curr node in call expression: %T\n", node.Function)

		if _, ok := node.Function.(*ast.Identifier); ok {
			// node.Function is of type *ast.Identifier
			currentFuncName := node.Function.String()
			if builtinWithExists(currentFuncName) {
				c.compileBuiltInFunction(node)
				// TODO: uncomment the below after installing llvm
				if c.ctx.leaveBlock != nil {
					return c.insertRuntimeFunctions(node, c.ctx.leaveBlock)
				} else {
					return c.insertRuntimeFunctions(node, c.ctx.Block)
				}
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
					var funcBlock *ir.Block
					if c.ctx.leaveBlock != nil {
						funcBlock = c.ctx.leaveBlock
					} else {
						funcBlock = c.ctx.Block
					}
					funcBlockCaller := funcMap[currentFuncName]
					fnArgs := make([]value.Value, len(node.Arguments))
					for i, arg := range node.Arguments {
						crExpVal, err := c.Compile(arg)
						if err != nil {
							return CompileResult{}, fmt.Errorf("error compiling argument %d: %w", i, err)
						}
						if crExpVal.Val.Type().Equal(types.NewPointer(llirgen.GetLLVMIRType(crExpVal.Type))) {
							print("i was right\n")
							fnArgs[i] = crExpVal.Val
						} else {
							newAlloca := funcBlock.NewAlloca(llirgen.GetLLVMIRType(crExpVal.Type))
							// I changed this to have a global counter for uniqueness
							// newAlloca.SetName("arg" + strconv.Itoa(i))
							newAlloca.SetName("arg" + strconv.Itoa(c.argCounter))
							funcBlock.NewStore(crExpVal.Val, newAlloca)
							fnArgs[i] = newAlloca
							c.argCounter += 1
						}
					}
					callInst := funcBlock.NewCall(funcBlockCaller.function, fnArgs...)
					print(callInst.String(), " ", callInst.LLString(), " : callInst\n")
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

func (c *Compiler) FindIdentifierInScopes(symbol Symbol) (CompileResult, error) {
	// Start from the current context (function, block, etc.)
	currentContext := c.ctx

	// Traverse upwards through parent scopes until you find the symbol or reach the global scope
	for currentContext != nil {
		// Check if the current context is a function with parameters
		if currentContext.Parent != nil {
			// Search parameters
			currentFunction := currentContext.Parent
			for _, param := range currentFunction.Params {
				if param.Name() == symbol.Name {
					// Found the symbol as a parameter
					print(param.String(), " : param name\n")
					return CompileResult{Type: symbol.Type, Val: param}, nil
				}
			}
		}

		// If not found, check local variables (if applicable) - TODO: could be needed in future
		// if val, exists := currentContext.vars[symbol.Name]; exists {
		// 	// Found the symbol as a local variable
		// 	return CompileResult{Type: symbol.Type, Val: val}, nil
		// }

		// Move to the parent context
		currentContext = currentContext.parent
	}

	// If reached here, symbol is not found in any context
	print("Symbol not found in any context\n")
	return CompileResult{}, fmt.Errorf("identifier '%s' not found in any scope", symbol.Name)
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

func (c *Compiler) InsertPutStringAtRuntime() {
	m := c.LLVMModule
	putstring := m.NewFunc("putstring", types.I1)
	putstring.Params = append(putstring.Params, ir.NewParam("paramValue", types.NewPointer(types.I8)))

	putstringEntry := putstring.NewBlock("putstring.entry")
	formatStrGlobal := m.NewGlobalDef("putstring.str", constant.NewCharArrayFromString("%s\n\x00"))

	indices := []value.Value{
		constant.NewInt(types.I64, 0),
		constant.NewInt(types.I64, 0),
	}
	formatStrPtr := putstringEntry.NewGetElementPtr(formatStrGlobal.Typ.ElemType, formatStrGlobal, indices...)

	printf := m.NewFunc("printf", types.I32, ir.NewParam("format", types.NewPointer(types.I8)))
	putstringEntry.NewCall(printf, formatStrPtr, putstring.Params[0])

	putstringEntry.NewRet(constant.NewInt(types.I1, 1))
}

func (c *Compiler) CompareStrings(str1, str2 value.Value) *ir.InstCall {
	m := c.LLVMModule
	for _, f := range m.Funcs {
		if f.Name() == "strcmp" {
			if _, ok := str1.(*ir.Global); ok {
				str1 = c.ctx.NewLoad(types.I8Ptr, str1)
			}
			return c.ctx.NewCall(f, str1, str2)
		}
	}
	return nil
}

func (c *Compiler) DeclareStrCmp() {
	m := c.LLVMModule
	strcmp := m.NewFunc("strcmp", types.I32,
		ir.NewParam("s1", types.I8Ptr),
		ir.NewParam("s2", types.I8Ptr))
	strcmp.Sig.Variadic = false
}

func (c *Compiler) insertRuntimeFunctions(node *ast.CallExpression, block *ir.Block) (CompileResult, error) {
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
		var printf *ir.Func
		for _, f := range m.Funcs {
			if f.Name() == "printf" {
				printf = f
			}
		}

		putintegerEntry.NewCall(
			// m.NewFunc("printf", types.I32, ir.NewParam("format", types.NewPointer(types.I8))),
			printf,
			formatStrPtr,
			loaded, // Pass the integer argument to printf.
		)
		putintegerEntry.NewRet(constant.NewInt(types.I1, 1))

		// callInstFromMain := funcMap["entry"].block.NewCall(putinteger, cr.Val)
		callInstFromBlock := block.NewCall(putinteger, cr.Val)
		return CompileResult{Type: "bool", Val: callInstFromBlock}, nil

	case "putstring":
		// putstring := m.NewFunc("putstring", types.I1)                                                      // Define the function returning bool
		// putstring.Params = append(putstring.Params, ir.NewParam("paramValue", types.NewPointer(types.I8))) // Pointer to i8 for C-style strings

		// putstringEntry := putstring.NewBlock("putstring.entry")
		// formatStrGlobal := m.NewGlobalDef("putstring.str", constant.NewCharArrayFromString("%s\n\x00")) // Format string for printf

		// // Define indices for the getelementptr instruction for accessing the string
		// indices := []value.Value{
		// 	constant.NewInt(types.I64, 0), // Index 0 to access the first element of the array
		// 	constant.NewInt(types.I64, 0), // Index 0 again, since it's a flat array
		// }
		// formatStrPtr := putstringEntry.NewGetElementPtr(formatStrGlobal.Typ.ElemType, formatStrGlobal, indices...)

		// cr, err := c.Compile(node.Arguments[0]) // Assume the argument is the string to print
		// if err != nil {
		// 	return CompileResult{}, fmt.Errorf("error compiling argument for putstring: %w", err)
		// }

		// // Call printf with the format string pointer and the string argument
		// putstringEntry.NewCall(
		// 	m.NewFunc("printf", types.I32, ir.NewParam("format", types.NewPointer(types.I8))),
		// 	formatStrPtr,
		// 	cr.Val, // Pass the string argument to printf
		// )
		// putstringEntry.NewRet(constant.NewInt(types.I1, 1)) // Return true (1)
		cr, err := c.Compile(node.Arguments[0]) // Compile the argument to get the value
		if err != nil {
			return CompileResult{}, err
		}
		var call *ir.InstCall

		for _, f := range m.Funcs {
			if f.Name() == "putstring" {
				call = block.NewCall(f, cr.Val)
			}
		}
		// Use the function in another context (e.g., calling it from a main block)
		// callInstFromBlock := block.NewCall(putstring, cr.Val)
		return CompileResult{Type: "bool", Val: call}, nil

	case "getbool":
		fmt.Print("Enter a boolean (0 or 1): ")
		var b int
		fmt.Scan(&b)
		// Eat newline
		fmt.Scanln()
		// return b != 0
		return CompileResult{Type: "bool", Val: constant.NewInt(types.I1, int64(b))}, nil
	case "getstring":
		// define global stdinp
		stdinp := m.NewGlobalDef("__stdinp", constant.NewNull(types.I8Ptr))
		stdinp.Typ = types.NewPointer(types.I8) // Setting the type to ptr (i8*)

		getline := m.NewFunc("getline", types.I64,
			ir.NewParam("buf", types.NewPointer(types.I8)),   // ptr %0
			ir.NewParam("size", types.NewPointer(types.I64)), // ptr %1
			ir.NewParam("file", types.I8Ptr),                 // ptr %2
		)
		// Define `getstring` function
		getString := m.NewFunc("getstring", types.I8Ptr)
		entry := getString.NewBlock("entry")

		bufPtr := entry.NewAlloca(types.I8Ptr)
		lenPtr := entry.NewAlloca(types.I64)
		read := entry.NewAlloca(types.I64)

		// Initialize pointers and integers
		entry.NewStore(constant.NewNull(types.I8Ptr), bufPtr)
		entry.NewStore(constant.NewInt(types.I64, 0), lenPtr)

		// Load external global stdin pointer
		stdinVal := entry.NewLoad(stdinp.Typ, stdinp)

		// Call getline function
		call := entry.NewCall(getline, bufPtr, lenPtr, stdinVal)

		// Store the result of getline
		entry.NewStore(call, read)

		// Load the buffer pointer for further use
		loadedBufPtr := entry.NewLoad(types.I8Ptr, bufPtr)
		loadedRead := entry.NewLoad(types.I64, read)

		// Adjust the string by setting the last character to null (presumed newline removal)
		lastIndex := entry.NewSub(loadedRead, constant.NewInt(types.I64, 1))
		charPtr := entry.NewGetElementPtr(types.I8, loadedBufPtr, lastIndex)
		entry.NewStore(constant.NewInt(types.I8, 0), charPtr)

		// Load the buffer pointer again for return
		finalBufPtr := entry.NewLoad(types.I8Ptr, bufPtr)
		entry.NewRet(finalBufPtr)

		getStringCall := block.NewCall(getString)

		print(entry.LLString(), " : getString-block\n")

		return CompileResult{Type: "string", Val: getStringCall}, nil

	case "getinteger":
		var i int
		fmt.Print("Enter an integer: ")
		fmt.Scan(&i)
		// Eat newline
		fmt.Scanln()
		return CompileResult{Type: "integer", Val: constant.NewInt(types.I64, int64(i))}, nil
	case "getfloat":
		fmt.Print("Enter a float: ")
		var f float64
		fmt.Scan(&f)
		// Eat newline
		fmt.Scanln()
		return CompileResult{Type: "float", Val: constant.NewFloat(types.Float, f)}, nil
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

	// print(node.Arguments[0].String(), " - arguments lmao\n")

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

func (c *Compiler) CompilePrefixExpression(block *ir.Block, node *ast.PrefixExpression, cr CompileResult) value.Value {
	switch node.Operator {
	case "-":
		return block.NewSub(constant.NewInt(types.I64, 0), cr.Val)
	default:
		panic("Unimplemented prefix expression")
	}
}

func (c *Compiler) CompileInfixExpression(block *ir.Block, node *ast.InfixExpression, cr CompileResult, cr_ CompileResult) value.Value {
	switch node.Operator {
	case "+":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		// Perform addition using the loaded values.
		return block.NewAdd(crVal, crVal_)
	case "-":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		// Perform addition using the loaded values.
		return block.NewSub(crVal, crVal_)
	case "*":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		// Perform addition using the loaded values.
		return block.NewMul(crVal, crVal_)
	case "/":
		return block.NewSDiv(cr.Val, cr_.Val)
	case "<":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		// Perform addition using the loaded values.
		return block.NewICmp(enum.IPredSLT, crVal, crVal_)
	case ">":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		return block.NewICmp(enum.IPredSGT, crVal, crVal_)
	case "<=":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		// Perform addition using the loaded values.
		return block.NewICmp(enum.IPredSLE, crVal, crVal_)
	case ">=":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		return block.NewICmp(enum.IPredSGE, crVal, crVal_)
	case "==":
		var crVal, crVal_ value.Value

		if cr_.Type == "string" {
			// Compare strings
			compareResult := c.CompareStrings(cr.Val, cr_.Val)
			// Check if strcmp returned 0 (strings are equal)
			isEqual := c.ctx.NewICmp(enum.IPredEQ, compareResult, constant.NewInt(types.I32, 0))
			return isEqual
		}
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		// Perform comparison using the loaded values.
		return block.NewICmp(enum.IPredEQ, crVal, crVal_)
	case "!=":
		var crVal, crVal_ value.Value
		// Check if cr.Val is an alloca.
		if cr.Val.Type().Equal(types.NewPointer(types.I64)) {
			// If it's an alloca, load its value.
			crVal = block.NewLoad(types.I64, cr.Val)
		} else {
			// If it's already a load instruction, use it directly.
			crVal = cr.Val
		}

		// Repeat the same process for cr_.Val.
		if cr_.Val.Type().Equal(types.NewPointer(types.I64)) {
			crVal_ = block.NewLoad(types.I64, cr_.Val)
		} else {
			crVal_ = cr_.Val
		}

		return block.NewICmp(enum.IPredNE, crVal, crVal_)
	default:
		panic("Unimplemented infix expression")
	}
}
