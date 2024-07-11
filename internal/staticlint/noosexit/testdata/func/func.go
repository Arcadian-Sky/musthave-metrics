package main

import (
	"fmt"
	"os"
)

func main() {
	runExit()
	func() { os.Exit(1) }() // want "os.Exit call is not allowed in main function"
}

func runExit() {
	fmt.Print("os.Exit in func")
	os.Exit(1)
}
