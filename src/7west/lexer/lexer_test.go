package lexer

import (
	"a-compiler-in-go/src/7west/src/7west/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	// input := `let five = 5;
	// let ten = 10;
	// let doublenumber = 10.50;

	// let add = fn(x, y) {
	//   x + y;
	// };

	// let result = add(five, ten);
	// !-/+5;
	// 5 < 10 > 5;

	// if (5 < 10) {
	// 	return true;
	// } else {
	// 	return false;
	// }

	// // a sample slash comment with random text

	// /*
	// 	a sample block comment with
	// 	random text
	// 	over multiple lines
	// */

	// 10 == 10;
	// 10 != 9;
	// let s = "hello";
	// `

	input := `
	/* ok so it's also time to test
	/* nested comments; /*can you do it??*/
	*/
	*/
	PROGRAM grant IS
    global variable JAKE : integer;
	global variable ryan : integer[3];
    global variable zach : integer;
    variable tmp : integer;

	procedure if_proc : integer()
        variable i : integer;
		procedure dummy: float()
		// this should just hide the i of the outter environment
		variable i: float;
		variable tst : bool;
		begin
			i := 4.5;
			tst := putString("passed");
			return (0);
		end procedure;
		begin
			if(true) then jake := jake + 1;
			end if;
			return (0);
	end procedure;

	procedure for_proc : integer()
                variable i : integer;
		begin
			for(i := 0; i < zach)
			end for;
			return 1;
		end procedure;

	begin
		tmp := if_proc();
		tmp := for_proc();
	end program
	`
	// tests := []struct {
	// 	expectedType    token.TokenType
	// 	expectedLiteral string
	// }{
	// 	{token.LET, "let"},
	// 	{token.IDENT, "five"},
	// 	{token.ASSIGN, "="},
	// 	{token.INT, "5"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.LET, "let"},
	// 	{token.IDENT, "ten"},
	// 	{token.ASSIGN, "="},
	// 	{token.INT, "10"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.LET, "let"},
	// 	{token.IDENT, "doublenumber"},
	// 	{token.ASSIGN, "="},
	// 	{token.INT, "10.50"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.LET, "let"},
	// 	{token.IDENT, "add"},
	// 	{token.ASSIGN, "="},
	// 	{token.FUNCTION, "fn"},
	// 	{token.LPAREN, "("},
	// 	{token.IDENT, "x"},
	// 	{token.COMMA, ","},
	// 	{token.IDENT, "y"},
	// 	{token.RPAREN, ")"},
	// 	{token.LBRACE, "{"},
	// 	{token.IDENT, "x"},
	// 	{token.PLUS, "+"},
	// 	{token.IDENT, "y"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.RBRACE, "}"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.LET, "let"},
	// 	{token.IDENT, "result"},
	// 	{token.ASSIGN, "="},
	// 	{token.IDENT, "add"},
	// 	{token.LPAREN, "("},
	// 	{token.IDENT, "five"},
	// 	{token.COMMA, ","},
	// 	{token.IDENT, "ten"},
	// 	{token.RPAREN, ")"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.BANG, "!"},
	// 	{token.MINUS, "-"},
	// 	{token.SLASH, "/"},
	// 	{token.PLUS, "+"},
	// 	{token.INT, "5"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.INT, "5"},
	// 	{token.LT, "<"},
	// 	{token.INT, "10"},
	// 	{token.GT, ">"},
	// 	{token.INT, "5"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.IF, "if"},
	// 	{token.LPAREN, "("},
	// 	{token.INT, "5"},
	// 	{token.LT, "<"},
	// 	{token.INT, "10"},
	// 	{token.RPAREN, ")"},
	// 	{token.LBRACE, "{"},
	// 	{token.RETURN, "return"},
	// 	{token.TRUE, "true"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.RBRACE, "}"},
	// 	{token.ELSE, "else"},
	// 	{token.LBRACE, "{"},
	// 	{token.RETURN, "return"},
	// 	{token.FALSE, "false"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.RBRACE, "}"},
	// 	{token.INT, "10"},
	// 	{token.EQ, "=="},
	// 	{token.INT, "10"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.INT, "10"},
	// 	{token.NOT_EQ, "!="},
	// 	{token.INT, "9"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.LET, "let"},
	// 	{token.IDENT, "s"},
	// 	{token.ASSIGN, "="},
	// 	{token.STRING, "hello"},
	// 	{token.SEMICOLON, ";"},
	// 	{token.EOF, ""},
	// }

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PROGRAM, "program"},
		{token.IDENT, "grant"},
		{token.IS, "is"},
		{token.GLOBAL, "global"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "jake"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.SEMICOLON, ";"},
		{token.GLOBAL, "global"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "ryan"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.LSQBRACE, "["},
		{token.INTEGER, "3"},
		{token.RSQBRACE, "]"},
		{token.SEMICOLON, ";"},
		{token.GLOBAL, "global"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "zach"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.SEMICOLON, ";"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "tmp"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.SEMICOLON, ";"},
		{token.PROCEDURE, "procedure"},
		{token.IDENT, "if_proc"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "i"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.SEMICOLON, ";"},
		{token.PROCEDURE, "procedure"},
		{token.IDENT, "dummy"},
		{token.COLON, ":"},
		{token.FLOAT, "float"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "i"},
		{token.COLON, ":"},
		{token.FLOAT, "float"},
		{token.SEMICOLON, ";"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "tst"},
		{token.COLON, ":"},
		{token.BOOLEAN, "bool"},
		{token.SEMICOLON, ";"},
		{token.BEGIN, "begin"},
		{token.IDENT, "i"},
		{token.ASSIGN, ":="},
		{token.FLOAT, "4.5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "tst"},
		{token.ASSIGN, ":="},
		{token.IDENT, "putstring"},
		{token.LPAREN, "("},
		{token.STRING, "passed"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.LPAREN, "("},
		{token.INTEGER, "0"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.PROCEDURE, "procedure"},
		{token.SEMICOLON, ";"},
		{token.BEGIN, "begin"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.TRUE, "true"},
		{token.RPAREN, ")"},
		{token.THEN, "then"},
		{token.IDENT, "jake"},
		{token.ASSIGN, ":="},
		{token.IDENT, "jake"},
		{token.PLUS, "+"},
		{token.INTEGER, "1"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.IF, "if"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.LPAREN, "("},
		{token.INTEGER, "0"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.PROCEDURE, "procedure"},
		{token.SEMICOLON, ";"},
		{token.PROCEDURE, "procedure"},
		{token.IDENT, "for_proc"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.VARIABLE, "variable"},
		{token.IDENT, "i"},
		{token.COLON, ":"},
		{token.INTEGER, "integer"},
		{token.SEMICOLON, ";"},
		{token.BEGIN, "begin"},
		{token.FOR, "for"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.ASSIGN, ":="},
		{token.INTEGER, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.IDENT, "zach"},
		{token.RPAREN, ")"},
		{token.END, "end"},
		{token.FOR, "for"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.INTEGER, "1"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.PROCEDURE, "procedure"},
		{token.SEMICOLON, ";"},
		{token.BEGIN, "begin"},
		{token.IDENT, "tmp"},
		{token.ASSIGN, ":="},
		{token.IDENT, "if_proc"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "tmp"},
		{token.ASSIGN, ":="},
		{token.IDENT, "for_proc"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.END, "end"},
		{token.PROGRAM, "program"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
