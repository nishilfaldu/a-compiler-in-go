package main

import (
	"a-compiler-in-go/src/7west/src/7west/compiler"
	"a-compiler-in-go/src/7west/src/7west/lexer"
	"a-compiler-in-go/src/7west/src/7west/parser"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// the below code can be uncommented to be used like REPL
// func main() {
// 	// get the current user
// 	user, err := user.Current()
// 	if err != nil {
// 		panic(err)
// 	}
// 	// print a welcome message
// 	fmt.Printf("Hello %s! This is the 7West programming language!\n", user.Username)
// 	// util.Run("/Users/happyhome/Desktop/a-compiler-in-go/tests/correct")
// 	// start the REPL
// 	repl.Start(os.Stdin, os.Stdout)
// }

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: 7West <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read file: %s\n", err)
		os.Exit(1)
	}

	process(content)
}

func process(input []byte) {
	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	comp := compiler.NewWithState()
	_, err := comp.Compile(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compilation failed:\n%s\n", err)
		return
	}

	// Save the LLVM IR to a file
	outputFilename := "output.ll"
	err = ioutil.WriteFile(outputFilename, []byte(comp.LLVMModule.String()), 0644)
	if err != nil {
		fmt.Printf("Failed to write LLVM IR to file: %s\n", err)
		return
	}

	fmt.Println("Compiled LLVM IR is written to", outputFilename)

	// Compile LLVM IR to object file using llc
	objFilename := "output.o"
	cmd := exec.Command("llc", "-filetype=obj", "-o", objFilename, outputFilename)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to compile LLVM IR to object file: %s\n", err)
		return
	}

	// Compile object file to executable using clang
	executableFilename := "program"
	cmd = exec.Command("clang", "-o", executableFilename, objFilename)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to compile object file to executable: %s\n", err)
		return
	}

	fmt.Println("Executable is created:", executableFilename)

	// Execute the program
	cmd = exec.Command("./" + executableFilename)
	cmd.Stdin = os.Stdin // Connect the program's stdin to the terminal's stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running the program:")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to run the program: %s\n", err)
		return
	}

	fmt.Println("Program executed successfully")
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Fprintln(out, "Woops! We ran into some not-so-nice business here!\n parser errors:")
	for _, msg := range errors {
		fmt.Fprintln(out, "\t"+msg)
	}
}
