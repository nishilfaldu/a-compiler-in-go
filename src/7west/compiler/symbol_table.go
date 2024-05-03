package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir/value"
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
	Value     value.Value
}

type SymbolTable struct {
	// for bidirectional traversal - inner is helpful during call expression
	Outer          *SymbolTable
	Inner          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
	FuncIndex      int
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	symbolTable := NewSymbolTable()
	symbolTable.Outer = outer
	symbolTable.Outer.Inner = symbolTable
	return symbolTable
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string, type_ string, param bool) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Type: type_}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else if param {
		symbol.Scope = ParamLocalScope
		symbol.Index = s.FuncIndex
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

func (s *SymbolTable) DefineBuiltin(index int, name string, returnType string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope, Type: returnType}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string, returnType string) Symbol {
	s.FuncIndex++
	symbol := Symbol{Name: name, Index: s.FuncIndex, Scope: FunctionScope, Type: returnType}
	s.store[name] = symbol
	return symbol
}

func getInnermostSymbolTable(symbolTable *SymbolTable) *SymbolTable {
	current := symbolTable
	for current.Inner != nil {
		current = current.Inner
	}
	return current
}

func (s *SymbolTable) getCurrentFunction() (FunctionType, bool) {
	for name := range s.store {
		if s.store[name].Scope == FunctionScope {
			return FunctionType{Name: name, ReturnType: s.store[name].Type}, true
		} else if s.store[name].Scope == BuiltinScope {
			return FunctionType{Name: "entry", ReturnType: "integer"}, true
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

	// Recursively print symbols in outer scopes
	if s.Outer != nil {
		fmt.Println("\nOuter Scope:")
		fmt.Println("------------")
		PrintSymbolTable(s.Outer)
	}
}
