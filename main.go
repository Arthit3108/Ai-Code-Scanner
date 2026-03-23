package main

import (
	"fmt"
	"ai-code-scanner/scanner"
)

func main() {
	fmt.Println("hello world go")

	target := "./petplace"
	scanner.RunTrivy(target, "json")
}

