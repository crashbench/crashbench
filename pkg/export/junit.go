package export

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/crashbench/crashbench/pkg/gauntlet"
)

// JUnitTestSuites represents the root JUnit XML document.
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	Name       string           `xml:"name,attr"`
	Tests      int              `xml:"tests,attr"`
	Failures   int              `xml:"failures,attr"`
	Time       string           `xml:"time,attr"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Time      string          `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	Classname string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
}

type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Body    string `xml:",cdata"`
}

// GenerateJUnit compiles an AgentScorecard into JUnit XML bytes.
func GenerateJUnit(scorecard gauntlet.AgentScorecard) ([]byte, error) {
	totalTests := len(scorecard.Results)
	failuresCount := scorecard.ScenariosTotal - scorecard.ScenariosPassed
	totalTimeSec := fmt.Sprintf("%.3f", float64(scorecard.TotalDurationMs)/1000.0)

	testCases := make([]JUnitTestCase, 0, totalTests)
	for _, r := range scorecard.Results {
		durationSec := fmt.Sprintf("%.3f", float64(r.DurationMs)/1000.0)
		tc := JUnitTestCase{
			Classname: fmt.Sprintf("CrashBench.%s", r.Category),
			Name:      fmt.Sprintf("[%s] %s (%s)", r.ScenarioID, r.ScenarioName, r.CWE),
			Time:      durationSec,
		}

		if !r.Passed {
			tc.Failure = &JUnitFailure{
				Message: r.FailureReason,
				Type:    r.CWE,
				Body:    r.TerminalReplay,
			}
		}

		testCases = append(testCases, tc)
	}

	suite := JUnitTestSuite{
		Name:      fmt.Sprintf("CrashBench Gauntlet: %s", scorecard.AgentName),
		Tests:     totalTests,
		Failures:  failuresCount,
		Time:      totalTimeSec,
		TestCases: testCases,
	}

	root := JUnitTestSuites{
		Name:       "CrashBench Operational Chaos Benchmark",
		Tests:      totalTests,
		Failures:   failuresCount,
		Time:       totalTimeSec,
		TestSuites: []JUnitTestSuite{suite},
	}

	data, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), data...), nil
}

// WriteJUnitFile exports JUnit XML directly to a file.
func WriteJUnitFile(filePath string, scorecard gauntlet.AgentScorecard) error {
	data, err := GenerateJUnit(scorecard)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
