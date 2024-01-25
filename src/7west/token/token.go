package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	line    int // to track what line is the token on
}

const (
	// special types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	ERROR   = "ERROR"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...

	// Operators
	// ASSIGN   = "="
	ASSIGN   = ":="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ         = "=="
	NOT_EQ     = "!="
	LESS_EQ    = "<="
	GREATER_EQ = ">="

	AND = "&"
	OR  = "|"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LSQBRACE = "["
	RSQBRACE = "]"

	// loops
	FOR = "FOR"

	// Keywords
	// FUNCTION = "FUNCTION"
	LET   = "LET"
	TRUE  = "TRUE"
	FALSE = "FALSE"
	// ELSE      = "ELSE"
	RETURN    = "RETURN"
	GLOBAL    = "GLOBAL"
	PROGRAM   = "PROGRAM"
	IS        = "IS"
	VARIABLE  = "VARIABLE"
	IF        = "IF"
	THEN      = "THEN"
	PROCEDURE = "PROCEDURE"
	BEGIN     = "BEGIN"
	END       = "END"

	// Data types
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"
	STRING  = "STRING"
	BOOLEAN = "BOOLEAN"
	NOT     = "NOT"
)

var keywords = map[string]TokenType{
	// "fn":       FUNCTION,
	// "else":      ELSE,
	"let":       LET,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"then":      THEN,
	"return":    RETURN,
	"global":    GLOBAL,
	"program":   PROGRAM,
	"is":        IS,
	"variable":  VARIABLE,
	"procedure": PROCEDURE,
	"begin":     BEGIN,
	"end":       END,
	"bool":      BOOLEAN,
	"integer":   INTEGER,
	"float":     FLOAT,
	"string":    STRING,
	"for":       FOR,
	"not":       NOT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	// The below is a longer version of code for a better understanding of Go
	// tok, ok := keywords[ident]
	// if ok {
	// 	return tok
	// }
	return IDENT
}
