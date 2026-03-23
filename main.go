package main

import (
	"encoding/json"
	"fmt"
	"os"

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

	analysis, err := ai.Gemini(vuln)
	if err != nil {
		fmt.Println("AI Analysis Error:", err)
	} else {
		// Map AI analysis to vulnerabilities
		for i := range vuln {
			for _, a := range analysis {
				if vuln[i].CVE == a.CVE {
					vuln[i].FixCommand = a.FixCommand
					vuln[i].FixExplanation = a.FixExplanation
				}
			}
		}

		jsonData, err := json.MarshalIndent(vuln, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling final results:", err)
		} else {
			err = os.WriteFile("scanResult.json", jsonData, 0644)
			if err != nil {
				fmt.Println("Error writing final JSON:", err)
			} else {
				fmt.Println("Successfully saved merged AI results to scanResult.json")
			}
		}
	}
}
