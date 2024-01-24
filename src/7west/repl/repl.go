package repl

import (
	"a-compiler-in-go/src/7west/src/7west/lexer"
	"a-compiler-in-go/src/7west/src/7west/token"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// print the prompt
		print(PROMPT)
		// read a line of input
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		// get the line of input
		line := scanner.Text()
		// create a new lexer
		l := lexer.New(line)

		// loop through the tokens until we reach the end of the input
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			// print the token type and literal
			fmt.Printf("%+v\n", tok)
		}
	}
}
