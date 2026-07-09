package recorder

import "testing"

func TestParsePickPayloadString(t *testing.T) {
	payload, err := parsePickPayload(`#login`)
	if err != nil || payload.Selector != `#login` {
		t.Fatalf("got %+v err=%v", payload, err)
	}
}

func TestParsePickPayloadMap(t *testing.T) {
	payload, err := parsePickPayload(map[string]any{
		"selector":         `[data-testid="save"]`,
		"suggested_action": "click",
		"warnings":         []any{"low-confidence"},
		"candidates": []any{
			map[string]any{
				"selector":       `[data-testid="save"]`,
				"strategy":       "testid",
				"score":          float64(0),
				"matches_count":  float64(1),
				"unique":         true,
				"visible":        true,
				"matches_picked": true,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if payload.Selector != `[data-testid="save"]` || payload.SuggestedAction != "click" {
		t.Fatalf("payload: %+v", payload)
	}
	if len(payload.Candidates) != 1 || payload.Candidates[0].Strategy != "testid" {
		t.Fatalf("candidates: %+v", payload.Candidates)
	}
}
