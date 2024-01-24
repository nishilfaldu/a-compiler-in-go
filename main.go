package main

import (
	"a-compiler-in-go/src/7west/src/7west/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	// get the current user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// print a welcome message
	fmt.Printf("Hello %s! This is the 7West programming language!\n", user.Username)
	// start the REPL
	repl.Start(os.Stdin, os.Stdout)
}
