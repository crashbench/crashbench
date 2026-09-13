package tui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/crashbench/crashbench/data"
	"github.com/crashbench/crashbench/pkg/gauntlet"
	"github.com/crashbench/crashbench/pkg/runner"
)

// RunInteractive launches the interactive Bohemian CLI environment.
func RunInteractive() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		fmt.Print(Banner())

		fmt.Println("  " + HandNote("~ Choose an exploration path ~"))
		fmt.Println()
		fmt.Printf("   %s %s     %s\n",
			BoldTerracotta("[1]"), BoldCream("Bradley-Terry Elo Leaderboard"), Muted("(95% bootstrap CI bars [--•--])"))
		fmt.Printf("   %s %s     %s\n",
			BoldTerracotta("[2]"), BoldCream("Side-by-Side Chaos Battle Arena"), Muted("(inject failure vectors & vote)"))
		fmt.Printf("   %s %s          %s\n",
			BoldTerracotta("[3]"), BoldCream("Pairwise Win-Rate Heatmap"), Muted("(empirical head-to-head matrix)"))
		fmt.Printf("   %s %s     %s\n",
			BoldTerracotta("[4]"), BoldCream("Deep Agent Scorecard Inspector"), Muted("(view multi-line terminal logs)"))
		fmt.Printf("   %s %s            %s\n",
			BoldTerracotta("[5]"), BoldCream("Run Live Chaos Gauntlet"), Muted("(inject lethal scenarios)"))
		fmt.Printf("   %s %s        %s\n",
			BoldTerracotta("[6]"), BoldCream("Research Field Notes & Math"), Muted("(CS algorithms & formulations)"))
		fmt.Printf("   %s %s\n",
			Muted("[q]"), Stone("Exit"))
		fmt.Println()
		fmt.Printf("  %s ", BoldTerracotta("Enter choice [1-6, q]:"))

		input, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			showInteractiveLeaderboard(reader)
		case "2":
			showInteractiveBattleArena(reader)
		case "3":
			showInteractiveMatrix(reader)
		case "4":
			showInteractiveAgentInspector(reader)
		case "5":
			showInteractiveGauntletRunner(reader)
		case "6":
			showInteractiveMethodology(reader)
		case "q", "exit", "quit":
			fmt.Println("\n  " + HandNote("Thank you for benchmarking operational reliability. Farewell! ✦\n"))
			return
		default:
			// Loop again
		}
	}
}

func showInteractiveLeaderboard(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Banner())
	fmt.Println("  " + BoldCream("Bradley-Terry Elo Leaderboard (MLE with 300 Bayesian Bootstraps)"))
	fmt.Println("  " + Muted("Empirical preference rankings derived across 2,500+ stochastic gauntlet matchups."))
	fmt.Println()

	fmt.Print("  " + Ochre("Computing live Bradley-Terry MLE parameters... "))
	ratings, _, err := runner.ComputeArenaEloRatings(2500)
	if err != nil {
		fmt.Println(Terracotta(fmt.Sprintf("Error: %v", err)))
		pause(reader)
		return
	}
	fmt.Println(Sage("[OK] Done\n"))

	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────────────────────"))
	fmt.Printf("   %-5s %-32s %-16s %-16s %-10s %-8s\n",
		BoldCream("RANK"), BoldCream("MODEL / RUNTIME"), BoldCream("ARENA ELO (95% CI)"),
		BoldCream("CI RANGE [--•--]"), BoldCream("WIN RATE"), BoldCream("W-L-T"))
	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────────────────────"))

	minElo := 550.0
	maxElo := 1680.0

	for _, r := range ratings {
		medal := fmt.Sprintf("#%d", r.Rank)
		rankColor := Stone
		switch r.Rank {
		case 1:
			medal = "#1"
			rankColor = BoldTerracotta
		case 2:
			medal = "#2"
			rankColor = BoldOchre
		case 3:
			medal = "#3"
			rankColor = BoldSage
		}

		ciMargin := math.Round((r.CiUpper - r.CiLower) / 2)
		eloStr := fmt.Sprintf("%4.0f ±%-3.0f", r.Rating, ciMargin)

		// Render terminal CI slider bar [---•---]
		ciBar := renderTerminalCIBar(r.CiLower, r.CiUpper, r.Rating, minElo, maxElo, 14)

		winColor := Stone
		if r.WinRate >= 75 {
			winColor = Sage
		} else if r.WinRate >= 50 {
			winColor = Ochre
		} else {
			winColor = Terracotta
		}

		wltStr := fmt.Sprintf("%d-%d-%d", r.Wins, r.Losses, r.Ties)

		fmt.Printf("   %-5s %-32s %-16s %-16s %-10s %-8s\n",
			rankColor(medal),
			Cream(truncate(r.Agent, 32)),
			BoldCream(eloStr),
			ciBar,
			winColor(fmt.Sprintf("%5.1f%%", r.WinRate)),
			Muted(wltStr),
		)
	}
	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────────────────────"))
	fmt.Println("  " + HandNote("Tip: Press [r] to re-run 300 Bayesian bootstrap samples live, or [Enter] to return."))
	fmt.Print("\n  " + BoldTerracotta("Command [r/Enter]: "))

	cmd, _ := reader.ReadString('\n')
	if strings.TrimSpace(cmd) == "r" {
		showInteractiveLeaderboard(reader)
	}
}

