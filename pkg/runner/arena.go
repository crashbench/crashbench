package runner

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/crashbench/crashbench/data"
	"github.com/crashbench/crashbench/pkg/algorithms"
	"github.com/crashbench/crashbench/pkg/gauntlet"
)

// ComputeArenaEloRatings calculates Bradley-Terry MLE Elo ratings and 95% bootstrap CIs
// based on thousands of pairwise scenario matchups between benchmarked agents.
func ComputeArenaEloRatings(numSimulatedBattles int) ([]algorithms.EloRating, *algorithms.PairwiseMatrix, error) {
	var scorecards []gauntlet.AgentScorecard
	if err := json.Unmarshal(data.LeaderboardData, &scorecards); err != nil {
		return nil, nil, err
	}

	if numSimulatedBattles <= 0 {
		numSimulatedBattles = 2500
	}

	// Index scorecards by agent name
	byName := make(map[string]gauntlet.AgentScorecard)
	agentNames := make([]string, 0, len(scorecards))
	for _, sc := range scorecards {
		byName[sc.AgentName] = sc
		agentNames = append(agentNames, sc.AgentName)
	}

	// Generate synthetic pairwise battle results based on scenario performance distributions
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	scenarios := gauntlet.GetStandardScenarios()
	matchups := make([]algorithms.Matchup, 0, numSimulatedBattles)

	for i := 0; i < numSimulatedBattles; i++ {
		// Pick two distinct agents at random
		idxA := rng.Intn(len(agentNames))
		idxB := rng.Intn(len(agentNames))
		for idxB == idxA {
			idxB = rng.Intn(len(agentNames))
		}

		agentA := byName[agentNames[idxA]]
		agentB := byName[agentNames[idxB]]
		sc := scenarios[rng.Intn(len(scenarios))]

		// Find score for this scenario in both agents
		scoreA := getAgentScenarioScore(agentA, sc.ID)
		scoreB := getAgentScenarioScore(agentB, sc.ID)

		// Add subtle stochastic noise (+- 5 points) to simulate real-world battle variance
		noiseA := (rng.Float64() - 0.5) * 10.0
		noiseB := (rng.Float64() - 0.5) * 10.0
		adjA := scoreA + noiseA
		adjB := scoreB + noiseB

		winner := "TIE"
		if adjA > adjB+3.0 {
			winner = agentA.AgentName
		} else if adjB > adjA+3.0 {
			winner = agentB.AgentName
		}

		matchups = append(matchups, algorithms.Matchup{
			AgentA:   agentA.AgentName,
			AgentB:   agentB.AgentName,
			Winner:   winner,
			ScoreA:   scoreA,
			ScoreB:   scoreB,
			Weight:   sc.Weight,
			Scenario: sc.ID,
		})
	}

	solver := algorithms.NewBradleyTerrySolver()
	pairwiseMatrix := algorithms.BuildPairwiseMatrix(matchups)
	ratings := solver.BootstrapConfidenceIntervals(matchups, 300)

	return ratings, pairwiseMatrix, nil
}

func getAgentScenarioScore(sc gauntlet.AgentScorecard, scnID string) float64 {
	for _, r := range sc.Results {
		if r.ScenarioID == scnID {
			return r.Score
		}
	}
	return sc.OverallScore
}
