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
			fmt.Printf("Trivy Scan Error: %v\n", err)
		}

		// Run Gitleaks
		fmt.Println("Running Gitleaks scan...")
		secrets, err := scanner.RunGitleaks(target)
		if err != nil {
			fmt.Printf("Gitleaks Scan Error: %v\n", err)
		}

		// Only call Gemini if there are findings to analyze
		if len(vuln) > 0 || len(secrets) > 0 {
			analysis, err := ai.Gemini(vuln, secrets)
			if err != nil {
				fmt.Printf("\n⚠️  AI Analysis Error (continuing with raw scan results): %v\n", err)
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
		} else {
			fmt.Println("\nNo vulnerabilities or secrets found. Skipping AI analysis.")
		}

		// Always write output file
		outputPath := "scanResult.json"
		if output != "" {
			outputPath = output
		}

		// Combine results for JSON output
		type CombinedResult struct {
			Vulnerabilities []scanner.CleanVuln      `json:"vulnerabilities"`
			Secrets         []scanner.GitleaksFinding `json:"secrets"`
		}

		combined := CombinedResult{
			Vulnerabilities: vuln,
			Secrets:         secrets,
		}

		var dataToWrite []byte
		if outputType == "json" {
			jsonData, err := json.MarshalIndent(combined, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling final results:", err)
				return
			}
			dataToWrite = jsonData
			fmt.Println(string(dataToWrite))
		} else {
			dataToWrite = rawOutput
		}

		// Print summary
		critical := 0
		high := 0
		medium := 0
		low := 0
		for _, v := range vuln {
			switch v.Severity {
			case "CRITICAL":
				critical++
			case "HIGH":
				high++
			case "MEDIUM":
				medium++
			case "LOW":
				low++
			}
		}

		fmt.Println("\n---------------------------------------")
		fmt.Println("Scan Results Summary:")
		fmt.Printf("Total Vulnerabilities (Trivy): %d\n", len(vuln))
		fmt.Printf("🔴 Critical: %d\n", critical)
		fmt.Printf("🟠 High: %d\n", high)
		fmt.Printf("🟡 Medium: %d\n", medium)
		fmt.Printf("🟢 Low: %d\n", low)
		fmt.Printf("\nTotal Secrets (Gitleaks): %d\n", len(secrets))
		fmt.Println("---------------------------------------")

		err = os.WriteFile(outputPath, dataToWrite, 0644)
		if err != nil {
			fmt.Println("Error writing final output:", err)
		} else {
			fmt.Printf("Successfully saved results to %s\n", outputPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&target, "target", "t", ".", "Target directory to scan")
	scanCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: scanResult.json)")
	scanCmd.Flags().StringVar(&severity, "severity", "HIGH,CRITICAL", "Severity levels to scan (e.g. HIGH,CRITICAL)")
	scanCmd.Flags().StringVar(&outputType, "output-type", "json", "Output type (default: json)")
}
