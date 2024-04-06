package compiler

import (
	"fmt"
)

type SymbolScope string

const (
	GlobalScope     SymbolScope = "GLOBAL"
	LocalScope      SymbolScope = "LOCAL"
	ParamLocalScope SymbolScope = "PARAMLOCAL"
	BuiltinScope    SymbolScope = "BUILTIN"
	FreeScope       SymbolScope = "FREE"
	FunctionScope   SymbolScope = "FUNCTION"
)

type FunctionType struct {
	Name       string
	ReturnType string
}

type Symbol struct {
	Name      string
	Scope     SymbolScope
	Index     int
	Type      string
	ArraySize int64
}

type SymbolTable struct {
	// for bidirectional traversal - inner is helpful during call expression
	Outer          *SymbolTable
	Inner          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
	paramIndex     int

	FreeSymbols       []Symbol
	functionTypeStack []FunctionType // Stack to track function names and return types
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()
	symbolTable.Outer = outer
	symbolTable.Outer.Inner = symbolTable
	return symbolTable
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s, paramIndex: 0}
}

// func (s *SymbolTable) Define(name string, type_ string, param bool) Symbol {
// 	var symbolIndex int
// 	if param {
// 		symbolIndex = s.paramIndex
// 		s.paramIndex++
// 	} else {
// 		symbolIndex = s.numDefinitions
// 		s.numDefinitions++
// 	}
// 	symbol := Symbol{Name: name, Index: symbolIndex, Type: type_}
// 	if s.Outer == nil {
// 		symbol.Scope = GlobalScope
// 	} else if param {
// 		symbol.Scope = ParamLocalScope
// 	} else {
// 		symbol.Scope = LocalScope
// 	}
// 	s.store[name] = symbol
// 	s.numDefinitions++
// 	return symbol
// }

func (s *SymbolTable) Define(name string, type_ string, param bool) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Type: type_}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else if param {
		symbol.Scope = ParamLocalScope
	} else {
		symbol.Scope = LocalScope
	}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) DefineGlobal(name string, type_ string) Symbol {
	// Check if the variable already exists in the global scope
	if symbol, ok := s.store[name]; ok && symbol.Scope == GlobalScope {
		return symbol // Return the existing global variable symbol
	}

	// If not, define the variable in the global scope
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope, Type: type_}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// func (s *SymbolTable) NewChildSymbolTable() *SymbolTable {
// 	child := NewSymbolTable()
// 	child.Outer = s
// 	// s.Children = append(s.Children, child)
// 	return child
// }

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		// free := s.defineFree(obj)
		// return free, true
	}
	return obj, ok
}

func (s *SymbolTable) ResolveInner(name string) (Symbol, bool) {
	// Try to resolve the symbol in the current scope
	obj, ok := s.store[name]
	if ok {
		return obj, true
	}

	// If not found in the current scope, try inner scopes recursively
	if s.Inner != nil {
		obj, ok := s.Inner.ResolveInner(name)
		if ok {
			return obj, true
		}
	}

	// If not found in any inner scopes, return false
	return Symbol{}, false
}

func (s *SymbolTable) DefineArray(name string, typeName string, size int64, scope SymbolScope) Symbol {
	symbol := Symbol{Name: name, Type: typeName, ArraySize: size, Index: s.numDefinitions, Scope: scope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve a symbol by recursively searching in the current and outer scopes
// func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
// 	obj, ok := s.store[name]
// 	if !ok && s.Outer != nil {
// 		obj, ok = s.Outer.Resolve(name)
// 	}
// 	return obj, ok
// }

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	s.store[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) DefineBuiltin(index int, name string, returnType string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope, Type: returnType}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string, returnType string) Symbol {
	print("name: ", name, "\n")
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope, Type: returnType}
	s.store[name] = symbol
	return symbol
}

// Function to push a new function type onto the stack
func (s *SymbolTable) pushFunction(name, returnType string) {
	s.functionTypeStack = append(s.functionTypeStack, FunctionType{Name: name, ReturnType: returnType})
}

// Function to pop the top function type from the stack
func (s *SymbolTable) popFunction() {
	if len(s.functionTypeStack) > 0 {
		s.functionTypeStack = s.functionTypeStack[:len(s.functionTypeStack)-1]
	}
}

func (s *SymbolTable) printFunctionDetails() {
	print(len(s.functionTypeStack), "length\n")
	// for _, f := range s.functionTypeStack {
	// 	fmt.Println(f.Name, f.ReturnType)
	// }
}

func getInnermostSymbolTable(symbolTable *SymbolTable) *SymbolTable {
	current := symbolTable
	for current.Inner != nil {
		current = current.Inner
	}
	return current
}

// Function to get the top function type from the stack
//
//	func (s *SymbolTable) getCurrentFunction() (FunctionType, bool) {
//		if len(s.functionTypeStack) > 0 {
//			return s.functionTypeStack[len(s.functionTypeStack)-1], true
//		}
//		return FunctionType{}, false
//	}
func (s *SymbolTable) getCurrentFunction() (FunctionType, bool) {
	for name := range s.store {
		if s.store[name].Scope == FunctionScope {
			return FunctionType{Name: name, ReturnType: s.store[name].Type}, true
		}
	}
	// If no function name found
	return FunctionType{}, false
}

// Function to check if the current scope is global
func (s *SymbolTable) IsGlobalScope() bool {
	return s.Outer == nil
}

// PrintSymbolTable prints the contents of the symbol table along with labels.
func PrintSymbolTable(s *SymbolTable) {
	fmt.Println("Symbol Table:")
	fmt.Println("=============")

	// Print symbols defined in the current scope
	// Print symbols defined in the current scope
	fmt.Println("Current Scope:")
	fmt.Println("-------------")
	for name, sym := range s.store {
		fmt.Printf("Name: %-10s Type: %-10s Scope: %-10s Index: %-5d ArraySize: %-5d\n", name, sym.Type, sym.Scope, sym.Index, sym.ArraySize)
	}

	// Print free symbols
	if len(s.FreeSymbols) > 0 {
		fmt.Println("\nFree Symbols:")
		fmt.Println("-------------")
		for _, sym := range s.FreeSymbols {
			fmt.Printf("Name: %-10s Scope: %-10s Index: %-5d\n", sym.Name, sym.Scope, sym.Index)
		}
	}

	// Recursively print symbols in outer scopes
	if s.Outer != nil {
		fmt.Println("\nOuter Scope:")
		fmt.Println("------------")
		PrintSymbolTable(s.Outer)
	}
}
