package recorder

import (
	"encoding/json"
	"fmt"
)

// SelectorCandidate describes one ranked picker selector option.
type SelectorCandidate struct {
	Selector      string   `json:"selector"`
	Strategy      string   `json:"strategy"`
	Score         int      `json:"score"`
	MatchesCount  int      `json:"matches_count"`
	Unique        bool     `json:"unique"`
	Visible       bool     `json:"visible"`
	MatchesPicked bool     `json:"matches_picked"`
	Warnings      []string `json:"warnings,omitempty"`
}

// PickPayload is the structured result from the in-browser element picker.
type PickPayload struct {
	Selector        string              `json:"selector"`
	SuggestedAction string              `json:"suggested_action,omitempty"`
	Warnings        []string            `json:"warnings,omitempty"`
	Candidates      []SelectorCandidate `json:"candidates,omitempty"`
}

func parsePickPayload(arg any) (PickPayload, error) {
	switch v := arg.(type) {
	case nil:
		return PickPayload{}, fmt.Errorf("empty picker payload")
	case string:
		if v == "" {
			return PickPayload{}, fmt.Errorf("empty picker selector")
		}
		return PickPayload{Selector: v}, nil
	case map[string]any:
		return pickPayloadFromMap(v)
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return PickPayload{}, fmt.Errorf("picker payload: %w", err)
		}
		var payload PickPayload
		if err := json.Unmarshal(raw, &payload); err != nil {
			return PickPayload{}, fmt.Errorf("picker payload: %w", err)
		}
		if payload.Selector == "" {
			return PickPayload{}, fmt.Errorf("empty picker selector")
		}
		return payload, nil
	}
}

func pickPayloadFromMap(m map[string]any) (PickPayload, error) {
	selector, _ := m["selector"].(string)
	if selector == "" {
		return PickPayload{}, fmt.Errorf("empty picker selector")
	}
	payload := PickPayload{Selector: selector}
	if action, _ := m["suggested_action"].(string); action != "" {
		payload.SuggestedAction = action
	}
	if warnings, ok := m["warnings"].([]any); ok {
		for _, item := range warnings {
			if s, ok := item.(string); ok && s != "" {
				payload.Warnings = append(payload.Warnings, s)
			}
		}
	}
	if rawCandidates, ok := m["candidates"].([]any); ok {
		for _, item := range rawCandidates {
			candMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			cand := SelectorCandidate{}
			if s, _ := candMap["selector"].(string); s != "" {
				cand.Selector = s
			}
			if s, _ := candMap["strategy"].(string); s != "" {
				cand.Strategy = s
			}
			if n, ok := asInt(candMap["score"]); ok {
				cand.Score = n
			}
			if n, ok := asInt(candMap["matches_count"]); ok {
				cand.MatchesCount = n
			}
			if b, ok := candMap["unique"].(bool); ok {
				cand.Unique = b
			}
			if b, ok := candMap["visible"].(bool); ok {
				cand.Visible = b
			}
			if b, ok := candMap["matches_picked"].(bool); ok {
				cand.MatchesPicked = b
			}
			if warnings, ok := candMap["warnings"].([]any); ok {
				for _, w := range warnings {
					if s, ok := w.(string); ok && s != "" {
						cand.Warnings = append(cand.Warnings, s)
					}
				}
			}
			if cand.Selector != "" {
				payload.Candidates = append(payload.Candidates, cand)
			}
		}
	}
	return payload, nil
}

func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	default:
		return 0, false
	}
}
