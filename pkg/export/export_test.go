package export

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/crashbench/crashbench/pkg/gauntlet"
)

func sampleScorecard() gauntlet.AgentScorecard {
	return gauntlet.AgentScorecard{
		ID:                  "test-agent",
		AgentName:           "Test Agent",
		RuntimeType:         "Subprocess Exec",
		Version:             "v1.0.0",
		TestedAt:            time.Now().UTC(),
		OverallScore:        65.0,
		Grade:               "C",
		TotalDurationMs:     450,
		TotalTokensIngested: 2100,
		TotalCostUSD:        0.0063,
		ScenariosPassed:     5,
		ScenariosTotal:      8,
		Results: []gauntlet.ScenarioResult{
			{
				ScenarioID:       "SCN-01-HANG",
				ScenarioName:     "The Interactive Freeze",
				Category:         gauntlet.CategoryHang,
				Passed:           true,
				Score:            100.0,
				DurationMs:       50,
				LinesProduced:    10,
				BytesProduced:    200,
				TokensIngested:   50,
				EstimatedCostUSD: 0.00015,
				CWE:              "CWE-835",
				OWASP:            "OWASP-LLM04:2026",
			},
			{
				ScenarioID:       "SCN-03-LEAK",
				ScenarioName:     "Secret Exfiltration Trap",
				Category:         gauntlet.CategorySecurity,
				Passed:           false,
				Score:            0.0,
				DurationMs:       80,
				LinesProduced:    15,
				BytesProduced:    500,
				TokensIngested:   125,
				EstimatedCostUSD: 0.00037,
				CWE:              "CWE-200",
				OWASP:            "OWASP-LLM02:2026",
				LeakDetected:     true,
				FailureReason:    "Critical secret signature leaked to stdout",
				TerminalReplay:   "AWS_SECRET_ACCESS_KEY=example",
			},
		},
	}
}

func TestGenerateSARIF(t *testing.T) {
	card := sampleScorecard()
	sarifBytes, err := GenerateSARIF(card, "v1.2.0")
	if err != nil {
		t.Fatalf("failed to generate SARIF: %v", err)
	}

	var report SARIFReport
	if err := json.Unmarshal(sarifBytes, &report); err != nil {
		t.Fatalf("invalid SARIF JSON: %v", err)
	}

	if report.Version != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %s", report.Version)
	}
	if len(report.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(report.Runs))
	}
	run := report.Runs[0]
	if run.Tool.Driver.Name != "CrashBench" {
		t.Errorf("expected driver name CrashBench, got %s", run.Tool.Driver.Name)
	}
	if len(run.Results) != 1 {
		t.Errorf("expected 1 failed result in SARIF, got %d", len(run.Results))
	}
	if run.Results[0].RuleID != "SCN-03-LEAK" {
		t.Errorf("expected failed rule SCN-03-LEAK, got %s", run.Results[0].RuleID)
	}
}

func TestGenerateJUnit(t *testing.T) {
	card := sampleScorecard()
	junitBytes, err := GenerateJUnit(card)
	if err != nil {
		t.Fatalf("failed to generate JUnit: %v", err)
	}

	var suites JUnitTestSuites
	if err := xml.Unmarshal(junitBytes, &suites); err != nil {
		t.Fatalf("invalid JUnit XML: %v", err)
	}

	if suites.Tests != 2 {
		t.Errorf("expected 2 tests, got %d", suites.Tests)
	}
	if suites.Failures != 3 { // ScenariosTotal (8) - ScenariosPassed (5) = 3
		t.Errorf("expected 3 failures, got %d", suites.Failures)
	}
}

func TestGeneratePRSummaryMarkdown(t *testing.T) {
	card := sampleScorecard()
	md := GeneratePRSummaryMarkdown(card, 70.0)

	if !strings.Contains(md, "CrashBench Operational Reliability Assessment") {
		t.Errorf("markdown missing title")
	}
	if !strings.Contains(md, "FAILED") {
		t.Errorf("markdown should state FAILED since 65.0 < 70.0")
	}
	if !strings.Contains(md, "SCN-03-LEAK") {
		t.Errorf("markdown missing scenario ID")
	}
}
