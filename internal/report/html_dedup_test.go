package report

import "testing"

func TestDedupeScreenshotURLs(t *testing.T) {
	dup := "data:image/png;base64,AAAA"
	payload := htmlReportPayload{
		Scenarios: []htmlScenario{{
			Screenshot: dup,
			Steps: []htmlStep{
				{Screenshot: dup},
				{Screenshot: dup},
				{Screenshot: "data:image/png;base64,BBBB"},
			},
		}},
	}
	saved := dedupeScreenshotURLs(&payload)
	if saved != 2 {
		t.Fatalf("expected 2 deduped refs, got %d", saved)
	}
	if payload.Scenarios[0].Steps[0].Screenshot != dup {
		t.Fatal("first step should keep canonical url")
	}
	if payload.Scenarios[0].Steps[1].Screenshot != dup {
		t.Fatal("duplicate should reference same url")
	}
}
