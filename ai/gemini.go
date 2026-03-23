package ai

import (
	"ai-code-scanner/scanner"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func Gemini(vuln []scanner.CleanVuln) {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: os.Getenv("GEMINI_API_KEY"),
    })
    if err != nil {
        log.Fatal(err)
    }

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
		]`, vuln)

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text(prompt),
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Text())
}