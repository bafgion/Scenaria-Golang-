package report

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/bafgion/scenaria-golang/internal/player"
)

type junitTestSuite struct {
	XMLName  xml.Name        `xml:"testsuite"`
	Name     string          `xml:"name,attr"`
	Tests    int             `xml:"tests,attr"`
	Failures int             `xml:"failures,attr"`
	Skipped  int             `xml:"skipped,attr"`
	TestCase []junitTestCase `xml:"testcase"`
}

type junitTestCase struct {
	ClassName string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Failure   *junitFailure `xml:"failure,omitempty"`
	Skipped   *junitSkipped `xml:"skipped,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
}

type junitSkipped struct {
	Message string `xml:"message,attr"`
}

func WriteJUnit(path string, result player.ExecutionResult) error {
	suite := junitTestSuite{
		Name:     "scenaria",
		Tests:    len(result.ScenarioResults),
		Failures: 0,
		Skipped:  0,
		TestCase: make([]junitTestCase, 0, len(result.ScenarioResults)),
	}

	for _, scenario := range result.ScenarioResults {
		tc := junitTestCase{
			ClassName: scenario.FeaturePath,
			Name:      scenario.Scenario,
		}
		switch scenario.Status {
		case "failed":
			suite.Failures++
			tc.Failure = &junitFailure{Message: scenario.Message}
		case "skipped":
			suite.Skipped++
			tc.Skipped = &junitSkipped{Message: scenario.Message}
		}
		suite.TestCase = append(suite.TestCase, tc)
	}

	payload, err := xml.MarshalIndent(suite, "", "  ")
	if err != nil {
		return fmt.Errorf("encode junit report %q: %w", path, err)
	}
	content := append([]byte(xml.Header), payload...)
	content = append(content, '\n')
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write junit report %q: %w", path, err)
	}
	return nil
}
