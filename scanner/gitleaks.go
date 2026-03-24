package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type GitleaksFinding struct {
	Description string `json:"Description"`
	StartLine   int    `json:"StartLine"`
	EndLine     int    `json:"EndLine"`
	File        string `json:"File"`
	RuleID      string `json:"RuleID"`
	Secret      string `json:"Secret"`
	Match       string `json:"Match"`
}

func RunGitleaks(target string) ([]GitleaksFinding, error) {
	reportFile := "gitleaks-report.json"
	// Run gitleaks detect on the target directory
	// Note: using --no-git to avoid git ownership issues (e.g. dubious ownership in WSL/mounted drives)
	cmd := exec.Command("gitleaks", "detect", "--source", target, "--report-format", "json", "--report-path", reportFile, "--exit-code", "0", "--no-git")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gitleaks error: %v, output: %s", err, string(output))
	}

	// Check if report file exists
	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		// No leaks found or error
		return []GitleaksFinding{}, nil
	}

	data, err := os.ReadFile(reportFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read gitleaks report: %w", err)
	}

	var findings []GitleaksFinding
	if err := json.Unmarshal(data, &findings); err != nil {
		return nil, fmt.Errorf("failed to parse gitleaks report: %w", err)
	}

	// Clean up
	_ = os.Remove(reportFile)

	return findings, nil
}
