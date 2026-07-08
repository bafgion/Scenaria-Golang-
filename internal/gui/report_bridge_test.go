package gui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"
)

func bridgePOST(t *testing.T, svc *Service, path string, body []byte) *http.Response {
	t.Helper()
	token := svc.ReportBridgeToken()
	if token == "" {
		t.Fatal("empty bridge token")
	}
	req, err := http.NewRequest(http.MethodPost, svc.ReportBridgeURL()+path, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Scenaria-Bridge-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func TestReportBridgeGoto(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	feature := filepath.Join(root, "a.feature")
	got := make(chan ReportGotoRequest, 1)
	if err := svc.EnsureReportBridge(func(req ReportGotoRequest) { got <- req }, nil, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	body, err := json.Marshal(ReportGotoRequest{
		FeaturePath: feature,
		Scenario:    "S1",
		LeafIndex:   2,
		Line:        12,
	})
	if err != nil {
		t.Fatal(err)
	}
	resp := bridgePOST(t, svc, "/goto", body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	select {
	case req := <-got:
		if req.Line != 12 || req.LeafIndex != 2 {
			t.Fatalf("unexpected req: %+v", req)
		}
	default:
		t.Fatal("handler not called")
	}
}

func TestReportBridgeRejectsMissingToken(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	if err := svc.EnsureReportBridge(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	resp, err := http.Post(svc.ReportBridgeURL()+"/goto", "application/json", bytes.NewReader([]byte(`{"feature_path":"x"}`)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestReportBridgeRotatesToken(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	if err := svc.EnsureReportBridge(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	oldToken := svc.ReportBridgeToken()
	if oldToken == "" {
		t.Fatal("empty old token")
	}
	_, newToken := svc.ReportBridgeCredentials()
	if newToken == "" || newToken == oldToken {
		t.Fatalf("token was not rotated: old=%q new=%q", oldToken, newToken)
	}
	body := []byte(`{"feature_path":"x.feature"}`)
	req, err := http.NewRequest(http.MethodPost, svc.ReportBridgeURL()+"/goto", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Scenaria-Bridge-Token", oldToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("old token status %d", resp.StatusCode)
	}
}

func TestReportBridgeDropsOldestTokenWhenBoundExceeded(t *testing.T) {
	svc := NewService()
	if err := svc.EnsureReportBridge(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()

	oldest := svc.ReportBridgeToken()
	for i := 0; i < maxBridgeTokens; i++ {
		_, _ = svc.ReportBridgeCredentials()
	}

	req, err := http.NewRequest(http.MethodPost, svc.ReportBridgeURL()+"/goto", bytes.NewReader([]byte(`{"feature_path":"x.feature"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Scenaria-Bridge-Token", oldest)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected oldest token to expire, got %d", resp.StatusCode)
	}
}

func TestReportBridgeRerun(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	feature := filepath.Join(root, "b.feature")
	got := make(chan ReportRerunRequest, 1)
	if err := svc.EnsureReportBridge(nil, func(req ReportRerunRequest) { got <- req }, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	body, err := json.Marshal(ReportRerunRequest{FeaturePath: feature, Scenario: "Login"})
	if err != nil {
		t.Fatal(err)
	}
	resp := bridgePOST(t, svc, "/rerun", body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	select {
	case req := <-got:
		if req.Scenario != "Login" {
			t.Fatalf("unexpected: %+v", req)
		}
	default:
		t.Fatal("handler not called")
	}
}

func TestReportBridgeTrace(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	traceRel := filepath.Join(".scenaria", "traces", "a.zip")
	got := make(chan ReportTraceRequest, 1)
	if err := svc.EnsureReportBridge(nil, nil, func(req ReportTraceRequest) { got <- req }); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	body, err := json.Marshal(ReportTraceRequest{TracePath: traceRel, ReportDir: ".scenaria"})
	if err != nil {
		t.Fatal(err)
	}
	resp := bridgePOST(t, svc, "/trace", body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	select {
	case req := <-got:
		if req.TracePath != traceRel {
			t.Fatalf("unexpected: %+v", req)
		}
	default:
		t.Fatal("handler not called")
	}
}

func TestReportBridgeCORSOptions(t *testing.T) {
	svc := NewService()
	if err := svc.EnsureReportBridge(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	req, err := http.NewRequest(http.MethodOptions, svc.ReportBridgeURL()+"/goto", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "null")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("options status %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Headers") == "" {
		t.Fatal("missing CORS allow-headers")
	}
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "null" {
		t.Fatalf("allow-origin %q", got)
	}
}

func TestReportBridgeCORSOptionsRejectsRemoteOrigin(t *testing.T) {
	svc := NewService()
	if err := svc.EnsureReportBridge(nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	defer svc.CloseReportBridge()
	req, err := http.NewRequest(http.MethodOptions, svc.ReportBridgeURL()+"/goto", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "https://evil.example")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("options status %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unexpected allow-origin %q", got)
	}
}

func TestAllowedBridgeOriginRejectsRemote(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:1/health", nil)
	req.Header.Set("Origin", "https://evil.example")
	if allowedBridgeOrigin(req) != "" {
		t.Fatal("expected remote origin rejected")
	}
	req.Header.Set("Origin", "http://127.0.0.1:8765")
	if allowedBridgeOrigin(req) != "http://127.0.0.1:8765" {
		t.Fatal("expected localhost origin allowed")
	}
}

func TestResolveTraceFromReport(t *testing.T) {
	dir := `C:\proj\.scenaria`
	rel := `traces\run.zip`
	got := filepath.Join(dir, rel)
	if !stringsContains(got, "run.zip") {
		t.Fatalf("unexpected %s", got)
	}
}

func stringsContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
