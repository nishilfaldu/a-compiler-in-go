package lexer

import (
	"a-compiler-in-go/src/7west/src/7west/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // line count
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
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
		if l.peekChar() == '=' {
			ch := l.ch
			// read the next character
			l.readChar()
			// and set the token to the equality operator
			// cannot use newToken here as it requires single byte for the literal argument
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			// otherwise, it's just a regular assignment operator
			tok = newToken(token.ASSIGN, l.ch)
		}
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
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		// if it's an exclamation point, check whether it's part of the "not equals" operator
		if l.peekChar() == '=' {
			ch := l.ch
			// read the next character
			l.readChar()
			// and set the token to the not-equals operator
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			// otherwise, it's just a regular bang operator
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			// read the next character /
			l.readChar()
			// if it's a double slash, skip the entire comment
			l.skipSlashComment()
			// and return the next token which is important otherwise the program hangs
			return l.NextToken()
		} else if l.peekChar() == '*' {
			// read the next character *
			l.readChar()
			// if it's a slash star, skip the entire comment
			l.skipBlockComment()
			// and continue from where we left off in the code and skip */
			l.readChar()
			l.readChar()
			// and return the next token which is important otherwise the program hangs
			return l.NextToken()
		}

		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		// read the next character to skip the first "
		l.readChar()
		str := l.readString()
		// if it's a double quote, read the entire string
		if str == "" {
			// if we've reached the end of the input, return an error
			tok = newErrorToken(token.ERROR, "unterminated string")
		} else {
			tok.Type = token.STRING
			tok.Literal = str
		}
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
			// and return early
			return tok
		} else {
			// if it's not a letter, it could be an integer literal

			// if it's not a digit, it's an illegal character or you can throw an error
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

func (l *Lexer) readString() string {
	// remember our current position in the input string
	position := l.position
	// keep reading until we encounter a double quote
	for {
		if l.ch == '"' {
			break
		} else if l.ch == 0 {
			// if we've reached the end of the input, return an empty string
			return ""
		} else if l.ch == '\n' {
			// if we've reached a newline, increment the line count
			l.line += 1
		}
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

// func (l *Lexer) isAtEnd() bool {
// 	return l.ch == 0
// }

func newErrorToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipSlashComment() {
	// keep reading until we encounter a newline
	for l.ch != '\n' {
		// advance our position in the input string
		l.readChar()
	}
	// increment the line count
	l.line += 1
}

func (l *Lexer) skipBlockComment() {
	// keep reading until we encounter a newline
	l.readChar()
	print(l.ch, "\n")
	for l.ch != '*' && l.peekChar() != '/' {
		// advance our position in the input string
		l.readChar()
	}
	// increment the line count - TODO: this needs to be edited
	l.line += 1
}

func (l *Lexer) peekChar() byte {
	// if we've reached the end of the input, return 0
	// peeking doesn't need incrementing
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		// otherwise, set ch to the next character in the input
		return l.input[l.readPosition]
	}
}
