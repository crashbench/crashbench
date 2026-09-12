package algorithms

import (
	"bytes"
	"compress/flate"
	"math"
)

// EntropyProfile holds information-theoretic metrics for an execution stream.
type EntropyProfile struct {
	TotalBytes         int     `json:"total_bytes"`
	ShannonEntropy     float64 `json:"shannon_entropy"`      // Bits per byte (0.0 to 8.0)
	CompressionRatio   float64 `json:"compression_ratio"`    // Flate compressed / Raw size
	RepetitionIndex    float64 `json:"repetition_index"`     // 0.0 (random/healthy) to 1.0 (pure loop)
	IsContextBomb      bool    `json:"is_context_bomb"`      // Exceeds runaway token threshold
	IsRepetitiveFreeze bool    `json:"is_repetitive_freeze"` // Low entropy infinite recursion
}

// ComputeShannonEntropy calculates the Shannon information entropy H(X) = -sum(P(x) * log2(P(x)))
// of a byte slice. For natural English/code text, H is typically between 3.5 and 5.0 bits/byte.
// Extremely low H (< 1.8) indicates runaway infinite printing loops or repetitive memory crashes.
func ComputeShannonEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}

	var freq [256]int
	for _, b := range data {
		freq[b]++
	}

	total := float64(len(data))
	entropy := 0.0

	for _, count := range freq {
		if count > 0 {
			p := float64(count) / total
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

// ComputeCompressionRatio measures LZ77/Deflate compressibility.
// Highly repetitive compiler dumps compress by >95% (ratio < 0.05).
func ComputeCompressionRatio(data []byte) float64 {
	if len(data) == 0 {
		return 1.0
	}

	var buf bytes.Buffer
	zw, err := flate.NewWriter(&buf, flate.BestSpeed)
	if err != nil {
		return 1.0
	}
	_, _ = zw.Write(data)
	_ = zw.Close()

	return float64(buf.Len()) / float64(len(data))
}

// AnalyzeContextEntropy evaluates terminal stream dynamics for context bomb risks.
func AnalyzeContextEntropy(data []byte) EntropyProfile {
	shannon := ComputeShannonEntropy(data)
	compRatio := ComputeCompressionRatio(data)

	// Repetition index: inverse of entropy normalized between 0 and 8 bits
	repetition := 1.0 - (shannon / 8.0)
	if repetition < 0 {
		repetition = 0
	}

	isBomb := len(data) > 50000 || (len(data) > 10000 && compRatio < 0.08)
	isFreeze := shannon < 2.0 && len(data) > 2000

	return EntropyProfile{
		TotalBytes:         len(data),
		ShannonEntropy:     shannon,
		CompressionRatio:   compRatio,
		RepetitionIndex:    repetition,
		IsContextBomb:      isBomb,
		IsRepetitiveFreeze: isFreeze,
	}
}
