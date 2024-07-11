package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("os.Exit in main")
	os.Exit(1) // want "os.Exit call is not allowed in main function"
}
