package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"ai-code-scanner/ai"
	"ai-code-scanner/scanner"
)

var target string
var output string
var severity string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a target",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load(".env")

		vuln, err := scanner.RunTrivy(target, "json", severity)
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
			fmt.Println(string(jsonData))
			if err != nil {
				fmt.Println("Error marshaling final results:", err)
			} else {
				outputPath := "scanResult.json"
				if output != "" {
					outputPath = output
				}
				err = os.WriteFile(outputPath, jsonData, 0644)
				if err != nil {
					fmt.Println("Error writing final JSON:", err)
				} else {
					fmt.Println("Successfully saved merged AI results to", outputPath)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&target, "target", "t", ".", "Target directory to scan")
	scanCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: scanResult.json)")
	scanCmd.Flags().StringVar(&severity, "severity", "HIGH,CRITICAL", "Severity levels to scan (e.g. HIGH,CRITICAL)")
}
