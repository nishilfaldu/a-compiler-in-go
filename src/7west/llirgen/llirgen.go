package llirgen

import (
	"a-compiler-in-go/src/7west/src/7west/compiler"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

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
