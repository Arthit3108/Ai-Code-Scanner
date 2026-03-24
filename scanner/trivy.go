package scanner

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type CleanVuln struct {
	ID               int    `json:"id"`
	CVE              string `json:"cve"`
	Package          string `json:"package"`
	InstalledVersion string `json:"installed_version"`
	FixedVersion     string `json:"fixed_version"`
	Severity         string `json:"severity"`
	Title            string `json:"title"`
	FixCommand       string `json:"fix_command"`
	FixExplanation   string `json:"fix_explanation"`
}

type TrivyReport struct {
	Results []TrivyResult `json:"Results"`
}

type TrivyResult struct {
	Target          string               `json:"Target"`
	Vulnerabilities []TrivyVulnerability `json:"Vulnerabilities"`
}

type TrivyVulnerability struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
	Severity         string `json:"Severity"`
	Title            string `json:"Title"`
}

func RunTrivy(target, outputType, severity string) ([]CleanVuln, []byte, error) {
	var cmd *exec.Cmd

	switch outputType {
	case "json":
		fmt.Printf("output type %s\n", outputType)
		cmd = exec.Command("trivy", "fs", target, "-f", "json", "--severity", severity)
	case "sarif":
		fmt.Printf("output type %s\n", outputType)
		cmd = exec.Command("trivy", "fs", target, "-f", "sarif", "--severity", severity)
	default:
		return nil, nil, fmt.Errorf("unsupported output type: %s", outputType)
	}

	if cmd == nil {
		return nil, nil, fmt.Errorf("command not initialized")
	}

	scanOutput, err := cmd.Output()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run trivy: %w", err)
	}

	var vulns []CleanVuln

	if outputType == "json" {
		var report TrivyReport
		if err := json.Unmarshal(scanOutput, &report); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal trivy output: %w", err)
		}

		for _, result := range report.Results {
			for _, vuln := range result.Vulnerabilities {
				clean := CleanVuln{
					ID:               len(vulns) + 1,
					CVE:              vuln.VulnerabilityID,
					Package:          vuln.PkgName,
					InstalledVersion: vuln.InstalledVersion,
					FixedVersion:     vuln.FixedVersion,
					Severity:         vuln.Severity,
					Title:            vuln.Title,
				}
				vulns = append(vulns, clean)
			}
		}
	}

	return vulns, scanOutput, nil
}