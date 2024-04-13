package llirgen

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Expr interface{ isExpr() Expr }
type EConstant interface {
	Expr
	isEConstant() EConstant
}
type EVoid struct{ EConstant }
type EBool struct {
	EConstant
	V bool
}
type EI32 struct {
	EConstant
	V int64
}
type EVariable struct {
	Expr
	Name string
}
type EAdd struct {
	Expr
	Lhs, Rhs Expr
}
type ELessThan struct {
	Expr
	Lhs, Rhs Expr
}

type Stmt interface{ isStmt() Stmt }
type SDefine struct {
	Stmt
	Name string
	Typ  types.Type
	Expr Expr
}
type SRet struct {
	Stmt
	Val Expr
}

// context definition starts here
type Context struct {
	*ir.Block
	parent *Context
	vars   map[string]value.Value
}

// this can be used to create initial context
func NewContext(b *ir.Block) *Context {
	return &Context{
		Block:  b,
		parent: nil,
		vars:   make(map[string]value.Value),
	}
}

// this can used to create new context with parent context
func (c *Context) NewContext(b *ir.Block) *Context {
	ctx := NewContext(b)
	ctx.parent = c
	return ctx
}

func (c Context) lookupVariable(name string) value.Value {
	if v, ok := c.vars[name]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.lookupVariable(name)
	} else {
		fmt.Printf("variable: `%s`\n", name)
		panic("no such variable")
	}
}

// context definition ends here

// TODO: maybe a separate function for codegen of Arrays

// TODO: a new block will be needed for the general program - this will be the entry block ("main")
// removed value as a parameters
func LLVMIRGlobalVariable(m *ir.Module, name string, type_ string) *ir.Global {
	global := m.NewGlobalDef(name,
		// constant.NewInt(types.I64, value)
		GetLLVMIRConstant(type_),
	)
	return global
}

func LLVMIRModule() *ir.Module {
	return ir.NewModule()
}

func LLVMIRFuncMain(m *ir.Module) *ir.Func {
	return m.NewFunc("main", types.I64)
}

func LLVMIRAlloca(block *ir.Block, name string, typ string) *ir.InstAlloca {
	// typ types.Type
	alloca := block.NewAlloca(GetLLVMIRType(typ))
	alloca.SetName(name)
	return alloca
}

func LLVMIRStore(block *ir.Block, value value.Value, ptr *ir.InstAlloca) *ir.InstStore {
	return block.NewStore(value, ptr)
}

func LLVMIRFunctionDefinition(m *ir.Module, name string, returnType string, params []*ast.VariableDeclaration) *ir.Func {
	irParams := make([]*ir.Param, len(params))
	for i, p := range params {
		param := ir.NewParam(p.Name.Value, GetLLVMIRType(p.Type.Name))
		irParams[i] = param
	}
	funcDef := m.NewFunc(name, GetLLVMIRType(returnType), irParams...)
	return funcDef
}

func LLVMIRFunctionBlock(fn *ir.Func, name string) *ir.Block {
	return fn.NewBlock(name)
}

func LLVMIRReturn(block *ir.Block, value value.Value) *ir.TermRet {
	return block.NewRet(value)
}

// func LLVMIRCall(block *ir.Block, callee *ir.Func, args ...compiler.Argument) *ir.InstCall {
// 	// Convert the Argument struct to []value.Value.
// 	var argValues []value.Value
// 	for _, arg := range args {
// 		argValues = append(argValues, GetLLVMIRConstant(arg.Type, arg.Value))
// 	}
// 	// args - i want to
// 	return block.NewCall(callee, argValues...)
// }

// compilation with ctx logic starts here
func compileConstant(e EConstant) constant.Constant {
	switch e := e.(type) {
	case *EI32:
		return constant.NewInt(types.I32, e.V)
	case *EBool:
		// we have no boolean in LLVM IR
		if e.V {
			return constant.NewInt(types.I1, 1)
		} else {
			return constant.NewInt(types.I1, 0)
		}
	case *EVoid:
		return nil
	}
	panic("unknown expression")
}

