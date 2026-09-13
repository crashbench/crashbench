package gauntlet

import (
	"time"
)

// Category represents the failure mode domain being tested.
type Category string

const (
	CategoryHang       Category = "Hang Resistance"
	CategoryContext    Category = "Context Efficiency"
	CategorySecurity   Category = "Secret Safety"
	CategoryProcess    Category = "Process Hygiene"
	CategoryTerminal   Category = "Terminal Fidelity"
	CategoryBoundary   Category = "Boundary Isolation"
)

// Scenario defines a specific chaos stress test with industry taxonomy mapping.
type Scenario struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Category    Category      `json:"category"`
	Description string        `json:"description"`
	Impact      string        `json:"impact"`
	Timeout     time.Duration `json:"timeout"`
	Weight      float64       `json:"weight"` // Out of 100 total
	CWE         string        `json:"cwe"`    // e.g. "CWE-835", "CWE-200"
	OWASP       string        `json:"owasp"`  // e.g. "OWASP-LLM04:2026"
}

// ScenarioResult records the outcome of a single scenario execution.
type ScenarioResult struct {
	ScenarioID       string   `json:"scenario_id"`
	ScenarioName     string   `json:"scenario_name"`
	Category         Category `json:"category"`
	Passed           bool     `json:"passed"`
	Score            float64  `json:"score"` // 0.0 - 100.0
	DurationMs       int64    `json:"duration_ms"`
	LinesProduced    int      `json:"lines_produced"`
	BytesProduced    int      `json:"bytes_produced"`
	TokensIngested   int      `json:"tokens_ingested"`    // Context window token estimate (bytes / 4)
	EstimatedCostUSD float64  `json:"estimated_cost_usd"` // $3.00 / 1M input tokens
	CWE              string   `json:"cwe"`
	OWASP            string   `json:"owasp"`
	HangDetected     bool     `json:"hang_detected"`
	LeakDetected     bool     `json:"leak_detected"`
	ZombieDetected   bool     `json:"zombie_detected"`
	AnsiPollution    bool     `json:"ansi_pollution"`
	FailureReason    string   `json:"failure_reason,omitempty"`
	TerminalReplay   string   `json:"terminal_replay"`
}

// AgentScorecard represents the aggregate benchmark scorecard for an agent.
type AgentScorecard struct {
	ID                  string           `json:"id"`
	AgentName           string           `json:"agent_name"`
	RuntimeType         string           `json:"runtime_type"` // e.g., "Raw Shell", "Protected Runtime", "Custom Container"
	Version             string           `json:"version"`
	TestedAt            time.Time        `json:"tested_at"`
	OverallScore        float64          `json:"overall_score"` // 0.0 - 100.0
	Grade               string           `json:"grade"`         // S, A, B, C, D, F
	HangResistance      float64          `json:"hang_resistance"`
	ContextEfficiency   float64          `json:"context_efficiency"`
	SecretSafety        float64          `json:"secret_safety"`
	ProcessHygiene      float64          `json:"process_hygiene"`
	TerminalFidelity    float64          `json:"terminal_fidelity"`
	BoundaryIsolation   float64          `json:"boundary_isolation"`
	TotalDurationMs     int64            `json:"total_duration_ms"`
	TotalTokensIngested int              `json:"total_tokens_ingested"`
	TotalCostUSD        float64          `json:"total_cost_usd"`
	ScenariosPassed     int              `json:"scenarios_passed"`
	ScenariosTotal      int              `json:"scenarios_total"`
	Results             []ScenarioResult `json:"results"`
	BadgeMarkdown       string           `json:"badge_markdown"`
	Verified            bool             `json:"verified"`
}
