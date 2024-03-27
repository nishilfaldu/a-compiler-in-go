package ast

import (
	"a-compiler-in-go/src/7west/src/7west/token"
	"bytes"
	"fmt"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Declaration interface {
	Node
	declarationNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST our parser produces
type Program struct {
	Header *ProgramHeader
	Body   *ProgramBody
}

func (p *Program) TokenLiteral() string {
	if p.Header != nil {
		return p.Header.TokenLiteral()
	} else if p.Body != nil {
		return p.Body.TokenLiteral()
	} else {
		return ""
	}
}

// String returns a string representation of the Program node and all its children
func (p *Program) String() string {
	var out bytes.Buffer

	if p.Header != nil {
		out.WriteString(p.Header.String())
		out.WriteString("\n")
	}

	out.WriteString(p.Body.String())

	return out.String()
}

type ProgramHeader struct {
	Token      token.Token // the PROGRAM token
	Identifier *Identifier
}

func (ph *ProgramHeader) TokenLiteral() string { return ph.Token.Literal }
func (ph *ProgramHeader) String() string {
	var out bytes.Buffer

	out.WriteString(ph.TokenLiteral() + " ")
	out.WriteString(ph.Identifier.String() + " is ")

	return out.String()
}

type ProgramBody struct {
	Declarations []Declaration
	Statements   []Statement
}

// TokenLiteral returns the literal value of the token associated with this node
func (p *ProgramBody) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// String returns a string representation of the Program node and all its children
func (p *ProgramBody) String() string {
	var out bytes.Buffer

	for _, d := range p.Declarations {
		out.WriteString(d.String())
		out.WriteString(";")
	}
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString(";")
	}
	out.WriteString(" end program ")
	return out.String()
}

type VariableDeclaration struct {
	Token token.Token // the token.VARIABLE token
	Name  *Identifier
	Type  *TypeMark // the token.INTEGER | token.BOOLEAN | token.STRING | token.FLOAT token
}

func (vd *VariableDeclaration) declarationNode()     {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(vd.TokenLiteral() + " ")
	out.WriteString(vd.Name.String())
	out.WriteString(" : ")
	out.WriteString(vd.Type.String())
	out.WriteString(";")

	return out.String()
}

type TypeMark struct {
	Token token.Token
	Name  string
	Array *ArrayBound
}

func (tm *TypeMark) TokenLiteral() string { return tm.Token.Literal }
func (tm *TypeMark) String() string {
	if tm.Array != nil {
		return tm.Name + tm.Array.String()
	}
	return tm.Name
}

type ArrayBound struct {
	Value int64
}

func (ab *ArrayBound) String() string {
	return "[" + fmt.Sprint(ab.Value) + "]"
}

type Destination struct {
	Identifier *Identifier
	Expression Expression // This will represent the index expression (if present)
}

func (d *Destination) String() string {
	var out bytes.Buffer

	if d.Identifier != nil {
		out.WriteString(d.Identifier.String())
	}

	if d.Expression != nil {
		out.WriteString("[")
		out.WriteString(d.Expression.String())
		out.WriteString("]")
	}

	return out.String()
}

// <variable_declaration> ::= variable <identifier> = <expression> --- 1
type AssignmentStatement struct {
	// Token token.Token // the token.VARIABLE token
	Destination *Destination
	Value       Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return "AssignmentStatement" }

func (as *AssignmentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Destination.String())
	out.WriteString(" := ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type GlobalVariableDeclaration struct {
	Token               token.Token // the token.GLOBAL token
	VariableDeclaration *VariableDeclaration
}

func (gvd *GlobalVariableDeclaration) declarationNode()     {}
func (gvd *GlobalVariableDeclaration) TokenLiteral() string { return gvd.Token.Literal }
func (gvd *GlobalVariableDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(gvd.TokenLiteral() + " ")
	out.WriteString(gvd.VariableDeclaration.String())

	return out.String()
}

// type GlobalVariableStatement struct {
// 	Token    token.Token // the token.GLOBAL token
// 	Variable *VariableStatement
// }

// func (gvs *GlobalVariableStatement) statementNode()       {}
// func (gvs *GlobalVariableStatement) TokenLiteral() string { return gvs.Token.Literal }
// func (gvs *GlobalVariableStatement) String() string {
// 	var out bytes.Buffer

// 	out.WriteString(gvs.TokenLiteral() + " ")
// 	out.WriteString(gvs.Variable.String())

// 	return out.String()
// }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

// <return_statement> ::= return <expression> --- 2
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// here is where the actual value of the integer which was previously string
// going to contained as integer
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type StringLiteral struct {
	Token token.Token // the token.STRING token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *IfBlockStatement
	Alternative *IfBlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type IfBlockStatement struct {
	Token      token.Token // can be then or else?
	Statements []Statement
}

func (bs *IfBlockStatement) statementNode()       {}
func (bs *IfBlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *IfBlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ForBlockStatement struct {
	Statements []Statement
}

func (bs *ForBlockStatement) statementNode()       {}
func (bs *ForBlockStatement) TokenLiteral() string { return "for" }
func (bs *ForBlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LoopStatement struct {
	Token         token.Token          // the for token
	InitStatement *AssignmentStatement // the initialization statement
	Condition     Expression           // the loop condition expression
	Body          *ForBlockStatement   // the loop body
}

func (bs *LoopStatement) statementNode() {}

// TokenLiteral returns the literal value of the loop statement token
func (ls *LoopStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String returns the string representation of the loop statement
func (ls *LoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " (")
	out.WriteString(ls.InitStatement.String() + "; ")
	out.WriteString(ls.Condition.String() + ") {\n")
	out.WriteString(ls.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// type FunctionLiteral struct {
// 	Token      token.Token
// 	Parameters []*Identifier
// 	Body       *BlockStatement
// }

// func (fl *FunctionLiteral) expressionNode()      {}
// func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
// func (fl *FunctionLiteral) String() string {
// 	var out bytes.Buffer

// 	params := []string{}
// 	for _, p := range fl.Parameters {
// 		params = append(params, p.String())
// 	}

// 	out.WriteString(fl.TokenLiteral())
// 	out.WriteString("(")
// 	out.WriteString(strings.Join(params, ", "))
// 	out.WriteString(") ")
// 	out.WriteString(fl.Body.String())

//		return out.String()
//	}
type ProcedureDeclaration struct {
	Token  token.Token // the PROCEDURE token
	Header *ProcedureHeader
	Body   *ProcedureBody
}

// declarationNode marks this struct as a declaration node in the AST.
func (pd *ProcedureDeclaration) declarationNode()     {}
func (pd *ProcedureDeclaration) TokenLiteral() string { return pd.Token.Literal }
func (pd *ProcedureDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(pd.TokenLiteral() + " ")
	out.WriteString(pd.Header.String())
	out.WriteString(pd.Body.String())

	return out.String()
}

type ProcedureHeader struct {
	Token      token.Token // the PROCEDURE token
	Name       *Identifier
	TypeMark   *TypeMark
	Parameters []*VariableDeclaration
}

func (ph *ProcedureHeader) TokenLiteral() string { return ph.Token.Literal }

func (ph *ProcedureHeader) String() string {
	var out strings.Builder

	out.WriteString("procedure ")
	out.WriteString(ph.Name.String())
	out.WriteString(" : ")
	out.WriteString(ph.TypeMark.String())
	out.WriteString("(")
	for i, param := range ph.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(param.String())
	}
	out.WriteString(") ")

	return out.String()
}

type ProcedureBody struct {
	Declarations []Declaration
	Statements   []Statement
}

func (pb *ProcedureBody) TokenLiteral() string {
	if len(pb.Statements) > 0 {
		return pb.Statements[0].TokenLiteral()
	}
	return ""
}

func (pb *ProcedureBody) String() string {
	var out bytes.Buffer

	for _, decl := range pb.Declarations {
		out.WriteString(decl.String())
		out.WriteString(";")
	}
	out.WriteString(" begin ")
	for _, stmt := range pb.Statements {
		out.WriteString(stmt.String())
		out.WriteString(";")
	}
	out.WriteString(" end procedure")

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
