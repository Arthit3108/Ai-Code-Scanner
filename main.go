package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"ai-code-scanner/scanner"
	"ai-code-scanner/ai"

)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	fmt.Println("hello world go")

	target := "./petplace"
	vuln, err := scanner.RunTrivy(target, "json")
	if err != nil {
		fmt.Println(err)
	}

	ai.Gemini(vuln)
	// fmt.Println(vuln)

	

}