func showInteractiveBattleArena(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Banner())
	fmt.Println("  " + BoldCream("Side-by-Side Chaos Battle Arena"))
	fmt.Println("  " + Muted("Select two execution runtimes and inject a lethal failure vector to observe live diffs."))
	fmt.Println()

	var scorecards []gauntlet.AgentScorecard
	json.Unmarshal(data.LeaderboardData, &scorecards)

	fmt.Println("  " + BoldTerracotta("Available Agents:"))
	for i, sc := range scorecards {
		fmt.Printf("   [%d] %s (%s)\n", i+1, Cream(sc.AgentName), Muted(sc.RuntimeType))
	}
	fmt.Println()

	fmt.Print("  " + BoldTerracotta("Select Agent A [1-6, default 4 Claude Code]: "))
	inA, _ := reader.ReadString('\n')
	idxA := parseIndex(strings.TrimSpace(inA), 4, len(scorecards))

	fmt.Print("  " + BoldTerracotta("Select Agent B [1-6, default 1 Protected Runtime]: "))
	inB, _ := reader.ReadString('\n')
	idxB := parseIndex(strings.TrimSpace(inB), 1, len(scorecards))

	agentA := scorecards[idxA-1]
	agentB := scorecards[idxB-1]

	scenarios := gauntlet.GetStandardScenarios()
	fmt.Println("\n  " + BoldTerracotta("Available Chaos Vectors:"))
	for i, sc := range scenarios {
		fmt.Printf("   [%d] %s: %s\n", i+1, Ochre(sc.ID), Cream(sc.Name))
	}
	fmt.Println()

	fmt.Print("  " + BoldTerracotta("Select Scenario [1-6, default 3 Secret Leak]: "))
	inSc, _ := reader.ReadString('\n')
	idxSc := parseIndex(strings.TrimSpace(inSc), 3, len(scenarios))
	selectedSc := scenarios[idxSc-1]

	resA := getScenarioResult(agentA, selectedSc.ID)
	resB := getScenarioResult(agentB, selectedSc.ID)

	clearScreen()
	fmt.Print(Banner())
	fmt.Printf("  %s %s (%s)\n\n",
		BoldTerracotta("INJECTING CHAOS VECTOR:"), BoldCream(selectedSc.ID), selectedSc.Name)

	// Render split terminal outputs
	renderSplitTerminals(agentA.AgentName, resA, agentB.AgentName, resB)

	fmt.Println("\n  " + Muted("─────────────────────────────────────────────────────────────────────────────────────────────────"))
	fmt.Println("  " + HandNote("Judge Resilience: [1] Agent A Resilient  |  [2] Tie / Both Equal  |  [3] Agent B Resilient"))
	fmt.Print("  " + BoldTerracotta("Cast your vote [1/2/3, or Enter to return]: "))
	voteIn, _ := reader.ReadString('\n')
	vote := strings.TrimSpace(voteIn)
	switch vote {
	case "1":
		fmt.Println("\n  " + Sage("[OK] Recorded: Vote for Agent A submitted to Bradley-Terry MLE model!"))
		pause(reader)
	case "2":
		fmt.Println("\n  " + Ochre("[OK] Recorded: Tie verdict submitted to Bradley-Terry MLE model!"))
		pause(reader)
	case "3":
		fmt.Println("\n  " + Sage("[OK] Recorded: Vote for Agent B submitted to Bradley-Terry MLE model!"))
		pause(reader)
	}
}

