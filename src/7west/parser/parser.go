// Technique: We already know the drill: we build our AST node and then try to
// fill its field by calling other parsing functions.
package parser

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/lexer"
	"a-compiler-in-go/src/7west/src/7west/token"
	"fmt"
	"strconv"
)

// Parser represents a parser
type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// Here we use iota to give the following constants incrementing numbers as values.
// The blank identifier _ takes the zero value and the following constants get
// assigned the values 1 to 7.
const (
	_ int = iota
	// LOWEST is the lowest precedence
	LOWEST
	// EQUALS is the precedence of the equality operator
	EQUALS // ==
	// LESSGREATER is the precedence of the comparison operators
	LESSGREATER // > or <
	// SUM is the precedence of the sum operator
	SUM // +
	// PRODUCT is the precedence of the product operator
	PRODUCT // *
	// PREFIX is the precedence of prefix operators
	PREFIX // -X or !X
	// CALL is the precedence of the call operator
	CALL // myFunction(X)
	// INDEX is the precedence of the index operator
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.AND:      LESSGREATER,
	token.OR:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	// p.registerPrefix(token.PROCEDURE, p.parseFunctionLiteral)
	// p.registerPrefix(token.LSQBRACE, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	// p.registerInfix(token.MINUS, p.parseInfixExpression)
	// p.registerInfix(token.SLASH, p.parseInfixExpression)
	// p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	// p.registerInfix(token.LPAREN, p.parseCallExpression)
	// p.registerInfix(token.LSQBRACE, p.parseIndexExpression)
	// Read two tokens, so currentToken and peekToken are both set
	// p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances the tokens for the parser
func (p *Parser) nextToken() {
	// print(p.currentToken.Literal, " . ", p.peekToken.Literal, "\n")

	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
	// print(p.currentToken.Literal, " . ", p.peekToken.Literal, "\n")

}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	// Parse the program header
	program.Header = p.parseProgramHeader()

	// Parse the program body
	program.Body = p.parseProgramBody()

	return program
}

// parseProgramHeader parses the initial program structure: PROGRAM ID IS ...
func (p *Parser) parseProgramHeader() *ast.ProgramHeader {
	programHeader := &ast.ProgramHeader{}
	// Ensure that the next token is "program"
	if !p.expectPeek(token.PROGRAM) {
		// print("bla bla\n")
		return nil
	}
	programHeader.Token = p.currentToken

	// Ensure that the next token is an identifier
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	programHeader.Identifier = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// Ensure that the next token is "is"
	if !p.expectPeek(token.IS) {
		return nil
	}
	p.nextToken()

	return programHeader
}

// ParseProgramBody parses a program body
func (p *Parser) parseProgramBody() *ast.ProgramBody {
	programBody := &ast.ProgramBody{}
	programBody.Statements = []ast.Statement{}
	programBody.Declarations = []ast.Declaration{}

	print("here", p.currentToken.Literal, "\n")

	// Parse declarations until "begin" keyword
	for !p.currentTokenIs(token.BEGIN) && !p.currentTokenIs(token.EOF) {
		decl := p.parseDeclaration()
		if decl != nil {
			switch d := decl.(type) {
			case *ast.VariableDeclaration:
				print("whats up 1\n")
				programBody.Declarations = append(programBody.Declarations, d)
			case *ast.ProcedureDeclaration:
				print("whats up 2\n")
				programBody.Declarations = append(programBody.Declarations, d)
			}
		}
		print("here hello ", p.currentToken.Literal, "\n")
		// p.nextToken()
	}

	print("reached here with ", p.currentToken.Literal, "\n")
	// TODO: might have advance to another token?
	// Check for "begin" keyword
	if !p.currentTokenIs(token.BEGIN) {
		return nil
	}

	p.nextToken() // Consume "begin"
	print("reached here with 1 ", p.currentToken.Literal, "\n")
	// Parse statements until "end program" keyword
	for !p.currentTokenIs(token.END) && !p.peekTokenIs(token.PROGRAM) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			programBody.Statements = append(programBody.Statements, stmt)
		}
		p.nextToken()
	}

	// TODO: edit this to expect format
	p.nextToken() // Consume "end"
	print("reached here with ", p.currentToken.Literal, "\n")
	p.nextToken() // Consume "program"
	print("reached here with ", p.currentToken.Literal, "\n")

	if p.currentTokenIs(token.SEMICOLON) {
		p.nextToken() // Consume ";"
	}

	return programBody
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	// case token.IDENTIFIER:
	// 	return p.parseAssignmentStatement()
	// case token.IF:
	// 	return p.parseIfStatement()
	// case token.FOR:
	// 	return p.parseLoopStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// return p.parseExpressionStatement()
		return nil
	}
}

