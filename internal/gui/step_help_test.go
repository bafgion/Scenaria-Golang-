package gui

import "testing"

func TestDescribeEditorLineStep(t *testing.T) {
	entry, ok := DescribeEditorLine(`	Когда нажимаю "#login"`)
	if !ok || entry.Label != "нажимаю" {
		t.Fatalf("got %+v ok=%v", entry, ok)
	}
	if entry.Example == "" {
		t.Fatal("expected example")
	}
}

func TestDescribeEditorLineSkipsHeader(t *testing.T) {
	_, ok := DescribeEditorLine("  Сценарий: Demo")
	if ok {
		t.Fatal("expected no help for scenario header")
	}
}

func TestDescribeEditorLineSkipsComment(t *testing.T) {
	_, ok := DescribeEditorLine("  # comment")
	if ok {
		t.Fatal("expected no help for comment")
	}
}
