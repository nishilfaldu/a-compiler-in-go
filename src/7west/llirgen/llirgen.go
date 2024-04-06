package llirgen

import (
	"a-compiler-in-go/src/7west/src/7west/compiler"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// type Context struct {
// 	Block  *ir.Block
// 	parent *Context
// 	vars   map[string]value.Value
// }

// // this can be used to create initial context
// func NewContext(b *ir.Block) *Context {
// 	return &Context{
// 		Block:  b,
// 		parent: nil,
// 		vars:   make(map[string]value.Value),
// 	}
// }

// // this can used to create new context with parent context
// func (c *Context) NewContext(b *ir.Block) *Context {
// 	ctx := NewContext(b)
// 	ctx.parent = c
// 	return ctx
// }

func main() {
	m := ir.NewModule()

	// Generate global variable definition.
	// globalG := LLVMIRGlobalVariable(m, "g", 2)

	// Generate function definition for 'add'.
	// funcAdd := newAddFunction(m)

	// Generate function definition for 'main'.
	// funcMain := newMainFunction(m, funcAdd, globalG)

	// Print LLVM IR.
	println(m.String())
}

// TODO: a new block will be needed for the general program - this will be the entry block ("main")

func LLVMIRGlobalVariable(m *ir.Module, name string, value int64) *ir.Global {
	global := m.NewGlobalDef(name, constant.NewInt(types.I64, value))
	return global
}

func LLVMIRFunctionDefinition(m *ir.Module, name string, params []compiler.Parameter) *ir.Func {
	irParams := make([]*ir.Param, len(params))
	for i, p := range params {
		param := ir.NewParam(p.Name, GetLLVMIRType(p.Type))
		irParams[i] = param
	}
	funcDef := m.NewFunc(name, types.I64, irParams...)
	return funcDef
}

func LLVMIRFunctionBlock(fn *ir.Func, name string) *ir.Block {
	return fn.NewBlock(name)
}

func LLVMIRReturn(block *ir.Block, value value.Value) *ir.TermRet {
	return block.NewRet(value)
}

func LLVMIRCall(block *ir.Block, callee *ir.Func, args ...compiler.Argument) *ir.InstCall {
	// Convert the Argument struct to []value.Value.
	var argValues []value.Value
	for _, arg := range args {
		argValues = append(argValues, GetLLVMIRConstant(arg.Type, arg.Value))
	}
	// args - i want to
	return block.NewCall(callee, argValues...)
}

// helpers for condegen start here
func GetLLVMIRType(type_ string) types.Type {
	switch type_ {
	case "integer":
		return types.I64
	case "float":
		return types.Float
	case "void":
		return types.Void
	default:
		return nil // Default to nil if not specified in compiler.
	}
}

func GetLLVMIRConstant(type_ string, value interface{}) constant.Constant {
	switch type_ {
	case "integer":
		return constant.NewInt(types.I64, value.(int64))
	case "float":
		return constant.NewFloat(types.Float, value.(float64))
	default:
		return nil // Default to nil if not specified in compiler.
	}
}