func (p *Parser) parseProcedureDeclaration() *ast.ProcedureDeclaration {
	procedure := &ast.ProcedureDeclaration{}

	// Parse the program header
	procedure.Header = p.parseProcedureHeader()

	// Parse the program body
	procedure.Body = p.parseProcedureBody()

	return procedure
}

// parseDeclaration parses a statement
// With this change, parseDeclaration now returns an interface{},
// which allows it to return either *ast.Procedure or ast.Declaration.
// You can then handle the returned value accordingly in the calling code.
func (p *Parser) parseDeclaration() interface{} {
	print("running? parse declaration\n")
	switch p.currentToken.Type {
	case token.GLOBAL:
		return p.parseGlobalVariableDeclaration()
	case token.VARIABLE:
		return p.parseVariableDeclaration()
	case token.PROCEDURE:
		return p.parseProcedureDeclaration()
	// TODO: think about parseStatement here or not
	default:
		return nil
	}
}

// parseGlobalVariableStatement parses a global variable statement
func (p *Parser) parseGlobalVariableDeclaration() *ast.GlobalVariableDeclaration {
	gdecl := &ast.GlobalVariableDeclaration{Token: p.currentToken}

	if !p.expectPeek(token.VARIABLE) {
		return nil
	}
	gdecl.VariableDeclaration = p.parseVariableDeclaration()

	return gdecl

}

// parseVariableDeclaration parses a variable declaration
func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	print("here in variable declaration 1", p.currentToken.Literal, "\n")
	decl := &ast.VariableDeclaration{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	print("here in variable declaration 2", p.currentToken.Literal, "\n")
	decl.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}
	print("here in variable declaration 3", p.currentToken.Literal, "\n")
	p.nextToken() // consume the colon
	if !p.currentTokenIs(token.INTEGER) && !p.currentTokenIs(token.BOOLEAN) && !p.currentTokenIs(token.STRING) && !p.currentTokenIs(token.FLOAT) {
		return nil
	}
	print("bla bla bla\n")
	// Parse the type mark
	typeMark := &ast.TypeMark{Token: p.currentToken, Name: p.currentToken.Literal}

	// Optionally parse array bounds
	if p.peekTokenIs(token.LSQBRACE) {
		p.nextToken() // Consume '['

		// Parse the bound
		p.nextToken()
		bound, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
		if err != nil {
			return nil
		}

		typeMark.Array = &ast.ArrayBound{Value: bound}

		// Consume ']'
		if !p.expectPeek(token.RSQBRACE) {
			return nil
		}
	}
	decl.Type = typeMark

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	print("here in variable declaration 4", p.currentToken.Literal, "\n")
	print("whats token here: ", p.currentToken.Literal, "\n")

	p.nextToken() // consume the semicolon
	return decl
}

// parseVariableStatement parses a variable statement
func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	stmt := &ast.VariableStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// currentTokenIs checks if the current token is of a given type
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs checks if the peek token is of a given type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the peek token is of a given type
// and advances the tokens if it is
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// Errors returns an array of error messages
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

