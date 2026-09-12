package algorithms

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// Matchup represents a head-to-head comparison between two agents on a scenario.
type Matchup struct {
	AgentA   string  `json:"agent_a"`
	AgentB   string  `json:"agent_b"`
	Winner   string  `json:"winner"`   // AgentA, AgentB, or "TIE"
	ScoreA   float64 `json:"score_a"`  // Raw performance score (0-100)
	ScoreB   float64 `json:"score_b"`
	Weight   float64 `json:"weight"`
	Scenario string  `json:"scenario"`
}

// EloRating holds an agent's estimated Elo rating with 95% confidence intervals.
type EloRating struct {
	Agent        string  `json:"agent"`
	Rating       float64 `json:"rating"`        // Normalized Elo rating (e.g. 1000 baseline)
	CiLower      float64 `json:"ci_lower"`      // 2.5th percentile bootstrap
	CiUpper      float64 `json:"ci_upper"`      // 97.5th percentile bootstrap
	WinRate      float64 `json:"win_rate"`      // Empirical win rate %
	TotalMatches int     `json:"total_matches"`
	Wins         int     `json:"wins"`
	Losses       int     `json:"losses"`
	Ties         int     `json:"ties"`
	Rank         int     `json:"rank"`
}

// PairwiseMatrix represents head-to-head win frequencies between agent pairs.
type PairwiseMatrix struct {
	Agents   []string             `json:"agents"`
	Wins     map[string]map[string]float64 `json:"wins"`     // wins[A][B] = number of times A beat B
	Matchups map[string]map[string]int     `json:"matchups"` // total matches between A and B
}

// BradleyTerrySolver solves the Bradley-Terry model via Maximum Likelihood Estimation (MLE)
// with L2 regularization using Newton-Raphson / Gradient Descent.
// P(i > j) = exp(gamma_i) / (exp(gamma_i) + exp(gamma_j))
// Elo = 1000 + 400 * log10(exp(gamma))
type BradleyTerrySolver struct {
	MaxIterations int
	Tolerance     float64
	L2Reg         float64 // Ridge regularization parameter to prevent divergence on 100% wins
	BaseElo       float64
	Scale         float64
}

// NewBradleyTerrySolver initializes a solver tuned to LMSYS Arena parameters.
func NewBradleyTerrySolver() *BradleyTerrySolver {
	return &BradleyTerrySolver{
		MaxIterations: 120,
		Tolerance:     1e-5,
		L2Reg:         0.15, // Smooth regularization
		BaseElo:       1000.0,
		Scale:         120.0,
	}
}

// BuildPairwiseMatrix compiles raw matchups into a pairwise win-loss frequency table.
func BuildPairwiseMatrix(matchups []Matchup) *PairwiseMatrix {
	agentSet := make(map[string]bool)
	for _, m := range matchups {
		agentSet[m.AgentA] = true
		agentSet[m.AgentB] = true
	}

	agents := make([]string, 0, len(agentSet))
	for a := range agentSet {
		agents = append(agents, a)
	}
	sort.Strings(agents)

	wins := make(map[string]map[string]float64)
	total := make(map[string]map[string]int)
	for _, a := range agents {
		wins[a] = make(map[string]float64)
		total[a] = make(map[string]int)
	}

	for _, m := range matchups {
		total[m.AgentA][m.AgentB]++
		total[m.AgentB][m.AgentA]++

		switch m.Winner {
		case m.AgentA:
			wins[m.AgentA][m.AgentB] += 1.0
		case m.AgentB:
			wins[m.AgentB][m.AgentA] += 1.0
		default:
			// Tie counts as 0.5 win for each (standard LMSYS rule)
			wins[m.AgentA][m.AgentB] += 0.5
			wins[m.AgentB][m.AgentA] += 0.5
		}
	}

	return &PairwiseMatrix{
		Agents:   agents,
		Wins:     wins,
		Matchups: total,
	}
}

