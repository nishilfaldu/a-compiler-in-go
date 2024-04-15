package llirgen

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

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

// context definition ends here

// TODO: maybe a separate function for codegen of Arrays

// TODO: a new block will be needed for the general program - this will be the entry block ("main")
// removed value as a parameters
func LLVMIRGlobalVariable(m *ir.Module, name string, type_ string) *ir.Global {
	// ptrType := types.NewPointer(GetLLVMIRType(type_))
	global := m.NewGlobalDef(name,
		// constant.NewInt(types.I64, value)
		GetLLVMIRConstant(type_),
		// constant.NewNull(ptrType),
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

// func LLVMIRFunctionCall(block *ir.Block, fn *ir.Func, args ...ast.Expression) *ir.TermC {
// 	fnArgs := make([]*value.Value, len(args))
// 	for i, arg := range args {
// 		fnArgs[i] = &value.Value{}
// 	}

// 	block.NewCall(fn)
// 	return block.NewCall(fn, fnArgs...)
// }

func LLVMIRFunctionBlock(fn *ir.Func, name string) *ir.Block {
	return fn.NewBlock(name)
}

func LLVMIRReturn(block *ir.Block, value value.Value) *ir.TermRet {
	return block.NewRet(value)
}

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
	case "string[]":
		return constant.NewArray(&types.ArrayType{TypeName: types.NewArray(8, types.I8).TypeName}, constant.NewCharArrayFromString(""))
	case "float[]":
		return constant.NewArray(&types.ArrayType{TypeName: types.Float.TypeName}, constant.NewFloat(types.Float, 0.0))
	case "void":
		return nil
	default:
		panic("Invalid type name:" + type_)
	}
}
