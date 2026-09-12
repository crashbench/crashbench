package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/crashbench/crashbench/pkg/gauntlet"
)

// LoadLeaderboard reads the ground-truth benchmark dataset from disk.
func LoadLeaderboard(dataPath string) ([]gauntlet.AgentScorecard, error) {
	if dataPath == "" {
		// Try default locations
		candidates := []string{
			"data/leaderboard.json",
			"../data/leaderboard.json",
			"../../data/leaderboard.json",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				dataPath = c
				break
			}
		}
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read leaderboard data at %s: %w", dataPath, err)
	}

	var leaderboard []gauntlet.AgentScorecard
	if err := json.Unmarshal(data, &leaderboard); err != nil {
		return nil, fmt.Errorf("failed to parse leaderboard json: %w", err)
	}

	return leaderboard, nil
}

// RunLiveGauntlet executes the 5 lethal scenarios against a target command runner.
func RunLiveGauntlet(agentName, targetRunner string) (gauntlet.AgentScorecard, error) {
	scenarios := gauntlet.GetStandardScenarios()
	results := make([]gauntlet.ScenarioResult, 0, len(scenarios))

	fmt.Println("=================================================================")
	fmt.Println("             CRASHBENCH - AGENT CHAOS & RELIABILITY GAUNTLET      ")
	fmt.Printf(" Target: %s  |  Runner: %s\n", agentName, targetRunner)
	fmt.Println("=================================================================")

	for i, sc := range scenarios {
		fmt.Printf("\n[%d/5] Running %s (%s)...\n", i+1, sc.ID, sc.Name)
		start := time.Now()

		output, timedOut, zombieCount, err := executeScenarioTest(sc, targetRunner)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("  ⚠️ Execution error: %v\n", err)
		}

		res := EvaluateOutput(sc, output, elapsed, timedOut, zombieCount)
		results = append(results, res)

		if res.Passed {
			fmt.Printf("  ✅ PASSED (Score: %.1f | %dms)\n", res.Score, res.DurationMs)
		} else {
			fmt.Printf("  ❌ FAILED (Score: %.1f | %dms)\n", res.Score, res.DurationMs)
			fmt.Printf("     Reason: %s\n", res.FailureReason)
		}
	}

	scorecard := ComputeScorecard(agentName, "Custom Evaluation", "live-run", results)
	fmt.Println("\n=================================================================")
	fmt.Printf(" GAUNTLET COMPLETED: Overall Resilience Score: %.1f/100 (Grade: %s)\n", scorecard.OverallScore, scorecard.Grade)
	fmt.Printf(" Passed: %d/%d scenarios in %dms\n", scorecard.ScenariosPassed, scorecard.ScenariosTotal, scorecard.TotalDurationMs)
	fmt.Printf(" Badge Markdown: %s\n", scorecard.BadgeMarkdown)
	fmt.Println("=================================================================")

	return scorecard, nil
}

func executeScenarioTest(sc gauntlet.Scenario, runnerCmd string) (output string, timedOut bool, zombieCount int, err error) {
	// If runnerCmd is empty or "mock", return canonical payloads to simulate raw shell behavior
	if runnerCmd == "" || runnerCmd == "mock" {
		payload := gauntlet.GenerateScenarioPayload(sc.ID)
		switch sc.ID {
		case "SCN-01-HANG":
			return payload, true, 0, nil
		case "SCN-02-BOMB":
			return payload, false, 0, nil
		case "SCN-03-LEAK":
			return payload, false, 0, nil
		case "SCN-04-ZOMBIE":
			return payload, false, 1, nil
		case "SCN-05-ANSI":
			return payload, false, 0, nil
		}
		return payload, false, 0, nil
	}

	// Create temporary scenario script
	tmpDir, err := os.MkdirTemp("", "crashbench-*")
	if err != nil {
		return "", false, 0, err
	}
	defer os.RemoveAll(tmpDir)

	scriptFile := filepath.Join(tmpDir, "test_payload.txt")
	payload := gauntlet.GenerateScenarioPayload(sc.ID)
	if err := os.WriteFile(scriptFile, []byte(payload), 0644); err != nil {
		return "", false, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), sc.Timeout)
	defer cancel()

	// Execute target runner with the payload
	var cmd *exec.Cmd
	parts := strings.Fields(runnerCmd)
	if len(parts) == 1 {
		cmd = exec.CommandContext(ctx, parts[0], scriptFile)
	} else {
		args := append(parts[1:], scriptFile)
		cmd = exec.CommandContext(ctx, parts[0], args...)
	}

	outBytes, err := cmd.CombinedOutput()
	outStr := string(outBytes)

	if ctx.Err() == context.DeadlineExceeded {
		timedOut = true
		outStr += "\n[CRASHBENCH_TIMEOUT_EXCEEDED]"
	}

	return outStr, timedOut, 0, err
}