type (
	prefixParseFn func() ast.Expression
	// This argument is “left side” of the infix operator that’s being parsed
	infixParseFn func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))

	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// we check for an optional semicolon. Yes, it’s optional.
	// If the peekToken is a token.SEMICOLON, we advance so it’s the curToken.
	// If it’s not there, that’s okay too, we don’t add an error to the parser if it’s not there.
	// That’s because we want expression statements to have optional semicolons
	// (which makes it easier to type something like 5 + 5 into the REPL later on).
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFns[p.currentToken.Type]
	print(prefix, p.currentToken.Literal, "prefix\n")
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()
	print("after leftExp ", p.currentToken.Literal, " prefix\n")

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecendence() {
		print("this executes\n")
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		// we advance the tokens and call the infix function
		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseIdentifier parses an identifier (which is also an expression)
// Example: foobar;
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

// parseIntegerLiteral parses an integer literal (which is also an expression)
// Example: 5;
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))

	lit := &ast.IntegerLiteral{Token: p.currentToken}

	// we use the strconv.ParseInt function to convert the string literal to an integer value
	// and store it in the Value field of our IntegerLiteral struct
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		// if we encounter an error, we add it to the parser’s error list
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value

	return lit
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// // parsePrefixExpression parses an infix expression -5; !foobar; 5 + -10;
func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))

	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	// But in order to correctly parse a prefix expres- sion like -5 more than one
	// token has to be “consumed”. So after using p.curToken to build a
	// *ast.PrefixExpression node, the method advances the tokens and calls parseExpression
	// again. This time with the precedence of prefix operators as argument.
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression parses an infix expression 5 + 5; 5 - 5; 5 != 5; 5 == 5

func (p *Parser) peekPrecendence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecendence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))

	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// defer untrace(trace("parseGrouped Expression"))

	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// if (<condition>) <consequence> else <alternative>
func (p *Parser) parseIfExpression() ast.Expression {

	expression := &ast.IfExpression{Token: p.currentToken} // if token

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.THEN) {
		return nil
	}

	// if !p.expectPeek(token.LBRACE) {
	// 	return nil
	// }

	expression.Consequence = p.parseBlockStatement()

	// if the next token is an else token, we advance the tokens and parse the alternative block statement

	// The whole part of this method is constructed in a way that
	// allows an optional else but doesn’t add a parser error if there is none.
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // reach else token
		p.nextToken() // consume else token

		// if !p.expectPeek(token.LBRACE) {
		// 	return nil
		// }

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// This method also follows our parsing function protocol: the tokens
// get advanced just enough so that parseBlockStatement sits on the
// { with p.curToken being of type token.LBRACE
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// defer untrace(trace("parseBlockStatement"))

	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(token.END) && !p.peekTokenIs(token.IF) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if !p.currentTokenIs(token.END) && !p.peekTokenIs(token.IF) {
		return nil
	}
	p.nextToken() // Consume "end"
	p.nextToken() // Consume "if"

	return block
}

func (p *Parser) parseProcedureHeader() *ast.ProcedureHeader {
	header := &ast.ProcedureHeader{Token: p.currentToken}
	print("here in procedure header 1 ", p.currentToken.Literal, "\n")

	// Parse the procedure name
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	header.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// Ensure that the next token is ":"
	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken() // consume the colon

	if !p.currentTokenIs(token.INTEGER) && !p.currentTokenIs(token.BOOLEAN) && !p.currentTokenIs(token.STRING) && !p.currentTokenIs(token.FLOAT) {
		return nil
	}
	// Parse the type mark
	typeMark := &ast.TypeMark{Token: p.currentToken, Name: p.currentToken.Literal}

	header.TypeMark = typeMark

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	header.Parameters = p.parseProcedureParameters()

	// Ensure that the next token is ")"
	// if !p.expectPeek(token.RPAREN) {
	// 	return nil
	// }
	p.nextToken() // Consume the closing parenthesis ")"

	return header
}

