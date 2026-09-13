package export

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/crashbench/crashbench/pkg/gauntlet"
)

// SARIFReport represents the root of a SARIF 2.1.0 document.
type SARIFReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SARIFRule `json:"rules"`
}

type SARIFRule struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	ShortDescription SARIFText              `json:"shortDescription"`
	FullDescription  SARIFText              `json:"fullDescription"`
	HelpURI          string                 `json:"helpUri,omitempty"`
	Properties       map[string]interface{} `json:"properties,omitempty"`
}

type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	RuleIndex int             `json:"ruleIndex"`
	Level     string          `json:"level"` // "error", "warning", "note"
	Message   SARIFText       `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           SARIFRegion           `json:"region"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFRegion struct {
	StartLine int `json:"startLine"`
}

type SARIFText struct {
	Text string `json:"text"`
}

// GenerateSARIF converts an AgentScorecard into standard SARIF 2.1.0 bytes for GitHub Code Scanning.
func GenerateSARIF(scorecard gauntlet.AgentScorecard, toolVersion string) ([]byte, error) {
	scenarios := gauntlet.GetStandardScenarios()
	rules := make([]SARIFRule, 0, len(scenarios))
	ruleIndexMap := make(map[string]int)

	for idx, s := range scenarios {
		ruleIndexMap[s.ID] = idx
		tags := []string{"chaos-engineering", "agent-safety"}
		if s.CWE != "" {
			tags = append(tags, s.CWE)
		}
		if s.OWASP != "" {
			tags = append(tags, s.OWASP)
		}

		rules = append(rules, SARIFRule{
			ID:               s.ID,
			Name:             s.Name,
			ShortDescription: SARIFText{Text: s.Description},
			FullDescription:  SARIFText{Text: fmt.Sprintf("%s. Impact: %s", s.Description, s.Impact)},
			HelpURI:          "https://github.com/crashbench/crashbench",
			Properties: map[string]interface{}{
				"category": string(s.Category),
				"cwe":      s.CWE,
				"owasp":    s.OWASP,
				"tags":     tags,
			},
		})
	}

	results := make([]SARIFResult, 0, len(scorecard.Results))
	for _, r := range scorecard.Results {
		if !r.Passed {
			msg := r.FailureReason
			if msg == "" {
				msg = fmt.Sprintf("Agent failed operational chaos scenario %s", r.ScenarioID)
			}
			idx, ok := ruleIndexMap[r.ScenarioID]
			if !ok {
				idx = 0
			}

			results = append(results, SARIFResult{
				RuleID:    r.ScenarioID,
				RuleIndex: idx,
				Level:     "error",
				Message:   SARIFText{Text: msg},
				Locations: []SARIFLocation{
					{
						PhysicalLocation: SARIFPhysicalLocation{
							ArtifactLocation: SARIFArtifactLocation{URI: "terminal/stdout"},
							Region:           SARIFRegion{StartLine: 1},
						},
					},
				},
			})
		}
	}

	report := SARIFReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "CrashBench",
						Version:        toolVersion,
						InformationURI: "https://github.com/crashbench/crashbench",
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}

	return json.MarshalIndent(report, "", "  ")
}

// WriteSARIFFile exports SARIF 2.1.0 report directly to a file.
func WriteSARIFFile(filePath string, scorecard gauntlet.AgentScorecard, toolVersion string) error {
	data, err := GenerateSARIF(scorecard, toolVersion)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
