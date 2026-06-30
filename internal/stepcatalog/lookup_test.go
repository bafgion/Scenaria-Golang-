package stepcatalog

import "testing"

func TestLookupByActionClick(t *testing.T) {
	entry, ok := LookupByAction("click")
	if !ok || entry.Label != "нажимаю" {
		t.Fatalf("got %+v ok=%v", entry, ok)
	}
}

func TestLookupByActionNormalizesHyphen(t *testing.T) {
	entry, ok := LookupByAction("double-click")
	if !ok || entry.Action != "double_click" {
		t.Fatalf("got %+v ok=%v", entry, ok)
	}
}

func TestLookupByStepTextLabelPrefix(t *testing.T) {
	entry, ok := LookupByStepText(`нажимаю "#login"`)
	if !ok || entry.Label != "нажимаю" {
		t.Fatalf("got %+v ok=%v", entry, ok)
	}
}