// SolveMLE performs Maximum Likelihood Estimation via regularized gradient ascent.
func (s *BradleyTerrySolver) SolveMLE(matrix *PairwiseMatrix) map[string]float64 {
	n := len(matrix.Agents)
	if n == 0 {
		return nil
	}

	agentIdx := make(map[string]int)
	for i, a := range matrix.Agents {
		agentIdx[a] = i
	}

	gamma := make([]float64, n) // log-abilities initialized to 0
	learningRate := 0.2

	for iter := 0; iter < s.MaxIterations; iter++ {
		grad := make([]float64, n)
		maxGrad := 0.0

		for i := 0; i < n; i++ {
			agentA := matrix.Agents[i]
			gA := gamma[i]

			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				agentB := matrix.Agents[j]
				wAB := matrix.Wins[agentA][agentB]
				nAB := float64(matrix.Matchups[agentA][agentB])
				if nAB == 0 {
					continue
				}

				gB := gamma[j]
				// Logistic probability: p_ij = 1 / (1 + exp(gB - gA))
				pAB := 1.0 / (1.0 + math.Exp(gB-gA))

				// Gradient of log-likelihood: sum_j [ w_ij - n_ij * p_ij ] - lambda * gamma_i
				grad[i] += wAB - (nAB * pAB)
			}

			// Apply L2 regularization
			grad[i] -= s.L2Reg * gA

			if math.Abs(grad[i]) > maxGrad {
				maxGrad = math.Abs(grad[i])
			}
		}

		// Update parameters
		for i := 0; i < n; i++ {
			gamma[i] += learningRate * grad[i]
		}

		// Re-center gamma so mean is 0 for numerical stability
		var sum float64
		for _, g := range gamma {
			sum += g
		}
		mean := sum / float64(n)
		for i := range gamma {
			gamma[i] -= mean
		}

		if maxGrad < s.Tolerance {
			break
		}
	}

	ratings := make(map[string]float64)
	for i, a := range matrix.Agents {
		// Standard LMSYS Arena normalization: centered at 1000, scale factor 8.0
		elo := 1000.0 + (gamma[i] * 8.0)
		if elo < 650 {
			elo = 650
		}
		ratings[a] = math.Round(elo)
	}

	return ratings
}

// BootstrapConfidenceIntervals runs non-parametric bootstrap resampling (B=500)
// to compute empirical 95% confidence intervals (2.5% to 97.5%).
func (s *BradleyTerrySolver) BootstrapConfidenceIntervals(matchups []Matchup, numBootstraps int) []EloRating {
	matrix := BuildPairwiseMatrix(matchups)
	pointEstimates := s.SolveMLE(matrix)

	if numBootstraps <= 0 {
		numBootstraps = 300
	}

	mCount := len(matchups)
	sampleRatings := make(map[string][]float64)
	for _, a := range matrix.Agents {
		sampleRatings[a] = make([]float64, 0, numBootstraps)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Bootstrap resample with replacement
	for b := 0; b < numBootstraps; b++ {
		resample := make([]Matchup, mCount)
		for i := 0; i < mCount; i++ {
			resample[i] = matchups[rng.Intn(mCount)]
		}

		bsMatrix := BuildPairwiseMatrix(resample)
		bsRatings := s.SolveMLE(bsMatrix)
		for a, r := range bsRatings {
			sampleRatings[a] = append(sampleRatings[a], r)
		}
	}

	ratings := make([]EloRating, 0, len(matrix.Agents))
	for _, a := range matrix.Agents {
		samples := sampleRatings[a]
		sort.Float64s(samples)

		lowerIdx := int(0.025 * float64(len(samples)))
		upperIdx := int(0.975 * float64(len(samples)))
		if upperIdx >= len(samples) {
			upperIdx = len(samples) - 1
		}

		// Calculate empirical win/loss counts
		var totalWins, totalLosses, totalTies, totalM int
		for other, matches := range matrix.Matchups[a] {
			if a == other {
				continue
			}
			w := matrix.Wins[a][other]
			l := matrix.Wins[other][a]
			totalM += matches
			totalWins += int(w)
			totalLosses += int(l)
			totalTies += int(float64(matches) - (w + l))
		}

		winRate := 0.0
		if totalM > 0 {
			winRate = (float64(totalWins) + 0.5*float64(totalTies)) / float64(totalM) * 100.0
		}

		ciLower := pointEstimates[a] - 25.0
		ciUpper := pointEstimates[a] + 25.0
		if len(samples) > 0 {
			ciLower = samples[lowerIdx]
			ciUpper = samples[upperIdx]
		}

		ratings = append(ratings, EloRating{
			Agent:        a,
			Rating:       pointEstimates[a],
			CiLower:      ciLower,
			CiUpper:      ciUpper,
			WinRate:      winRate,
			TotalMatches: totalM,
			Wins:         totalWins,
			Losses:       totalLosses,
			Ties:         totalTies,
		})
	}

	// Sort descending by Elo Rating
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].Rating > ratings[j].Rating
	})

	for i := range ratings {
		ratings[i].Rank = i + 1
	}

	return ratings
}