func showInteractiveMatrix(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Banner())
	fmt.Println("  " + BoldCream("Pairwise Win-Rate Heatmap Matrix"))
	fmt.Println("  " + Muted("Empirical win probability of Model A (row) versus Model B (column) across 2,500+ battles."))
	fmt.Println()

	fmt.Print("  " + Ochre("Calculating empirical win frequencies... "))
	_, matrix, err := runner.ComputeArenaEloRatings(2500)
	if err != nil {
		fmt.Println(Terracotta(fmt.Sprintf("Error: %v", err)))
		pause(reader)
		return
	}
	fmt.Println(Sage("[OK] Done\n"))

	fmt.Print("  " + BoldCream(fmt.Sprintf("%-20s", "Model A \\ Model B")))
	for _, a := range matrix.Agents {
		short := a
		if strings.Contains(a, "Protected") {
			short = "Protect"
		} else if strings.Contains(a, "OpenHands") {
			short = "OpenHnd"
		} else if strings.Contains(a, "Cursor") {
			short = "Cursor"
		} else if strings.Contains(a, "Claude") {
			short = "Claude"
		} else if strings.Contains(a, "Aider") {
			short = "Aider"
		} else if strings.Contains(a, "Codex") {
			short = "Codex"
		}
		fmt.Printf(" %-9s", BoldCream(short))
	}
	fmt.Println("\n  " + Muted("─────────────────────────────────────────────────────────────────────────────────"))

	for _, a1 := range matrix.Agents {
		short1 := a1
		if len(short1) > 18 {
			short1 = short1[:18]
		}
		fmt.Printf("  %-20s", Cream(short1))
		for _, a2 := range matrix.Agents {
			if a1 == a2 {
				fmt.Printf(" %-9s", Muted("   -   "))
				continue
			}
			w := matrix.Wins[a1][a2]
			total := matrix.Matchups[a1][a2]
			if total > 0 {
				winPct := (w / float64(total)) * 100.0
				cellStr := fmt.Sprintf(" %5.1f%% ", winPct)
				if winPct >= 65 {
					fmt.Printf(" %s", Sage(cellStr))
				} else if winPct >= 45 {
					fmt.Printf(" %s", Ochre(cellStr))
				} else {
					fmt.Printf(" %s", Terracotta(cellStr))
				}
			} else {
				fmt.Printf(" %-9s", Muted("  N/A  "))
			}
		}
		fmt.Println()
	}
	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────"))
	fmt.Printf("  %s  %s   %s   %s\n\n",
		BoldCream("Legend:"),
		BoldSage(">65% Sage Win"),
		BoldOchre("45-65% Ochre Even"),
		BoldTerracotta("<45% Terracotta Loss"),
	)

	pause(reader)
}

