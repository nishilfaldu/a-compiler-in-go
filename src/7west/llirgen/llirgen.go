package llirgen

import (
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
