package repl

import (
	"a-compiler-in-go/src/7west/src/7west/lexer"
	"a-compiler-in-go/src/7west/src/7west/parser"
	"bufio"
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
		// create a new parser
		p := parser.New(l)

		program := p.ParseProgram()
		// check for errors
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

		// When just Lexer existed:
		// // loop through the tokens until we reach the end of the input
		// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		// 	// print the token type and literal
		// 	fmt.Printf("%+v\n", tok)
		// }
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some not-so-nice business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
