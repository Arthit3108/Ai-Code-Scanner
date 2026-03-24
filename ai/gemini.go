package ai

import (
	"ai-code-scanner/scanner"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/genai"
	"strings"
)

type AiAnalysis struct {
	ID             int    `json:"id"`
	FixCommand     string `json:"fix_command"`
	FixExplanation string `json:"fix_explanation"`
}

func Gemini(vuln []scanner.CleanVuln, secrets []scanner.GitleaksFinding) ([]AiAnalysis, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	vulnJSON, _ := json.MarshalIndent(vuln, "", "  ")
	secretsJSON, _ := json.MarshalIndent(secrets, "", "  ")

	prompt := fmt.Sprintf(`
		You are a security analysis assistant. 
		Analyze the following vulnerability scan results and secrets detection and provide a concise security report.

		Vulnerabilities found by Trivy:
		<vulns>
		%s
		</vulns>

		Secrets found by Gitleaks:
		<secrets>
		%s
		</secrets>

		Respond in this exact JSON array format, nothing else (no markdown block):
		[
			{
				"id": 1,
				"fix_command": "command here",
				"fix_explanation": "One sentence why this fixes it"
			}
		]
		Note: For Gitleaks findings, the ID should correspond to the order in the secrets list but we will handle mapping. 
		Actually, for simplicity, provide one analysis per item if possible or a general report.
		The current system maps 'id' to Trivy vulnerabilities. 
		Please provide analysis for each vulnerability. For secrets, you can provide recommendations as well.`, string(vulnJSON), string(secretsJSON))

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
	cleanText := result.Text()
	// Strip markdown blocks if present
	cleanText = strings.TrimPrefix(cleanText, "```json")
	cleanText = strings.TrimPrefix(cleanText, "```")
	cleanText = strings.TrimSuffix(cleanText, "```")
	cleanText = strings.TrimSpace(cleanText)

	if err := json.Unmarshal([]byte(cleanText), &analysis); err != nil {
		fmt.Printf("AI Response: %s\n", result.Text())
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return analysis, nil
}