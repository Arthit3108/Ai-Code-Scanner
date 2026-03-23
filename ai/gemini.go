package ai

import (
	"ai-code-scanner/scanner"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type AiAnalysis struct {
	CVE            string `json:"cve"`
	FixCommand     string `json:"fix_command"`
	FixExplanation string `json:"fix_explanation"`
}

func Gemini(vuln []scanner.CleanVuln) ([]AiAnalysis, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	vulnJSON, _ := json.MarshalIndent(vuln, "", "  ")

	prompt := fmt.Sprintf(`
		You are a security analysis assistant. 
		Analyze the following vulnerability scan results and provide a concise security report.

		Vulnerabilities found:
		<vulns>
		%s
		</vulns>

		Respond in this exact JSON array format, nothing else (no markdown block):
		[
			{
				"cve": "CVE-xxxx",
				"fix_command": "command here",
				"fix_explanation": "One sentence why this fixes it"
			}
		]`, string(vulnJSON))

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	var analysis []AiAnalysis
	if err := json.Unmarshal([]byte(result.Text()), &analysis); err != nil {
		fmt.Printf("AI Response: %s\n", result.Text())
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return analysis, nil
}