func (p *Parser) parseProcedureBody() *ast.ProcedureBody {
	print("in procedure body\n")
	procedureBody := &ast.ProcedureBody{}
	procedureBody.Statements = []ast.Statement{}
	procedureBody.Declarations = []ast.Declaration{}

	// for p.currentToken.Type != token.END {
	// 	// Parse either a declaration or a statement
	// 	if decl, ok := p.parseDeclaration().(*ast.VariableDeclaration); ok {
	// 		body.Declarations = append(body.Declarations, decl)
	// 	} else if proc, ok := p.parseDeclaration().(*ast.ProcedureDeclaration); ok {
	// 		body.Declarations = append(body.Declarations, proc)
	// 	} else if stmt := p.parseStatement(); stmt != nil {
	// 		body.Statements = append(body.Statements, stmt)
	// 	}
	// 	p.nextToken()
	// }

	// Parse declarations until "begin" keyword
	// TODO: condition needs to be updated
	for !p.currentTokenIs(token.BEGIN) && !p.currentTokenIs(token.EOF) {
		decl := p.parseDeclaration()
		if decl != nil {
			switch d := decl.(type) {
			case *ast.VariableDeclaration:
				print("whats up 3\n")
				procedureBody.Declarations = append(procedureBody.Declarations, d)
			case *ast.ProcedureDeclaration:
				print("whats up 4\n")
				procedureBody.Declarations = append(procedureBody.Declarations, d)
			}
		}
		print("here hello ", p.currentToken.Literal, "\n")
		// p.nextToken()
	}
	// // TODO: might have advance to another token?
	// // Check for "begin" keyword
	if !p.currentTokenIs(token.BEGIN) {
		return nil
	}

	print("reached here in proc body 1", p.currentToken.Literal, "\n")
	p.nextToken() // Consume "begin"

	// Parse statements until "end procedure" keyword
	for !p.currentTokenIs(token.END) && !p.peekTokenIs(token.PROCEDURE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			procedureBody.Statements = append(procedureBody.Statements, stmt)
		}
		p.nextToken()
	}

	print("reached here in proc body 2 ", p.currentToken.Literal, "\n")
	p.nextToken() // Consume "end"
	print("reached here in proc body 3 ", p.currentToken.Literal, "\n")
	p.nextToken() // Consume "procedure"
	print("reached here with ", p.currentToken.Literal, "\n")
	// p.nextToken() // Consume "end"
	// p.nextToken() // Consume "procedure"

	if p.currentTokenIs(token.SEMICOLON) {
		p.nextToken() // Consume ";"
	}

	return procedureBody
}

func (p *Parser) parseProcedureParameters() []*ast.VariableDeclaration {
	parameters := []*ast.VariableDeclaration{}

	// if the next token is a RPAREN, the parameters are empty and we return an empty array
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	decl := p.parseVariableDeclarationAsParameter()
	parameters = append(parameters, decl)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to comma
		p.nextToken() // consume the comma
		print("here in procedure parameters ", p.currentToken.Literal, "\n")
		decl := p.parseVariableDeclarationAsParameter()
		parameters = append(parameters, decl)
	}

	if !p.expectPeek(token.RPAREN) {
		print("haha\n")
		return nil
	}

	return parameters
}

func (p *Parser) parseVariableDeclarationAsParameter() *ast.VariableDeclaration {
	decl := &ast.VariableDeclaration{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	decl.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken() // consume the colon

	if !p.currentTokenIs(token.INTEGER) && !p.currentTokenIs(token.BOOLEAN) && !p.currentTokenIs(token.STRING) && !p.currentTokenIs(token.FLOAT) {
		return nil
	}
	// Parse the type mark
	typeMark := &ast.TypeMark{Token: p.currentToken, Name: p.currentToken.Literal}

	// Optionally parse array bounds
	if p.peekTokenIs(token.LSQBRACE) {
		p.nextToken() // Consume '['

		// Parse the bound
		p.nextToken()
		bound, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
		if err != nil {
			return nil
		}

		typeMark.Array = &ast.ArrayBound{Value: bound}

		// Consume ']'
		if !p.expectPeek(token.RSQBRACE) {
			return nil
		}
	}
	decl.Type = typeMark

	return decl
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	// exp.Arguments = p.parseCallArguments() - before
	exp.Arguments = p.parseExpressionList(token.RPAREN) // - after (switched to a generalized function)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// if the next token is a RPAREN, the arguments are empty and we return an empty array
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST)) // parse first argument

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST)) // parse next argument
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RSQBRACE)
	return array
}

// this is a generalized method for parseExpressionList in an array and parseCallArguments
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		// skip the comma
		p.nextToken()
		// skip the next token
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RSQBRACE) {
		return nil
	}

	return exp
}