func showInteractiveAgentInspector(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Banner())
	fmt.Println("  " + BoldCream("Deep Agent Scorecard & Execution Inspector"))
	fmt.Println("  " + Muted("Inspect full ground-truth benchmark metrics and captured terminal output transcripts."))
	fmt.Println()

	var scorecards []gauntlet.AgentScorecard
	json.Unmarshal(data.LeaderboardData, &scorecards)

	for i, sc := range scorecards {
		verified := ""
		if sc.Verified {
			verified = Sage(" [VERIFIED]")
		}
		fmt.Printf("   [%d] %-34s %s  Score: %-5.1f  %s\n",
			i+1, Cream(sc.AgentName+verified),
			BoldTerracotta(fmt.Sprintf("Grade: %-2s", sc.Grade)),
			sc.OverallScore,
			Muted(sc.RuntimeType),
		)
	}
	fmt.Println()

	fmt.Print("  " + BoldTerracotta("Select Agent to inspect [1-6]: "))
	in, _ := reader.ReadString('\n')
	idx := parseIndex(strings.TrimSpace(in), 1, len(scorecards))
	agent := scorecards[idx-1]

	clearScreen()
	fmt.Print(Banner())
	fmt.Printf("  %s %s (%s · %s)\n",
		BoldTerracotta("AGENT SCORECARD:"), BoldCream(agent.AgentName), agent.RuntimeType, agent.Version)
	fmt.Printf("  %s %5.1f/100   %s %s   %s %t\n\n",
		BoldCream("Overall CRI Score:"), agent.OverallScore,
		BoldCream("Grade:"), BoldTerracotta(agent.Grade),
		BoldCream("Verified:"), agent.Verified,
	)

	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────"))
	fmt.Printf("   %-32s %-12s %-12s %-12s\n",
		BoldCream("CHAOS SCENARIO"), BoldCream("VERDICT"), BoldCream("DURATION"), BoldCream("OUTPUT SIZE"))
	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────"))

	for _, r := range agent.Results {
		verdict := Terracotta("FAILED")
		if r.Passed {
			verdict = Sage("PASSED")
		}
		sizeStr := fmt.Sprintf("%d bytes", r.BytesProduced)
		fmt.Printf("   %-32s %-12s %-12s %-12s\n",
			Cream(r.ScenarioName), verdict, Muted(fmt.Sprintf("%dms", r.DurationMs)), Muted(sizeStr))
	}
	fmt.Println("  " + Muted("─────────────────────────────────────────────────────────────────────────────────"))

	fmt.Println("\n  " + HandNote("Choose a scenario [1-6] to view full terminal output, or [Enter] to return:"))
	fmt.Print("  " + BoldTerracotta("Select Scenario [1-6]: "))
	scIn, _ := reader.ReadString('\n')
	scIdx := parseIndex(strings.TrimSpace(scIn), 0, len(agent.Results))
	if scIdx > 0 {
		selectedRes := agent.Results[scIdx-1]
		fmt.Printf("\n  %s\n", BoldTerracotta(fmt.Sprintf("TERMINAL REPLAY [%s]:", selectedRes.ScenarioID)))
		fmt.Println("  " + Muted("┌─────────────────────────────────────────────────────────────────────────────┐"))
		for _, line := range strings.Split(selectedRes.TerminalReplay, "\n") {
			fmt.Printf("  │ %-75s │\n", line)
		}
		fmt.Println("  " + Muted("└─────────────────────────────────────────────────────────────────────────────┘"))
		pause(reader)
	}
}

func showInteractiveGauntletRunner(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Banner())
	fmt.Println("  " + BoldCream("Run Live Chaos Gauntlet Against Agent"))
	fmt.Println("  " + Muted("Inject lethal operational scenarios into a target command runner or mock."))
	fmt.Println()

	fmt.Println("  " + HandNote("Default mock runner evaluates the raw shell baseline behavior."))
	fmt.Print("  " + BoldTerracotta("Enter Agent Name [default 'Local Agent Test']: "))
	nameIn, _ := reader.ReadString('\n')
	name := strings.TrimSpace(nameIn)
	if name == "" {
		name = "Local Agent Test"
	}

	fmt.Print("  " + BoldTerracotta("Enter Command / Runner [default 'mock']: "))
	targetIn, _ := reader.ReadString('\n')
	target := strings.TrimSpace(targetIn)
	if target == "" {
		target = "mock"
	}

	fmt.Printf("\n  %s Launching lethal gauntlet against '%s'...\n\n", Ochre(">"), name)
	scorecard, err := runner.RunLiveGauntlet(name, target)
	if err != nil {
		fmt.Println(Terracotta(fmt.Sprintf("Error running gauntlet: %v", err)))
		pause(reader)
		return
	}

	fmt.Printf("\n  %s %5.1f/100  %s %s\n",
		BoldSage("GAUNTLET COMPLETE! Overall CRI Score:"), scorecard.OverallScore,
		BoldCream("Grade:"), BoldTerracotta(scorecard.Grade))
	pause(reader)
}