func (ctx *Context) compileExpr(e Expr) value.Value {
	switch e := e.(type) {
	case *EVariable:
		return ctx.lookupVariable(e.Name)
	case *EAdd:
		l, r := ctx.compileExpr(e.Lhs), ctx.compileExpr(e.Rhs)
		return ctx.NewAdd(l, r)
	case *ELessThan:
		l, r := ctx.compileExpr(e.Lhs), ctx.compileExpr(e.Rhs)
		return ctx.NewICmp(enum.IPredSLT, l, r)
	case EConstant:
		return compileConstant(e)
	}
	panic("unimplemented expression")
}

// func (ctx *Context) compileStmt(stmt Stmt) {
// 	if ctx.Parent != nil {
// 		return
// 	}
// 	f := ctx.Parent
// 	switch s := stmt.(type) {
// 	case *SDefine:
// 		v := ctx.NewAlloca(s.Typ)
// 		ctx.NewStore(ctx.compileExpr(s.Expr), v)
// 		ctx.vars[s.Name] = v
// 	case *SRet:
// 		ctx.NewRet(ctx.compileExpr(s.Val))
// 	case *SIf:
// 		thenCtx := ctx.NewContext(f.NewBlock("if.then"))
// 		thenCtx.compileStmt(s.Then)
// 		elseB := f.NewBlock("if.else")
// 		ctx.NewContext(elseB).compileStmt(s.Else)
// 		ctx.NewCondBr(ctx.compileExpr(s.Cond), thenCtx.Block, elseB)
// 		if !thenCtx.HasTerminator() {
// 			leaveB := f.NewBlock("leave.if")
// 			thenCtx.NewBr(leaveB)
// 		}
// 	case *SForLoop:
// 		loopCtx := ctx.NewContext(f.NewBlock("for.loop.body"))
// 		ctx.NewBr(loopCtx.Block)
// 		firstAppear := loopCtx.NewPhi(ir.NewIncoming(loopCtx.compileExpr(s.InitExpr), ctx.Block))
// 		loopCtx.vars[s.InitName] = firstAppear
// 		step := loopCtx.compileExpr(s.Step)
// 		firstAppear.Incs = append(firstAppear.Incs, ir.NewIncoming(step, loopCtx.Block))
// 		loopCtx.vars[s.InitName] = step
// 		leaveB := f.NewBlock("leave.for.loop")
// 		loopCtx.leaveBlock = leaveB
// 		loopCtx.compileStmt(s.Block)
// 		loopCtx.NewCondBr(loopCtx.compileExpr(s.Cond), loopCtx.Block, leaveB)
// 	}
// }

// compilation with ctx logic ends here

// helpers for condegen start here
func GetLLVMIRType(type_ string) types.Type {
	switch type_ {
	case "integer":
		return types.I64
	case "float":
		return types.Float
	case "bool":
		return types.I1
	case "string":
		// String = array of 8 bit ints
		return types.NewArray(8, types.I8)
	case "integer[]":
		return types.NewArray(8, types.I64)
	case "void":
		return types.Void
	default:
		panic("Invalid type name:" + type_)
		return nil // Default to nil if not specified in compiler.
	}
}

func GetLLVMIRConstant(type_ string) constant.Constant {
	// value interface{}
	switch type_ {
	case "integer":
		return constant.NewInt(types.I64, 0)
	case "float":
		return constant.NewFloat(types.Float, 0.0)
	case "bool":
		return constant.NewInt(types.I1, 1)
	case "string":
		return constant.NewCharArrayFromString("")
	case "integer[]":
		return constant.NewArray(&types.ArrayType{TypeName: types.I64.TypeName}, constant.NewInt(types.I64, 0))
	default:
		panic("Invalid type name:" + type_)
		return nil // Default to nil if not specified in compiler.
	}
}

// helpers for condegen end here
