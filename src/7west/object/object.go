package object

type BuiltinFunction func(args ...Object) Object

type ObjectType string

const (
	STRING_OBJ = "STRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Builtin struct {
	Fn BuiltinFunction
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
