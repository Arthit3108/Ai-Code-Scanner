![Release](https://img.shields.io/github/v/release/Arthit3108/Ai-Code-Scanner)
![Go](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Trivy](https://img.shields.io/badge/Trivy-supported-00C7B7)
![Gitleaks](https://img.shields.io/badge/Gitleaks-supported-orange)
![Gemini](https://img.shields.io/badge/Google%20Gemini-AI%20powered-4285F4)

# 🛡️ AI Code Scanner

**AI-powered code vulnerability scanner** that combines [Trivy](https://trivy.dev/) + [Gitleaks](https://gitleaks.io/) with [Google Gemini AI](https://ai.google.dev/) to scan for vulnerabilities, detect leaked secrets, and **automatically suggest fixes** using AI.

> _"It doesn't just find vulnerabilities — AI tells you how to fix them."_

---

## 📋 Table of Contents

- [✨ Features](#-features)
- [🏗️ Architecture](#️-architecture)
- [📦 Quick Start](#-quick-start)
- [⚙️ GitHub Action](#️-github-action)
- [🖥️ CLI Usage](#️-cli-usage)
- [📄 Output Format](#-output-format)
- [🔐 SARIF + GitHub Code Scanning](#-sarif--github-code-scanning)
- [🏗️ Jenkins Integration](#️-jenkins-integration)
- [🔧 Development](#-development)
- [🤝 Contributing](#-contributing)
- [📝 License](#-license)

---

## ✨ Features

| Feature | Description |
|---|---|
| 🔍 **Trivy Scanning** | Scans for known vulnerabilities (CVEs) across dependencies, packages, and the filesystem |
| 🔑 **Gitleaks Detection** | Detects leaked secrets in source code such as API keys, passwords, and tokens |
| 🤖 **AI Fix Suggestions** | Uses Google Gemini to analyze findings and generate `fix_command` + `fix_explanation` for each vulnerability |
| 🔌 **GitHub Action** | Drop into any CI/CD pipeline with a single step |
| 📊 **SARIF Support** | Upload scan results to GitHub's Security tab (Code Scanning) |
| 🚦 **Severity Thresholds** | Automatically fail pipelines when CRITICAL/HIGH vulnerabilities are found |
| 🏭 **Jenkins Support** | Ready-to-use Jenkinsfile for Jenkins CI/CD pipelines |
| 📦 **Auto Release** | Tag push triggers automatic build and release of the Linux binary via GitHub Actions |

---

## 🏗️ Architecture

```
ai-code-scanner/
├── main.go                      # Entry point → calls cmd.Execute()
├── cmd/
│   ├── root.go                  # Cobra root command (ai-code-scanner)
│   └── scan.go                  # Scan command — orchestrates Trivy + Gitleaks + Gemini AI
├── scanner/
│   ├── trivy.go                 # Trivy integration — scans CVEs, parses results into CleanVuln
│   └── gitleaks.go              # Gitleaks integration — detects leaked secrets in source code
├── ai/
│   └── gemini.go                # Google Gemini AI — analyzes vulnerabilities and suggests fixes
├── action.yml                   # GitHub Action definition (composite action)
├── Jenkinsfile                  # Jenkins pipeline definition
├── .gitleaks.toml               # Gitleaks allowlist configuration
├── .github/workflows/
│   └── release.yml              # Auto-build & release binary on tag push
├── go.mod / go.sum              # Go module dependencies
└── .env                         # Environment variables (GEMINI_API_KEY)
```

### 🔄 Scan Flow

```
┌─────────────┐     ┌──────────────┐     ┌──────────────────┐
│  CLI / CI   │────▶│  Trivy Scan  │────▶│  Parse CVEs      │──┐
│  (scan cmd) │     │  (filesystem)│     │  → CleanVuln[]   │  │
└─────────────┘     └──────────────┘     └──────────────────┘  │
       │                                                        │
       │            ┌──────────────┐     ┌──────────────────┐  │
       └───────────▶│  Gitleaks    │────▶│  Parse Secrets   │──┤
                    │  (detect)    │     │  → Finding[]     │  │
                    └──────────────┘     └──────────────────┘  │
                                                               │
                    ┌──────────────────────────────────────────┘
                    ▼
            ┌───────────────┐     ┌──────────────────────┐
            │  Gemini AI    │────▶│  fix_command +        │
            │  Analysis     │     │  fix_explanation      │
            └───────────────┘     │  mapped back to       │
                                  │  each vulnerability   │
                                  └──────────┬───────────┘
                                             ▼
                                  ┌──────────────────────┐
                                  │  Combined Output     │
                                  │  (JSON / SARIF)      │
                                  └──────────────────────┘
```

---

## 📦 Quick Start

### As a GitHub Action (Easiest)

Create `.github/workflows/security-scan.yml` in your repository:

```yaml
name: Security Scan

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    permissions:
      security-events: write  # Required for SARIF upload
      actions: read
      contents: read

    steps:
      - uses: actions/checkout@v4

      - uses: Arthit3108/ai-code-scanner@v1
        with:
          gemini_api_key: ${{ secrets.GEMINI_API_KEY }}
```

### As a CLI Tool

```bash
# Prerequisites: Go 1.24+, Trivy, Gitleaks
go build -o ai-scanner .

# Set Gemini API Key
export GEMINI_API_KEY="your-api-key"

# Run a scan
./ai-scanner scan --target ./my-project
```

---

## ⚙️ GitHub Action

### Inputs

| Input | Description | Required | Default |
|---|---|---|---|
| `gemini_api_key` | Google Gemini API Key | ✅ Yes | — |
| `target` | Target directory to scan | No | `.` |
| `severity` | Severity levels to scan (e.g. `HIGH,CRITICAL`) | No | `HIGH,CRITICAL` |
| `output_type` | Output format: `json` / `sarif` | No | `json` |
| `fail_on_critical` | Fail pipeline on CRITICAL vulnerabilities | No | `true` |
| `fail_on_high` | Fail pipeline on HIGH vulnerabilities | No | `true` |
| `fail_on_medium` | Fail pipeline on MEDIUM vulnerabilities | No | `false` |
| `fail_on_low` | Fail pipeline on LOW vulnerabilities | No | `false` |

### Full Usage Example

```yaml
- uses: Arthit3108/ai-code-scanner@v1
  with:
    gemini_api_key: ${{ secrets.GEMINI_API_KEY }}
    target: "."
    severity: "HIGH,CRITICAL,MEDIUM"
    output_type: json
    fail_on_critical: "true"
    fail_on_high: "true"
    fail_on_medium: "false"
    fail_on_low: "false"
```

### What the Action Does Internally

1. **Install Gitleaks** — Downloads and installs Gitleaks v8.18.2
2. **Install Trivy** — Downloads and installs the latest Trivy
3. **Download Scanner Binary** — Fetches the `ai-scanner` binary from GitHub Releases
4. **Run Scan** — Executes `ai-scanner scan` with the configured parameters
5. **Upload SARIF** *(if sarif output selected)* — Uploads results to GitHub Code Scanning
6. **Upload Report Artifact** — Stores the report as a workflow artifact
7. **Check Vulnerabilities** — Validates vulnerability counts against configured severity thresholds

---

## 🖥️ CLI Usage

### Basic Command

```bash
ai-scanner scan [flags]
```

### Flags

| Flag | Short | Description | Default |
|---|---|---|---|
| `--target` | `-t` | Target directory to scan | `.` |
| `--output` | `-o` | Output file path | `scanResult.json` |
| `--severity` | | Severity levels (comma-separated) | `HIGH,CRITICAL` |
| `--output-type` | | Output format: `json` / `sarif` | `json` |

### Examples

```bash
# Scan the current directory with default settings
./ai-scanner scan

# Scan a specific directory with custom severity levels
./ai-scanner scan -t ./my-project --severity "CRITICAL,HIGH,MEDIUM"

# Scan and output as SARIF format
./ai-scanner scan --output-type sarif -o results.sarif

# Scan with a custom output path
./ai-scanner scan -t ./api-server -o report.json
```

### Sample Terminal Output

```
output type json
Running Gitleaks scan...

---------------------------------------
Scan Results Summary:
Total Vulnerabilities (Trivy): 15
🔴 Critical: 2
🟠 High: 13
🟡 Medium: 0
🟢 Low: 0

Total Secrets (Gitleaks): 0
---------------------------------------
Successfully saved results to scanResult.json
```

---

## 📄 Output Format

### JSON Output

When using `--output-type json`, results are saved as a **Combined Result** containing both vulnerabilities and secrets:

```json
{
  "vulnerabilities": [
    {
      "id": 1,
      "cve": "CVE-2024-45337",
      "package": "golang.org/x/crypto",
      "installed_version": "v0.28.0",
      "fixed_version": "0.31.0",
      "severity": "CRITICAL",
      "title": "Misuse of ServerConfig.PublicKeyCallback may cause authorization bypass",
      "fix_command": "go get golang.org/x/crypto@v0.35.0",
      "fix_explanation": "Upgrading to v0.35.0 patches the authorization bypass in SSH."
    }
  ],
  "secrets": [
    {
      "Description": "Generic API Key",
      "StartLine": 10,
      "EndLine": 10,
      "File": "config.js",
      "RuleID": "generic-api-key",
      "Secret": "AKIAIOSFODNN7EXAMPLE",
      "Match": "api_key = 'AKIAIOSFODNN7EXAMPLE'"
    }
  ]
}
```

> 💡 The `fix_command` and `fix_explanation` fields are generated by **Gemini AI** — they are only present when AI analysis succeeds.

### SARIF Output

When using `--output-type sarif`, results are output in the SARIF format which can be uploaded directly to GitHub Code Scanning.

---

## 🔐 SARIF + GitHub Code Scanning

To display scan results in GitHub's **Security → Code Scanning** tab:

```yaml
jobs:
  scan:
    runs-on: ubuntu-latest
    permissions:
      security-events: write  # ⚠️ Required

    steps:
      - uses: actions/checkout@v4

      - uses: Arthit3108/ai-code-scanner@v1
        with:
          gemini_api_key: ${{ secrets.GEMINI_API_KEY }}
          output_type: sarif
```

> **Important:** You must set `permissions: security-events: write` at the job level.

---

## 🏗️ Jenkins Integration

### Prerequisites

| Requirement | Description |
|---|---|
| **Trivy** | Must be installed on the Jenkins node (the Jenkinsfile will attempt to install it if missing) |
| **jq** | Required for severity threshold checking |
| **Jenkins Credential** | Create a "Secret text" credential for the Gemini API Key |

### Setup

1. Go to **Manage Jenkins → Credentials → System → Global credentials**
2. Add a new credential:
   - **Kind**: Secret text
   - **Secret**: `your-gemini-api-key`
   - **ID**: `gemini-api-key`
3. Add the [`Jenkinsfile`](Jenkinsfile) from this repository to your project root
4. Create a "Pipeline" job and point it to your repository

### Jenkins Pipeline Stages

```
Install Prerequisites → Download Scanner → AI Security Scan → Check Severity Thresholds
```

| Stage | Description |
|---|---|
| **Install Prerequisites** | Installs Trivy and jq if not already present (supports Debian/RHEL) |
| **Download Scanner** | Downloads the `ai-scanner` binary from GitHub Releases |
| **AI Security Scan** | Runs `ai-scanner scan` and saves the results |
| **Check Severity Thresholds** | Checks CRITICAL/HIGH counts and fails the build if thresholds are exceeded |

> 📎 See full example: [Jenkinsfile](Jenkinsfile)

---

## 🔧 Development

### Prerequisites

| Tool | Version | Link |
|---|---|---|
| Go | 1.24+ | [golang.org](https://golang.org/dl/) |
| Trivy | Latest | [trivy.dev](https://trivy.dev/) |
| Gitleaks | 8.18+ | [gitleaks.io](https://gitleaks.io/) |
| Gemini API Key | — | [ai.google.dev](https://ai.google.dev/) |

### Setup & Build

```bash
# Clone the repository
git clone https://github.com/Arthit3108/Ai-Code-Scanner.git
cd Ai-Code-Scanner

# Create .env file
echo "GEMINI_API_KEY=your-key-here" > .env

# Build
go build -o ai-scanner .

# Run
./ai-scanner scan --target .
```

### Project Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI framework — commands, flags, and help text |
| `github.com/joho/godotenv` | Loads environment variables from `.env` file |
| `google.golang.org/genai` | Google Gemini AI SDK — generates fix suggestions |

### Key Source Files

| File | Description |
|---|---|
| [`main.go`](main.go) | Entry point — calls `cmd.Execute()` |
| [`cmd/root.go`](cmd/root.go) | Root command definition (`ai-code-scanner`) |
| [`cmd/scan.go`](cmd/scan.go) | Scan command — orchestrates scanners, AI analysis, and output |
| [`scanner/trivy.go`](scanner/trivy.go) | Trivy runner — parses Trivy results into `CleanVuln[]` |
| [`scanner/gitleaks.go`](scanner/gitleaks.go) | Gitleaks runner — detects leaked secrets |
| [`ai/gemini.go`](ai/gemini.go) | Gemini AI client — sends findings to AI and receives fix suggestions |
| [`action.yml`](action.yml) | GitHub Action composite action definition |
| [`Jenkinsfile`](Jenkinsfile) | Jenkins pipeline definition |
| [`.gitleaks.toml`](.gitleaks.toml) | Gitleaks allowlist — excludes specified files from scanning |

### Release Process

When you push a tag starting with `v` (e.g. `v1.0.0`):

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions will automatically:
1. Checkout the code
2. Set up Go 1.24
3. Build the Linux binary (`ai-scanner`)
4. Create a GitHub Release with the binary attached

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

MIT
