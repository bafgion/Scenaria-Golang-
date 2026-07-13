package winfocus

import "testing"

func TestScoreChromiumTitle(t *testing.T) {
	score := scoreChromiumTitle("2moodstore - online shop - Chromium", "2moodstore", "https://www.2moodstore.com/")
	if score < 100 {
		t.Fatalf("expected title match, got %d", score)
	}
	score = scoreChromiumTitle("example.com - Home - Chromium", "", "https://example.com/page")
	if score < 90 {
		t.Fatalf("expected host match, got %d", score)
	}
	score = scoreChromiumTitle("Gmail - Google Chrome", "", "https://example.com/page")
	if score != 0 {
		t.Fatalf("expected unrelated Chrome window to be rejected, got %d", score)
	}
	score = scoreChromiumTitle("", "", "")
	if score != 0 {
		t.Fatalf("expected empty window title to be rejected, got %d", score)
	}
	score = scoreChromiumTitle("about:blank - Chromium", "", "")
	if score < 10 {
		t.Fatalf("expected Chromium browser window fallback, got %d", score)
	}
}

func TestHostFromURL(t *testing.T) {
	if got := hostFromURL("https://www.2moodstore.com/"); got != "www.2moodstore.com" {
		t.Fatalf("unexpected host: %q", got)
	}
}

func TestCDPWindowID(t *testing.T) {
	for _, raw := range []any{1, int64(2), float64(3)} {
		if got, ok := cdpWindowID(raw); !ok || got <= 0 {
			t.Fatalf("expected valid window id for %#v, got %d ok=%v", raw, got, ok)
		}
	}
	for _, raw := range []any{0, int64(0), float64(0), "1", nil} {
		if got, ok := cdpWindowID(raw); ok || got != 0 {
			t.Fatalf("expected invalid window id for %#v, got %d ok=%v", raw, got, ok)
		}
	}
}

func TestBrowserProcessIDFromPayload(t *testing.T) {
	payload := map[string]any{
		"processInfo": []any{
			map[string]any{"type": "renderer", "id": float64(10)},
			map[string]any{"type": "browser", "id": float64(42)},
		},
	}
	if got := browserProcessIDFromPayload(payload); got != 42 {
		t.Fatalf("expected browser pid 42, got %d", got)
	}
	if got := browserProcessIDFromPayload(map[string]any{"processInfo": []any{map[string]any{"type": "renderer", "id": float64(10)}}}); got != 0 {
		t.Fatalf("expected missing browser pid to be 0, got %d", got)
	}
}
