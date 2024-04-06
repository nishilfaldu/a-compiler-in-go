package util

import (
	"a-compiler-in-go/src/7west/src/7west/compiler"
	"a-compiler-in-go/src/7west/src/7west/lexer"
	"a-compiler-in-go/src/7west/src/7west/parser"
	"os"
	"path/filepath"
	"strings"
)

func Run(dirPath string) {
	// Open the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		print("Error opening directory: %s\n", err)
		return
	}

	for _, file := range files {
		// Check if the file is a .correct file
		if strings.HasSuffix(file.Name(), ".src") {
			filePath := filepath.Join(dirPath, file.Name())
			// Read the contents of the file
			code, err := os.ReadFile(filePath)
			if err != nil {
				print("Error reading file %s: %s\n", filePath, err)
				continue
			}

			// Create a new lexer
			l := lexer.New(string(code))
			// Create a new parser
			p := parser.New(l)

			program := p.ParseProgram()

			// Check for errors
			if len(p.Errors()) != 0 {
				printParserErrors(p.Errors())
				continue
			}

			comp := compiler.NewWithState()
			_, err_ := comp.Compile(program)
			if err_ != nil {
				print("Woops! Compilation failed for file %s:\n %s\n", filePath, err)
				continue
			}

			print(program.String())
			print("\n")
		}
	}
}

func printParserErrors(errors []string) {
	print("Woops! We ran into some not-so-nice business here!\n")
	print(" parser errors:\n")
	for _, msg := range errors {
		print("\t" + msg + "\n")
	}
}
