package object

import (
	"fmt"
	"math"
)

var Builtins = []struct {
	Name       string
	Builtin    *Builtin
	returnType string
}{
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return nil
		},
		},
		"nil",
	},
	{
		"getbool",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to read bool value from standard input
			// For simplicity, we'll assume the input is provided correctly
			return &Boolean{Value: true}
		},
		},
		"bool",
	},
	{
		"getinteger",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to read integer value from standard input
			// For simplicity, we'll assume the input is provided correctly
			return &Integer{Value: 0} // Placeholder value
		},
		},
		"integer",
	},
	{
		"getfloat",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to read float value from standard input
			// For simplicity, we'll assume the input is provided correctly
			return &Float{Value: 0.0} // Placeholder value
		},
		},
		"float",
	},
	{
		"getstring",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to read string value from standard input
			// For simplicity, we'll assume the input is provided correctly
			return &String{Value: ""} // Placeholder value
		},
		},
		"string",
	},
	{
		"putbool",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to write bool value to standard output
			// For simplicity, we'll just return a boolean indicating success
			return &Boolean{Value: true}
		},
		},
		"bool",
	},
	{
		"putinteger",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to write integer value to standard output
			// For simplicity, we'll just return a boolean indicating success
			return &Boolean{Value: true}
		},
		},
		"bool",
	},
	{
		"putfloat",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to write float value to standard output
			// For simplicity, we'll just return a boolean indicating success
			return &Boolean{Value: true}
		},
		},
		"bool",
	},
	{
		"putstring",
		&Builtin{Fn: func(args ...Object) Object {
			// Implement logic to write string value to standard output
			// For simplicity, we'll just return a boolean indicating success
			return &Boolean{Value: true}
		},
		},
		"bool",
	},
	{
		"sqrt",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			if args[0].Type() != INTEGER_OBJ {
				return &Error{Message: fmt.Sprintf("argument to `sqrt` must be INTEGER, got %s", args[0].Type())}
			}
			val := args[0].(*Integer).Value
			if val < 0 {
				return &Error{Message: fmt.Sprintf("cannot take square root of negative number %d", val)}
			}
			return &Float{Value: math.Sqrt(float64(val))}
		},
		},
		"integer",
	},
}
