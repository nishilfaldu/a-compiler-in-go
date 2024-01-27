package parser

import (
	"a-compiler-in-go/src/7west/src/7west/ast"
	"a-compiler-in-go/src/7west/src/7west/lexer"
	"a-compiler-in-go/src/7west/src/7west/token"
)

// Parser represents a parser
type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
}

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances the tokens for the parser
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses a program
func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
