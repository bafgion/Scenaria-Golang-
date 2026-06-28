package vanessa

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type RunSnapshot struct {
	Cases            []CaseResult
	CurrentScenario  string
	CompletedCases   int
	TotalPlanned     int
}

type RunMonitor struct {
	junitDir      string
	scenarioLog   string
	totalPlanned  int
	seen          map[string]struct{}
	mu            sync.Mutex
}

func NewRunMonitor(runDir string, totalPlanned int) *RunMonitor {
	if totalPlanned < 1 {
		totalPlanned = 1
	}
	return &RunMonitor{
		junitDir:     filepath.Join(runDir, "junit"),
		scenarioLog:  filepath.Join(runDir, "scenario.log"),
		totalPlanned: totalPlanned,
		seen:         map[string]struct{}{},
	}
}

func (m *RunMonitor) Poll() RunSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	cases := make([]CaseResult, 0)
	for _, item := range ParseJUnitDir(m.junitDir) {
		key := item.Name + "\x00" + item.Path
		if _, ok := m.seen[key]; ok {
			continue
		}
		m.seen[key] = struct{}{}
		cases = append(cases, item)
	}
	current := readCurrentScenarioLabel(m.scenarioLog)
	if current == "" && len(cases) > 0 {
		current = cases[len(cases)-1].Name
	}
	completed := len(m.seen)
	total := m.totalPlanned
	if completed > total {
		total = completed
	}
	return RunSnapshot{
		Cases:           cases,
		CurrentScenario: current,
		CompletedCases:  completed,
		TotalPlanned:    total,
	}
}

func readCurrentScenarioLabel(path string) string {
	payload, err := os.ReadFile(path)
	if err != nil || len(payload) == 0 {
		return ""
	}
	const tail = 24 * 1024
	if len(payload) > tail {
		payload = payload[len(payload)-tail:]
	}
	lines := strings.Split(string(payload), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if idx := strings.Index(line, "Сценарий"); idx >= 0 {
			return strings.TrimSpace(line[idx:])
		}
		return line
	}
	return ""
}
