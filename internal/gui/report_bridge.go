package gui

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const bridgeTokenHeader = "X-Scenaria-Bridge-Token"
const maxBridgeTokens = 8

// ReportGotoRequest is sent from an HTML report to focus a step in the IDE.
type ReportGotoRequest struct {
	FeaturePath string `json:"feature_path"`
	Scenario    string `json:"scenario"`
	LeafIndex   int    `json:"leaf_index"`
	Line        int    `json:"line"`
	ReportID    string `json:"report_id,omitempty"`
}

// ReportRerunRequest asks the IDE to re-run one scenario.
type ReportRerunRequest struct {
	FeaturePath string `json:"feature_path"`
	Scenario    string `json:"scenario"`
	ReportID    string `json:"report_id,omitempty"`
}

// ReportTraceRequest opens a Playwright trace archive in the IDE.
type ReportTraceRequest struct {
	TracePath     string `json:"trace_path"`
	ReportDir     string `json:"report_dir"`
	TraceOffsetMS int64  `json:"trace_offset_ms,omitempty"`
	StepIndex     int    `json:"step_index,omitempty"`
	ReportID      string `json:"report_id,omitempty"`
}

type reportBridge struct {
	mu      sync.Mutex
	server  *http.Server
	baseURL string
	token   string
	tokens  map[string]struct{}
	order   []string
	onGoto  func(ReportGotoRequest)
	onRerun func(ReportRerunRequest)
	onTrace func(ReportTraceRequest)
}

func newBridgeToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func (b *reportBridge) rotateToken() (string, error) {
	if b == nil {
		return "", nil
	}
	token, err := newBridgeToken()
	if err != nil {
		return "", err
	}
	b.mu.Lock()
	b.addTokenLocked(token)
	b.mu.Unlock()
	return token, nil
}

