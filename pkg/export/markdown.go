package export

import (
	"fmt"
	"os"
	"strings"

	"github.com/crashbench/crashbench/pkg/gauntlet"
)

// GeneratePRSummaryMarkdown generates a GitHub Pull Request markdown comment.
func GeneratePRSummaryMarkdown(scorecard gauntlet.AgentScorecard, failUnder float64) string {
	var b strings.Builder

	statusBadge := "PASSED"
	statusColor := "green"
	if scorecard.OverallScore < failUnder {
		statusBadge = "FAILED"
		statusColor = "red"
	}

	b.WriteString("## CrashBench Operational Reliability Assessment\n\n")
	b.WriteString(fmt.Sprintf("> **Status:** `%s` (Threshold: %.1f | Score: **%.1f/100** | Grade: **%s**)\n\n",
		statusBadge, failUnder, scorecard.OverallScore, scorecard.Grade))

	b.WriteString("### Summary Telemetry\n\n")
	b.WriteString("| Metric | Value |\n")
	b.WriteString("| :--- | :--- |\n")
	b.WriteString(fmt.Sprintf("| **Agent Under Test** | `%s` (%s) |\n", scorecard.AgentName, scorecard.RuntimeType))
	b.WriteString(fmt.Sprintf("| **Resilience Index (CRI)** | **%.1f / 100** (Grade `%s`) |\n", scorecard.OverallScore, scorecard.Grade))
	b.WriteString(fmt.Sprintf("| **Scenarios Passed** | %d / %d (%.0f%%) |\n", scorecard.ScenariosPassed, scorecard.ScenariosTotal, (float64(scorecard.ScenariosPassed)/float64(scorecard.ScenariosTotal))*100))
	b.WriteString(fmt.Sprintf("| **Context Ingestion Bloat** | %d tokens (%d bytes) |\n", scorecard.TotalTokensIngested, scorecard.TotalTokensIngested*4))
	b.WriteString(fmt.Sprintf("| **Est. API Runaway Cost** | $%.4f USD |\n", scorecard.TotalCostUSD))
	b.WriteString(fmt.Sprintf("| **Total Execution Duration** | %dms |\n\n", scorecard.TotalDurationMs))

	b.WriteString("### Scenario Breakdown (CWE & OWASP GenAI Taxonomy)\n\n")
	b.WriteString("| Vector ID | Scenario Name | Category | Taxonomy | Verdict | Duration | Tokens | Cost |\n")
	b.WriteString("| :--- | :--- | :--- | :--- | :---: | :---: | :---: | :---: |\n")

	for _, r := range scorecard.Results {
		verdict := "PASS"
		if !r.Passed {
			verdict = "**FAIL**"
		}
		taxonomy := fmt.Sprintf("`%s` · `%s`", r.CWE, r.OWASP)
		b.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s | %dms | %d | $%.4f |\n",
			r.ScenarioID, r.ScenarioName, r.Category, taxonomy, verdict, r.DurationMs, r.TokensIngested, r.EstimatedCostUSD))
	}

	b.WriteString("\n---\n")
	b.WriteString("*Automated operational safety evaluation powered by [CrashBench](https://github.com/crashbench/crashbench).*\n")

	_ = statusColor // reserved for badge url

	return b.String()
}

// WritePRSummaryFile writes the PR summary markdown to file.
func WritePRSummaryFile(filePath string, scorecard gauntlet.AgentScorecard, failUnder float64) error {
	content := GeneratePRSummaryMarkdown(scorecard, failUnder)
	return os.WriteFile(filePath, []byte(content), 0644)
}