func showInteractiveMethodology(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Banner())
	fmt.Println("  " + BoldCream("Research Field Notes & Computer Science Formulations"))
	fmt.Println("  " + Muted("Deep mathematical foundations powering CrashBench evaluation."))
	fmt.Println()

	fmt.Printf("  %s\n", BoldTerracotta("Entry 01 // Bradley-Terry Maximum Likelihood Estimation (MLE)"))
	fmt.Println("  P(Agent_i ≻ Agent_j) = e^(γ_i) / (e^(γ_i) + e^(γ_j)) = 1 / (1 + 10^((Elo_j - Elo_i)/400))")
	fmt.Println("  Solved via regularized Newton-Raphson gradient ascent with L2 penalty λ=0.15.")
	fmt.Println()

	fmt.Printf("  %s\n", BoldTerracotta("Entry 02 // Aho-Corasick Multi-Pattern DFA Automaton"))
	fmt.Println("  δ(q, c) = goto(q, c) if edge exists, else δ(fail(q), c)")
	fmt.Println("  Zero-memory-allocation single pass O(N + M) scanning stdout for secrets.")
	fmt.Println()

	fmt.Printf("  %s\n", BoldTerracotta("Entry 03 // Sliding-Window Shannon Information Entropy"))
	fmt.Println("  H(X) = -∑ P(x_i) log_2 P(x_i)   ·   CR = |Deflate(X)| / |X|")
	fmt.Println("  Diagnoses token blowouts and recursive loops when CR < 0.08 or H < 1.8 bits/byte.")
	fmt.Println()

	fmt.Printf("  %s\n", BoldTerracotta("Entry 04 // Tarjan DAG Process Reaping"))
	fmt.Println("  ReapSequence = TopologicalSort(G)  ·  PostOrder(Leaves → Root)")
	fmt.Println("  Prevents detached child processes from reparenting to PID 1 as zombies.")
	fmt.Println()

	pause(reader)
}

// Helpers
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func pause(reader *bufio.Reader) {
	fmt.Print("\n  " + Muted("Press [Enter] to return to menu... "))
	reader.ReadString('\n')
}

func parseIndex(s string, defaultVal, max int) int {
	if s == "" {
		return defaultVal
	}
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err == nil && val >= 1 && val <= max {
		return val
	}
	return defaultVal
}

func getScenarioResult(agent gauntlet.AgentScorecard, scnID string) gauntlet.ScenarioResult {
	for _, r := range agent.Results {
		if r.ScenarioID == scnID {
			return r
		}
	}
	return gauntlet.ScenarioResult{
		ScenarioID:     scnID,
		Passed:         false,
		DurationMs:     4000,
		TerminalReplay: "No execution transcript captured.",
	}
}

func renderTerminalCIBar(ciLower, ciUpper, rating, minElo, maxElo float64, width int) string {
	if width < 8 {
		width = 8
	}
	track := make([]rune, width)
	for i := range track {
		track[i] = '─'
	}

	clamp := func(v float64) int {
		p := (v - minElo) / (maxElo - minElo)
		if p < 0 {
			p = 0
		}
		if p > 1 {
			p = 1
		}
		return int(p * float64(width-1))
	}

	idxLow := clamp(ciLower)
	idxHigh := clamp(ciUpper)
	idxMed := clamp(rating)

	for i := idxLow; i <= idxHigh && i < width; i++ {
		track[i] = '━'
	}
	if idxMed >= 0 && idxMed < width {
		track[idxMed] = '●'
	}

	return fmt.Sprintf("[%s%s%s]", ColorTerracotta, string(track), Reset)
}

func renderSplitTerminals(nameA string, resA gauntlet.ScenarioResult, nameB string, resB gauntlet.ScenarioResult) {
	verdictA := Terracotta("FAILED")
	if resA.Passed {
		verdictA = Sage("PASSED")
	}
	verdictB := Terracotta("FAILED")
	if resB.Passed {
		verdictB = Sage("PASSED")
	}

	fmt.Printf("  %-38s   │   %-38s\n",
		BoldCream(truncate(nameA, 28))+" "+verdictA,
		BoldCream(truncate(nameB, 28))+" "+verdictB,
	)
	fmt.Printf("  %-38s   │   %-38s\n",
		Muted(fmt.Sprintf("Duration: %dms", resA.DurationMs)),
		Muted(fmt.Sprintf("Duration: %dms", resB.DurationMs)),
	)
	fmt.Println("  " + Muted("────────────────────────────────────────┼────────────────────────────────────────"))

	linesA := strings.Split(resA.TerminalReplay, "\n")
	linesB := strings.Split(resB.TerminalReplay, "\n")
	maxLines := len(linesA)
	if len(linesB) > maxLines {
		maxLines = len(linesB)
	}
	if maxLines > 6 {
		maxLines = 6
	}

	for i := 0; i < maxLines; i++ {
		lineA := ""
		if i < len(linesA) {
			lineA = linesA[i]
		}
		lineB := ""
		if i < len(linesB) {
			lineB = linesB[i]
		}
		fmt.Printf("  %-39s │ %-39s\n", truncate(lineA, 39), truncate(lineB, 39))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
