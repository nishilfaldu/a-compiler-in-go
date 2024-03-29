package compiler

import (
	"fmt"
)

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name      string
	Scope     SymbolScope
	Index     int
	Type      string
	ArraySize int64
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int

	FreeSymbols []Symbol
	// slice to store child symbol tables representing nested scopes
	// Children []*SymbolTable
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()
	symbolTable.Outer = outer
	return symbolTable
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string, type_ string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Type: type_}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
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

func (s *SymbolTable) DefineArray(name string, typeName string, size int64) Symbol {
	symbol := Symbol{Name: name, Type: typeName, ArraySize: size, Index: s.numDefinitions}
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

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	print("name: ", name, "\n")
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	s.store[name] = symbol
	return symbol
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
		fmt.Printf("Name: %-10s Type: %-10s Scope: %-10s Index: %-5d\n", name, sym.Type, sym.Scope, sym.Index)
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
