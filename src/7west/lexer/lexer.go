package lexer

import (
	"a-compiler-in-go/src/7west/src/7west/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// if we've reached the end of the input, set ch to 0, which is the ASCII code for the "NUL" character and is a
	// common way of saying "we don't have a value here"
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		// otherwise, set ch to the next character in the input
		l.ch = l.input[l.readPosition]
	}
	// advance our position in the input string
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// skip any whitespace
	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		// we're at the end of the input
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// if it's not one of the above characters, it could be the start of an identifier or an integer literal
		if isLetter(l.ch) {
			// if it's a letter, read in the entire identifier
			tok.Literal = l.readIdentifier()
			// and check whether it's a keyword
			tok.Type = token.LookupIdent(tok.Literal)
			// and return early
			// The early exit here, our return tok statement, is necessary because when calling readIdenti- fier(),
			// we call readChar() repeatedly and advance our readPosition and position fields past the last
			// character of the current identifier.
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			// if it's not a letter, it could be an integer literal

			// if it's not a digit, it's an illegal character
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	// remember our current position in the input string
	position := l.position
	// keep reading until we encounter a non-letter-character
	for isLetter(l.ch) {
		// advance our position in the input string
		l.readChar()
	}
	// return the substring of the input string from our starting position to our current position
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	// remember our current position in the input string
	position := l.position
	// keep reading until we encounter a non-digit-character
	for isDigit(l.ch) {
		// advance our position in the input string
		l.readChar()
	}
	// return the substring of the input string from our starting position to our current position
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	// we're only supporting ASCII characters
	return 'a' <= ch && ch <= 'z' || 'A' <= ch &&
		ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	// we're only supporting ASCII characters
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
