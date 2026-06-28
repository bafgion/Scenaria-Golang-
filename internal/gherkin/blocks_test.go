package gherkin

import "testing"

func TestDetectBlockHeaders(t *testing.T) {
	ifStep := Step{Line: 1, Text: `Если вижу "#modal"`}
	header, err := detectBlockHeader(ifStep)
	if err != nil || header == nil || header.Kind != BlockIf {
		t.Fatalf("unexpected if header: %+v err=%v", header, err)
	}

	repeatStep := Step{Line: 2, Text: "Повторяю 3 раза"}
	header, err = detectBlockHeader(repeatStep)
	if err != nil || header == nil || header.Kind != BlockRepeat || header.Count != 3 {
		t.Fatalf("unexpected repeat header: %+v err=%v", header, err)
	}
}

func TestBuildStepTree_WithNestedIf(t *testing.T) {
	flat := []Step{
		{Line: 1, Text: `Если вижу "#modal"`, Indent: 0},
		{Line: 2, Text: `нажимаю "#ok"`, Indent: 1},
		{Line: 3, Text: `вижу "#done"`, Indent: 0},
	}
	tree, err := buildStepTree(flat)
	if err != nil {
		t.Fatalf("buildStepTree failed: %v", err)
	}
	if len(tree) != 2 || tree[0].Block != BlockIf || len(tree[0].Children) != 1 {
		t.Fatalf("unexpected tree: %+v", tree)
	}
}

func TestParseTestClientName(t *testing.T) {
	name, err := ParseTestClientName([]Step{{Line: 1, Text: `я подключаю TestClient "Demo"`}})
	if err != nil || name != "Demo" {
		t.Fatalf("unexpected test client: %q err=%v", name, err)
	}
}
