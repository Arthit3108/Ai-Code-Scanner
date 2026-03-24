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
var outputType string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a target",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load(".env")

		vuln, rawOutput, err := scanner.RunTrivy(target, outputType, severity)
		if err != nil {
			fmt.Println(err)
		}

		analysis, err := ai.Gemini(vuln)
		if err != nil {
			fmt.Println("AI Analysis Error (continuing without AI):", err)
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
		}

		// Always write output file
		outputPath := "scanResult.json"
		if output != "" {
			outputPath = output
		}

		var dataToWrite []byte
		if outputType == "json" {
			jsonData, err := json.MarshalIndent(vuln, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling final results:", err)
				return
			}
			dataToWrite = jsonData
			fmt.Println(string(dataToWrite))
		} else {
			dataToWrite = rawOutput
		}

		err = os.WriteFile(outputPath, dataToWrite, 0644)
		if err != nil {
			fmt.Println("Error writing final output:", err)
		} else {
			fmt.Println("Successfully saved results to", outputPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&target, "target", "t", ".", "Target directory to scan")
	scanCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: scanResult.json)")
	scanCmd.Flags().StringVar(&severity, "severity", "HIGH,CRITICAL", "Severity levels to scan (e.g. HIGH,CRITICAL)")
	scanCmd.Flags().StringVar(&outputType, "output_type", "json", "Output type (default: json)")
}
