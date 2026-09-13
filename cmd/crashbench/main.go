
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/crashbench/crashbench/data"
	"github.com/crashbench/crashbench/pkg/gauntlet"
	"github.com/crashbench/crashbench/pkg/runner"
	"github.com/crashbench/crashbench/pkg/tui"
	"github.com/crashbench/crashbench/pkg/version"
)

// Version resolves dynamically from git describe, release tags, or CHANGELOG.md
var Version = version.Get()

func init() {
	tui.SetVersion(Version)
}

func main() {
	if len(os.Args) < 2 {
		tui.RunInteractive()
		return
	}

	command := os.Args[1]

	switch command {
	case "interactive", "ui", "tui":
		tui.RunInteractive()

	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		target := runCmd.String("target", "mock", "Target agent or command runner to benchmark (e.g. 'npx claude-code' or 'mock')")
		name := runCmd.String("name", "Custom Agent", "Name of the agent or runtime under test")
		outputFile := runCmd.String("output", "", "Optional path to export results JSON (e.g. scorecard.json)")
		runCmd.Parse(os.Args[2:])

		executeRun(*name, *target, *outputFile)

	case "elo", "arena":
		executeEloArenaCLI()

	case "matrix", "pairwise":
		executeMatrixCLI()

	case "leaderboard", "top":
		executeLeaderboardCLI()

	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.Int("port", 4040, "HTTP port to serve API on")
		serveCmd.Parse(os.Args[2:])

		executeServe(*port)

	case "version", "-v", "--version":
		fmt.Printf("CrashBench %s (Bradley-Terry MLE Engine · Go)\n", Version)

	case "help", "-h", "--help":
		printHelp()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(tui.Banner())
	fmt.Println("\nUSAGE:")
	fmt.Println("  crashbench run [flags]         Run the 5-scenario lethal chaos gauntlet against an agent")
	fmt.Println("  crashbench elo                 Display LMSYS-style Bradley-Terry Elo ratings & 95% CIs")
	fmt.Println("  crashbench matrix              Display Head-to-Head Pairwise Win Rate matrix")
	fmt.Println("  crashbench leaderboard         Display the standard CrashBench Resilience Index (CRI)")
	fmt.Println("  crashbench serve [flags]       Start API server on localhost")
	fmt.Println("  crashbench version             Print CrashBench version")
	fmt.Println("\nFLAGS FOR 'run':")
	fmt.Println("  --target <cmd>   Target command to test (default: 'mock' for raw shell baseline)")
	fmt.Println("  --name <name>    Agent name (default: 'Custom Agent')")
	fmt.Println("  --output <file>  Export JSON scorecard to file")
	fmt.Println("\nEXAMPLES:")
	fmt.Println("  crashbench run --target 'npx claude-code' --name 'Claude Code'")
	fmt.Println("  crashbench elo")
	fmt.Println("  crashbench matrix")
}

func executeRun(name, target, outputFile string) {
	fmt.Print(tui.Banner())
	scorecard, err := runner.RunLiveGauntlet(name, target)
	if err != nil {
		fmt.Printf("Error running gauntlet: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "" {
		dataBytes, err := json.MarshalIndent(scorecard, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding output JSON: %v\n", err)
		} else {
			if err := os.WriteFile(outputFile, dataBytes, 0644); err != nil {
				fmt.Printf("Error writing output file %s: %v\n", outputFile, err)
			} else {
				fmt.Printf("Successfully exported scorecard to: %s\n", outputFile)
			}
		}
	}
}

func executeEloArenaCLI() {
	fmt.Print(tui.Banner())
	fmt.Println("\nComputing Bradley-Terry MLE Elo Ratings (with 300 Bayesian Bootstraps)...")
	ratings, _, err := runner.ComputeArenaEloRatings(2500)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\n=========================================================================================")
	fmt.Printf(" %-5s %-32s %-16s %-12s %-12s %-8s\n", "RANK", "MODEL / RUNTIME", "ARENA ELO (95% CI)", "WIN RATE", "MATCHUPS", "W-L-T")
	fmt.Println("-----------------------------------------------------------------------------------------")

	for _, r := range ratings {
		ciStr := fmt.Sprintf("%.0f (%.0f-%.0f)", r.Rating, r.CiLower, r.CiUpper)
		wltStr := fmt.Sprintf("%d-%d-%d", r.Wins, r.Losses, r.Ties)
		fmt.Printf(" #%-4d %-32s %-18s %-10s %-10d %-8s\n",
			r.Rank,
			truncate(r.Agent, 32),
			ciStr,
			fmt.Sprintf("%.1f%%", r.WinRate),
			r.TotalMatches,
			wltStr,
		)
	}
	fmt.Println("=========================================================================================")
	fmt.Println(" Mathematical Formulation: Bradley-Terry MLE with L2 regularization (LMSYS Arena standard)")
}

func executeMatrixCLI() {
	fmt.Print(tui.Banner())
	fmt.Println("\nComputing Pairwise Head-to-Head Win-Rate Matrix...")
	_, matrix, err := runner.ComputeArenaEloRatings(2500)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\nPairwise Win Frequency Matrix (Row vs Column):")
	fmt.Println("=========================================================================================")
	fmt.Printf("%-24s", "Model")
	for _, a := range matrix.Agents {
		fmt.Printf(" %-8s", truncate(a, 8))
	}
	fmt.Println("\n-----------------------------------------------------------------------------------------")

	for _, a1 := range matrix.Agents {
		fmt.Printf("%-24s", truncate(a1, 24))
		for _, a2 := range matrix.Agents {
			if a1 == a2 {
				fmt.Printf(" %-8s", "  -  ")
				continue
			}
			w := matrix.Wins[a1][a2]
			total := matrix.Matchups[a1][a2]
			if total > 0 {
				winPct := (w / float64(total)) * 100.0
				fmt.Printf(" %-7.1f%%", winPct)
			} else {
				fmt.Printf(" %-8s", "  N/A")
			}
		}
		fmt.Println()
	}
	fmt.Println("=========================================================================================")
}

func executeLeaderboardCLI() {
	fmt.Print(tui.Banner())
	var board []gauntlet.AgentScorecard
	if err := json.Unmarshal(data.LeaderboardData, &board); err != nil {
		fmt.Printf("Error parsing leaderboard data: %v\n", err)
		return
	}

	fmt.Println("\n=========================================================================================")
	fmt.Printf(" %-5s %-32s %-8s %-7s %-12s %-12s\n", "RANK", "AGENT / RUNTIME", "SCORE", "GRADE", "HANG RESIST", "SECRET SAFE")
	fmt.Println("-----------------------------------------------------------------------------------------")

	for i, a := range board {
		verified := ""
		if a.Verified {
			verified = " [VERIFIED]"
		}
		fmt.Printf(" #%-4d %-32s %-8.1f %-7s %-12s %-12s\n",
			i+1,
			truncate(a.AgentName+verified, 32),
			a.OverallScore,
			a.Grade,
			fmt.Sprintf("%.0f%%", a.HangResistance),
			fmt.Sprintf("%.0f%%", a.SecretSafety),
		)
	}
	fmt.Println("=========================================================================================")
}

func executeServe(port int) {
	fmt.Print(tui.Banner())
	fmt.Printf("\nBooting CrashBench API Engine on http://localhost:%d\n", port)
	fmt.Println("   Press Ctrl+C to stop.")

	mux := http.NewServeMux()

	// API: Standard Leaderboard
	mux.HandleFunc("/api/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data.LeaderboardData)
	})

	// API: LMSYS Arena Bradley-Terry Elo Ratings
	mux.HandleFunc("/api/arena/elo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ratings, matrix, err := runner.ComputeArenaEloRatings(2500)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ratings": ratings,
			"matrix":  matrix,
		})
	})

	// Static Web Arena (Better-Auth + LMSYS UI)
	var webDir string
	candidates := []string{
		"../crashbench-web/dist",
		"crashbench-web/dist",
		"./crashbench-web/dist",
		"C:\\Users\\tesse\\Desktop\\newthings\\crashbench-web\\dist",
		"../crashbench-web",
		"crashbench-web",
		"./crashbench-web",
		"C:\\Users\\tesse\\Desktop\\newthings\\crashbench-web",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			webDir = c
			break
		}
	}

	if webDir != "" {
		fileServer := http.FileServer(http.Dir(webDir))
		mux.Handle("/", fileServer)
		fmt.Printf("   Serving static web arena from: %s\n", webDir)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server failed: %v\n", err)
		os.Exit(1)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
