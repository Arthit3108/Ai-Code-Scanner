package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"encoding/json"
	"os"

	"github.com/joho/godotenv"

	"ai-code-scanner/ai"
	"ai-code-scanner/scanner"
)

var target string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a target",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from scan")

		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("Error loading .env file:", err)
		}

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
					if vuln[i].ID == a.ID {
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
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&target, "target", "t", "", "Target directory to scan")
}
