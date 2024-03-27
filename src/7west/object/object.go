package object

type BuiltinFunction func(args ...Object) Object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Builtin struct {
	Fn BuiltinFunction
}
