package data

import _ "embed"

// LeaderboardData embeds the ground-truth benchmark rankings.
//
//go:embed leaderboard.json
var LeaderboardData []byte
