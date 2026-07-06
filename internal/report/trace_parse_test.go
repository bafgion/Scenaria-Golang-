package report

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestParseTraceActions(t *testing.T) {
	zipData := minimalTraceZIP(t, []string{
		`{"type":"action","startTime":1000,"class":"Frame","method":"goto","params":{"url":"https://example.com"}}`,
		`{"type":"action","startTime":1250,"class":"Frame","method":"click","params":{"selector":"#btn"}}`,
		`{"type":"event","name":"foo"}`,
	})
	events := parseTraceActions(zipData, 10)
	if len(events) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(events))
	}
	if events[0].OffsetMS != 0 || !strings.Contains(events[0].Title, "goto") {
		t.Fatalf("first: %+v", events[0])
	}
	if events[1].OffsetMS != 250 {
		t.Fatalf("second offset want 250 got %d", events[1].OffsetMS)
	}
}

func minimalTraceZIP(t *testing.T, lines []string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("trace.trace")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range lines {
		if _, err := w.Write([]byte(line + "\n")); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
