package vanessa

import "testing"

func TestMergeVAParamsExcludeTags(t *testing.T) {
	t.Parallel()
	merged, _, err := MergeVAParams(DefaultSettings(), RunRequest{
		ExcludeTags: []string{"wip", "@skip"},
	}, t.TempDir())
	if err != nil {
		t.Fatalf("MergeVAParams: %v", err)
	}
	raw, ok := merged["СписокТеговИсключение"].([]string)
	if !ok {
		t.Fatalf("expected exclude tags list, got %#v", merged["СписокТеговИсключение"])
	}
	if len(raw) != 2 || raw[0] != "@wip" || raw[1] != "@skip" {
		t.Fatalf("unexpected exclude tags: %#v", raw)
	}
}
