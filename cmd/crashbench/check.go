package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/crashbench/crashbench/pkg/export"
	"github.com/crashbench/crashbench/pkg/runner"
	"github.com/crashbench/crashbench/pkg/tui"
)

func executeCheck(args []string) {
	cmd := flag.NewFlagSet("check", flag.ExitOnError)
	target := cmd.String("target", "mock", "Target agent or command runner to benchmark (default: 'mock')")
	name := cmd.String("name", "Agent CI Gate", "Name of the agent or runtime under test")
	failUnder := cmd.Float64("fail-under", 70.0, "Minimum required Resilience Index (CRI) score (0-100). Exits with code 1 if failed.")
	sarifPath := cmd.String("sarif", "", "Export results as SARIF 2.1.0 for GitHub Code Scanning alerts")
	junitPath := cmd.String("junit", "", "Export results as JUnit XML for CI/CD test dashboards")
	summaryMdPath := cmd.String("summary-md", "", "Export formatted PR comment markdown table")
	outputPath := cmd.String("output", "", "Export full scorecard JSON")

	cmd.Parse(args)

	fmt.Print(tui.Banner())
	fmt.Printf("\nExecuting Operational Chaos Gate for '%s' (Fail threshold: %.1f)...\n\n", *name, *failUnder)

	scorecard, err := runner.RunLiveGauntlet(*name, *target)
	if err != nil {
		fmt.Printf("Error running chaos gate: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=========================================================================================")
	fmt.Printf(" CRASHBENCH RESILIENCE INDEX: %.1f / 100   GRADE: %s   STATUS: ", scorecard.OverallScore, scorecard.Grade)
	if scorecard.OverallScore >= *failUnder {
		fmt.Printf("[PASS]\n")
	} else {
		fmt.Printf("[FAIL]\n")
	}
	fmt.Println("-----------------------------------------------------------------------------------------")
	fmt.Printf(" Scenarios Passed:      %d / %d (%.0f%%)\n", scorecard.ScenariosPassed, scorecard.ScenariosTotal, (float64(scorecard.ScenariosPassed)/float64(scorecard.ScenariosTotal))*100)
	fmt.Printf(" Context Token Bloat:   %d tokens (%d bytes)\n", scorecard.TotalTokensIngested, scorecard.TotalTokensIngested*4)
	fmt.Printf(" Estimated API Waste:   $%.4f USD (@ $3.00/1M tokens)\n", scorecard.TotalCostUSD)
	fmt.Printf(" Total Wall-Clock Time: %dms\n", scorecard.TotalDurationMs)
	fmt.Println("=========================================================================================")

	// Export SARIF 2.1.0
	if *sarifPath != "" {
		if err := export.WriteSARIFFile(*sarifPath, scorecard, Version); err != nil {
			fmt.Printf("Error writing SARIF file %s: %v\n", *sarifPath, err)
		} else {
			fmt.Printf("[OK] Exported SARIF 2.1.0 to: %s\n", *sarifPath)
		}
	}

	// Export JUnit XML
	if *junitPath != "" {
		if err := export.WriteJUnitFile(*junitPath, scorecard); err != nil {
			fmt.Printf("Error writing JUnit file %s: %v\n", *junitPath, err)
		} else {
			fmt.Printf("[OK] Exported JUnit XML to: %s\n", *junitPath)
		}
	}

	// Export Summary Markdown
	if *summaryMdPath != "" {
		if err := export.WritePRSummaryFile(*summaryMdPath, scorecard, *failUnder); err != nil {
			fmt.Printf("Error writing PR markdown summary %s: %v\n", *summaryMdPath, err)
		} else {
			fmt.Printf("[OK] Exported PR comment markdown to: %s\n", *summaryMdPath)
		}
	}

	// Export Scorecard JSON
	if *outputPath != "" {
		dataBytes, err := json.MarshalIndent(scorecard, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding scorecard JSON: %v\n", err)
		} else {
			if err := os.WriteFile(*outputPath, dataBytes, 0644); err != nil {
				fmt.Printf("Error writing scorecard JSON %s: %v\n", *outputPath, err)
			} else {
				fmt.Printf("[OK] Exported Scorecard JSON to: %s\n", *outputPath)
			}
		}
	}

	// Fail-under enforcement
	if scorecard.OverallScore < *failUnder {
		fmt.Printf("\n[FAIL] Operational Resilience score (%.1f) failed required threshold (%.1f). Gating failed.\n",
			scorecard.OverallScore, *failUnder)
		os.Exit(1)
	}

	fmt.Printf("\n[PASS] Operational Resilience score (%.1f) satisfied threshold (%.1f). Gating passed.\n",
		scorecard.OverallScore, *failUnder)
}

func executeInitCI() {
	workflowDir := ".github/workflows"
	workflowPath := filepath.Join(workflowDir, "crashbench.yml")

	workflowContent := `name: CrashBench Agent Reliability Gate

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  crashbench-gate:
    name: Operational Chaos Gate
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '>=1.22'

      - name: Install CrashBench
        run: |
          go install github.com/crashbench/crashbench/cmd/crashbench@latest

      - name: Run Operational Reliability Check
        run: |
          crashbench check \
            --target "mock" \
            --name "CI-Agent" \
            --fail-under 70 \
            --sarif results.sarif \
            --junit junit.xml \
            --summary-md pr-comment.md

      - name: Upload SARIF to GitHub Code Scanning Alerts
        if: always()
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif
`

	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		fmt.Printf("Error creating %s directory: %v\n", workflowDir, err)
		os.Exit(1)
	}

	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", workflowPath, err)
		os.Exit(1)
	}

	fmt.Print(tui.Banner())
	fmt.Printf("\n[OK] Successfully initialized GitHub Actions workflow at: %s\n", workflowPath)
	fmt.Println("     This workflow runs 'crashbench check' on push/PR and uploads SARIF 2.1.0 security alerts.")
}
