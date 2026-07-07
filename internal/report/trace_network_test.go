package report

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestParseTraceNetworkFailures(t *testing.T) {
	zipData := minimalNetworkZIP(t, `{"type":"resource-snapshot","snapshot":{"_monotonicTime":1000,"request":{"method":"GET","url":"https://api.test/x"},"response":{"status":500}}}`)
	events := parseTraceNetworkFailures(zipData, 10)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].OffsetMS != 0 {
		t.Fatalf("offset want 0 got %d", events[0].OffsetMS)
	}
	if events[0].Snippet == "" {
		t.Fatal("empty snippet")
	}
}

func TestEnrichStepsNetworkFromTrace(t *testing.T) {
	steps := []htmlStep{
		{Index: 0, Status: "passed", TraceOffsetMS: 0},
		{Index: 1, Status: "failed", TraceOffsetMS: 200},
	}
	enrichStepsNetworkFromTrace(steps, []htmlNetworkEvent{
		{OffsetMS: 180, Snippet: "GET https://api.test/x — HTTP 500"},
	})
	if steps[0].Network == "" {
		t.Fatal("expected network on passed step near trace event")
	}
	if steps[1].Network == "" {
		t.Fatal("expected network on failed step")
	}
}

func minimalNetworkZIP(t *testing.T, line string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("trace.network")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(line + "\n")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