func startReportBridge(onGoto func(ReportGotoRequest), onRerun func(ReportRerunRequest), onTrace func(ReportTraceRequest)) (*reportBridge, error) {
	token, err := newBridgeToken()
	if err != nil {
		return nil, fmt.Errorf("bridge token: %w", err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("listen report bridge: %w", err)
	}
	b := &reportBridge{
		token:   token,
		tokens:  map[string]struct{}{token: {}},
		order:   []string{token},
		onGoto:  onGoto,
		onRerun: onRerun,
		onTrace: onTrace,
		baseURL: "http://" + ln.Addr().String(),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/goto", b.handleGoto)
	mux.HandleFunc("/rerun", b.handleRerun)
	mux.HandleFunc("/trace", b.handleTrace)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		setReportBridgeCORS(w, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	b.server = &http.Server{Handler: mux}
	go func() {
		_ = b.server.Serve(ln)
	}()
	return b, nil
}

func (b *reportBridge) addTokenLocked(token string) {
	if b == nil || token == "" {
		return
	}
	if b.tokens == nil {
		b.tokens = make(map[string]struct{})
	}
	if _, exists := b.tokens[token]; !exists {
		b.order = append(b.order, token)
	}
	b.tokens[token] = struct{}{}
	b.token = token
	for len(b.order) > maxBridgeTokens {
		oldest := b.order[0]
		b.order = b.order[1:]
		delete(b.tokens, oldest)
	}
}

func (b *reportBridge) close() error {
	if b == nil || b.server == nil {
		return nil
	}
	return b.server.Close()
}

// EnsureReportBridge starts the localhost bridge used by HTML reports (idempotent).
func (s *Service) EnsureReportBridge(onGoto func(ReportGotoRequest), onRerun func(ReportRerunRequest), onTrace func(ReportTraceRequest)) error {
	s.reportBridgeMu.Lock()
	defer s.reportBridgeMu.Unlock()
	if s.reportBridge != nil {
		return nil
	}
	wrapGoto := func(req ReportGotoRequest) {
		if _, err := s.confineFeaturePath(req.FeaturePath); err != nil {
			return
		}
		if onGoto != nil {
			onGoto(req)
		}
	}
	wrapRerun := func(req ReportRerunRequest) {
		if _, err := s.confineFeaturePath(req.FeaturePath); err != nil {
			return
		}
		if onRerun != nil {
			onRerun(req)
		}
	}
	wrapTrace := func(req ReportTraceRequest) {
		if _, err := s.confineArtifactPath(req.TracePath); err != nil {
			return
		}
		if req.ReportDir != "" {
			if _, err := s.confineArtifactPath(req.ReportDir); err != nil {
				return
			}
		}
		if onTrace != nil {
			onTrace(req)
		}
	}
	b, err := startReportBridge(wrapGoto, wrapRerun, wrapTrace)
	if err != nil {
		return err
	}
	s.reportBridge = b
	return nil
}

// ReportBridgeURL returns the base URL for IDE bridge requests from HTML reports.
func (s *Service) ReportBridgeURL() string {
	s.reportBridgeMu.Lock()
	defer s.reportBridgeMu.Unlock()
	if s.reportBridge == nil {
		return ""
	}
	return s.reportBridge.URL()
}

// ReportBridgeToken returns the per-session secret required for bridge POST requests.
func (s *Service) ReportBridgeToken() string {
	s.reportBridgeMu.Lock()
	defer s.reportBridgeMu.Unlock()
	if s.reportBridge == nil {
		return ""
	}
	return s.reportBridge.tokenValue()
}

// ReportBridgeCredentials rotates the bridge token and returns credentials for a newly written report.
func (s *Service) ReportBridgeCredentials() (string, string) {
	s.reportBridgeMu.Lock()
	defer s.reportBridgeMu.Unlock()
	if s.reportBridge == nil {
		return "", ""
	}
	token, err := s.reportBridge.rotateToken()
	if err != nil {
		return s.reportBridge.URL(), ""
	}
	return s.reportBridge.URL(), token
}

// CloseReportBridge stops the report bridge HTTP server.
func (s *Service) CloseReportBridge() {
	s.reportBridgeMu.Lock()
	defer s.reportBridgeMu.Unlock()
	if s.reportBridge != nil {
		_ = s.reportBridge.close()
		s.reportBridge = nil
	}
}

func (b *reportBridge) URL() string {
	if b == nil {
		return ""
	}
	return b.baseURL
}

func (b *reportBridge) tokenValue() string {
	if b == nil {
		return ""
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.token
}

func setReportBridgeCORS(w http.ResponseWriter, r *http.Request) {
	if origin := allowedBridgeOrigin(r); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+bridgeTokenHeader)
}

func allowedBridgeOrigin(r *http.Request) string {
	if r == nil {
		return ""
	}
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" || origin == "null" {
		return "null"
	}
	u, err := url.Parse(origin)
	if err != nil {
		return ""
	}
	host := strings.ToLower(u.Hostname())
	if host != "127.0.0.1" && host != "localhost" {
		return ""
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	return origin
}

func (b *reportBridge) authorize(r *http.Request) bool {
	if b == nil || r == nil {
		return false
	}
	got := r.Header.Get(bridgeTokenHeader)
	b.mu.Lock()
	_, ok := b.tokens[got]
	b.mu.Unlock()
	return ok
}

func rejectDisallowedBridgeOrigin(w http.ResponseWriter, r *http.Request) bool {
	if r == nil {
		return false
	}
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return false
	}
	if allowedBridgeOrigin(r) != "" {
		return false
	}
	http.Error(w, "origin not allowed", http.StatusForbidden)
	return true
}

func (b *reportBridge) handleGoto(w http.ResponseWriter, r *http.Request) {
	setReportBridgeCORS(w, r)
	if r.Method == http.MethodOptions {
		if rejectDisallowedBridgeOrigin(w, r) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !b.authorize(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req ReportGotoRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	req.FeaturePath = strings.TrimSpace(req.FeaturePath)
	if req.FeaturePath == "" {
		http.Error(w, "feature_path required", http.StatusBadRequest)
		return
	}
	if b.onGoto != nil {
		b.onGoto(req)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true}`))
}

func (b *reportBridge) handleRerun(w http.ResponseWriter, r *http.Request) {
	setReportBridgeCORS(w, r)
	if r.Method == http.MethodOptions {
		if rejectDisallowedBridgeOrigin(w, r) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !b.authorize(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req ReportRerunRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	req.FeaturePath = strings.TrimSpace(req.FeaturePath)
	req.Scenario = strings.TrimSpace(req.Scenario)
	if req.FeaturePath == "" || req.Scenario == "" {
		http.Error(w, "feature_path and scenario required", http.StatusBadRequest)
		return
	}
	if b.onRerun != nil {
		b.onRerun(req)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true}`))
}

func (b *reportBridge) handleTrace(w http.ResponseWriter, r *http.Request) {
	setReportBridgeCORS(w, r)
	if r.Method == http.MethodOptions {
		if rejectDisallowedBridgeOrigin(w, r) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !b.authorize(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req ReportTraceRequest
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	req.TracePath = strings.TrimSpace(req.TracePath)
	req.ReportDir = strings.TrimSpace(req.ReportDir)
	if req.TracePath == "" {
		http.Error(w, "trace_path required", http.StatusBadRequest)
		return
	}
	if b.onTrace != nil {
		b.onTrace(req)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true}`))
}
