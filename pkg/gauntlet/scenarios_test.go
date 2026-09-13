package gauntlet

import (
	"strings"
	"testing"
)

func TestStandardScenarios(t *testing.T) {
	scenarios := GetStandardScenarios()
	if len(scenarios) != 8 {
		t.Fatalf("expected 8 standard scenarios, got %d", len(scenarios))
	}

	totalWeight := 0.0
	seenIDs := make(map[string]bool)

	for _, s := range scenarios {
		if s.ID == "" {
			t.Errorf("scenario ID is empty")
		}
		if seenIDs[s.ID] {
			t.Errorf("duplicate scenario ID: %s", s.ID)
		}
		seenIDs[s.ID] = true

		if s.Name == "" {
			t.Errorf("scenario %s has empty Name", s.ID)
		}
		if s.Category == "" {
			t.Errorf("scenario %s has empty Category", s.ID)
		}
		if !strings.HasPrefix(s.CWE, "CWE-") {
			t.Errorf("scenario %s has invalid CWE: %s", s.ID, s.CWE)
		}
		if !strings.HasPrefix(s.OWASP, "OWASP-") {
			t.Errorf("scenario %s has invalid OWASP: %s", s.ID, s.OWASP)
		}
		if s.Weight <= 0 {
			t.Errorf("scenario %s has non-positive weight: %f", s.ID, s.Weight)
		}

		payload := GenerateScenarioPayload(s.ID)
		if payload == "" {
			t.Errorf("scenario %s has empty payload", s.ID)
		}

		totalWeight += s.Weight
	}

	if totalWeight != 100.0 {
		t.Errorf("total scenario weights should equal 100.0, got %f", totalWeight)
	}
}
