//go:build integration

package player

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func newIsolationServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		action, id := parts[0], parts[1]
		switch action {
		case "set":
			http.SetCookie(w, &http.Cookie{Name: "scenaria_iso", Value: id, Path: "/"})
			_, _ = w.Write([]byte(`<html><body>set</body></html>`))
		case "verify":
			cookie, err := r.Cookie("scenaria_iso")
			if err != nil || cookie.Value != id {
				_, _ = w.Write([]byte(`<html><body>LEAK</body></html>`))
				return
			}
			_, _ = w.Write([]byte(`<html><body>MATCH</body></html>`))
		default:
			http.NotFound(w, r)
		}
	}))
}

func isolationCase(name, base, id string) RunCase {
	return RunCase{
		FeaturePath: "isolation.feature",
		Name:        name,
		CaseID:      BuildCaseID("isolation.feature", name, 0),
		Steps: []gherkin.Step{
			{Keyword: "Допустим", Text: fmt.Sprintf(`открыт "%s/set/%s"`, base, id), Line: 1},
			{Keyword: "Допустим", Text: fmt.Sprintf(`открыт "%s/verify/%s"`, base, id), Line: 2},
			{Keyword: "Тогда", Text: `проверяю текст "MATCH" в "body"`, Line: 3},
		},
	}
}

func TestSequentialScenariosDoNotShareBrowserState(t *testing.T) {
	srv := newIsolationServer(t)
	defer srv.Close()

	plan := ExecutionPlan{
		Cases: []RunCase{
			isolationCase("Scenario A", srv.URL, "A"),
			isolationCase("Scenario B", srv.URL, "B"),
		},
	}
	runner := BrowserRunner{
		Executor: NewPlaywrightExecutor(PlaywrightExecutorOptions{
			BrowserName:   "chromium",
			Headless:      true,
			CloseAfterRun: true,
		}),
		ParallelWorkers: 1,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := runner.Execute(ctx, plan)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, sr := range result.ScenarioResults {
		if sr.Status != "passed" {
			t.Fatalf("scenario %q failed: %s", sr.Scenario, sr.Message)
		}
	}
}

func TestParallelPoolDoesNotShareBrowserStorage(t *testing.T) {
	srv := newIsolationServer(t)
	defer srv.Close()

	cases := make([]RunCase, 0, 6)
	for i := 1; i <= 6; i++ {
		name := fmt.Sprintf("Scenario %d", i)
		cases = append(cases, isolationCase(name, srv.URL, fmt.Sprintf("%d", i)))
	}
	plan := ExecutionPlan{Cases: cases}
	runner := BrowserRunner{
		Executor: NewPlaywrightExecutor(PlaywrightExecutorOptions{
			BrowserName:   "chromium",
			Headless:      true,
			CloseAfterRun: true,
		}),
		ParallelWorkers: 2,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	result, err := runner.Execute(ctx, plan)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, sr := range result.ScenarioResults {
		if sr.Status != "passed" {
			t.Fatalf("scenario %q failed: %s", sr.Scenario, sr.Message)
		}
	}
}

func TestRunContextsUseDistinctDownloadDirs(t *testing.T) {
	root := t.TempDir()
	a := NewRunContext(nil, 11, root)
	b := NewRunContext(nil, 22, root)
	if a.DownloadDir() == b.DownloadDir() {
		t.Fatalf("expected distinct download dirs, got %q", a.DownloadDir())
	}
}
