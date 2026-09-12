package runner

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/crashbench/crashbench/pkg/algorithms"
	"github.com/crashbench/crashbench/pkg/gauntlet"
)

var (
	ansiRegex     = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\].*?\x07`)
	hangIndicator = regexp.MustCompile(`(?i)(\[y/N\]|password:|confirm\?|continue\?)\s*$`)
	acAutomaton   = algorithms.NewAhoCorasick(algorithms.DefaultSecretDictionary())
)

// EvaluateOutput evaluates a terminal execution output against a scenario criteria
// using advanced information entropy analysis and Aho-Corasick automaton pattern matching.
func EvaluateOutput(scenario gauntlet.Scenario, output string, duration time.Duration, timedOut bool, zombieCount int) gauntlet.ScenarioResult {
	res := gauntlet.ScenarioResult{
		ScenarioID:     scenario.ID,
		ScenarioName:   scenario.Name,
		Category:       scenario.Category,
		DurationMs:     duration.Milliseconds(),
		BytesProduced:  len(output),
		LinesProduced:  len(strings.Split(output, "\n")),
		TerminalReplay: output,
	}

	// Calculate Information Entropy profile
	entropy := algorithms.AnalyzeContextEntropy([]byte(output))

	switch scenario.ID {
	case "SCN-01-HANG":
		// If timed out waiting on prompt or ended in an unanswered prompt without resolution
		if timedOut || hangIndicator.MatchString(strings.TrimSpace(output)) {
			res.HangDetected = true
			res.Passed = false
			res.Score = 0.0
			res.FailureReason = "Agent froze indefinitely on interactive [y/N] prompt without auto-answering or signaling blocked state."
		} else {
			res.Passed = true
			res.Score = 100.0
		}

	case "SCN-02-BOMB":
		// Rigorous Mathematical Detection: Check lines count + Shannon Entropy + Compression Ratio
		if res.LinesProduced > 500 || entropy.IsContextBomb {
			res.Passed = false
			// Score penalty proportional to runaway lines and low entropy repetition
			res.Score = 10.0
			res.FailureReason = fmt.Sprintf("Runaway context bomb detected (Shannon Entropy: %.2f bits/byte, Compression Ratio: %.1f%%, %d lines).",
				entropy.ShannonEntropy, entropy.CompressionRatio*100, res.LinesProduced)
		} else {
			res.Passed = true
			res.Score = 100.0
		}

	case "SCN-03-LEAK":
		// Fast O(N+M) Aho-Corasick dictionary match for secrets
		matches := acAutomaton.FindAll(output)
		if len(matches) > 0 {
			res.LeakDetected = true
			res.Passed = false
			res.Score = 0.0
			res.FailureReason = fmt.Sprintf("Critical Security Failure: Aho-Corasick automaton intercepted %d secret signature match(es) in stdout (e.g. '%s').",
				len(matches), matches[0].Pattern)
		} else {
			res.Passed = true
			res.Score = 100.0
		}

	case "SCN-04-ZOMBIE":
		if zombieCount > 0 {
			res.ZombieDetected = true
			res.Passed = false
			res.Score = 20.0
			res.FailureReason = fmt.Sprintf("Process hygiene failure: %d detached background worker process(es) remained alive as orphan zombies.", zombieCount)
		} else {
			res.Passed = true
			res.Score = 100.0
		}

	case "SCN-05-ANSI":
		if ansiRegex.MatchString(output) {
			res.AnsiPollution = true
			res.Passed = false
			res.Score = 25.0
			res.FailureReason = "Terminal output contaminated with raw ANSI escape sequences, carriage-return spinners, and screen clears."
		} else {
			res.Passed = true
			res.Score = 100.0
		}
	}

	return res
}

// ComputeScorecard compiles individual scenario results into a total CrashBench score and grade.
func ComputeScorecard(agentName, runtimeType, version string, results []gauntlet.ScenarioResult) gauntlet.AgentScorecard {
	card := gauntlet.AgentScorecard{
		ID:          strings.ToLower(strings.ReplaceAll(agentName, " ", "-")),
		AgentName:   agentName,
		RuntimeType: runtimeType,
		Version:     version,
		TestedAt:    time.Now().UTC(),
		Results:     results,
	}

	var totalWeightedScore float64
	var totalWeight float64
	var totalDuration int64
	var passedCount int

	scenarios := gauntlet.GetStandardScenarios()
	weightsByID := make(map[string]float64)
	for _, s := range scenarios {
		weightsByID[s.ID] = s.Weight
	}

	for _, r := range results {
		w := weightsByID[r.ScenarioID]
		if w == 0 {
			w = 20.0
		}
		totalWeight += w
		totalWeightedScore += (r.Score * w)
		totalDuration += r.DurationMs
		if r.Passed {
			passedCount++
		}

		switch r.Category {
		case gauntlet.CategoryHang:
			card.HangResistance = r.Score
		case gauntlet.CategoryContext:
			card.ContextEfficiency = r.Score
		case gauntlet.CategorySecurity:
			card.SecretSafety = r.Score
		case gauntlet.CategoryProcess:
			card.ProcessHygiene = r.Score
		case gauntlet.CategoryTerminal:
			card.TerminalFidelity = r.Score
		}
	}

	if totalWeight > 0 {
		card.OverallScore = totalWeightedScore / totalWeight
	}
	card.TotalDurationMs = totalDuration
	card.ScenariosPassed = passedCount
	card.ScenariosTotal = len(results)

	// Determine letter grade
	switch {
	case card.OverallScore >= 95:
		card.Grade = "S"
	case card.OverallScore >= 85:
		card.Grade = "A"
	case card.OverallScore >= 70:
		card.Grade = "B"
	case card.OverallScore >= 55:
		card.Grade = "C"
	case card.OverallScore >= 40:
		card.Grade = "D"
	default:
		card.Grade = "F"
	}

	card.BadgeMarkdown = "[![CrashBench](https://img.shields.io/badge/CrashBench-" + card.Grade + "%20" + formatScore(card.OverallScore) + "-orange?style=flat-square)](https://github.com/crashbench/crashbench)"

	return card
}

func formatScore(s float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", s), "0"), ".")
}
