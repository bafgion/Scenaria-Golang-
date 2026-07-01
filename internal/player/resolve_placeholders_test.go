package player

import (
	"testing"
)

func TestResolveNestedPlaceholders(t *testing.T) {
	ctx := NewRunContext(map[string]string{
		"a": "{{b}}",
		"b": "done",
	}, 1, "")
	got, err := ctx.ResolveText("prefix {{a}} suffix")
	if err != nil || got != "prefix done suffix" {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestResolveMultiplePlaceholdersInOnePass(t *testing.T) {
	ctx := NewRunContext(map[string]string{
		"first": "Ada",
		"last":  "Lovelace",
	}, 2, "")
	got, err := ctx.ResolveText("{{first}} {{last}}")
	if err != nil || got != "Ada Lovelace" {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestResolveLeavesUnclosedPlaceholderLiteral(t *testing.T) {
	ctx := NewRunContext(nil, 3, "")
	got, err := ctx.ResolveText(`value {{open`)
	if err != nil || got != `value {{open` {
		t.Fatalf("ResolveText() = %q, %v", got, err)
	}
}

func TestAllocateDownloadPathUnique(t *testing.T) {
	ctx := NewRunContext(nil, 42, t.TempDir())
	p1 := ctx.allocateDownloadPath("report.pdf")
	p2 := ctx.allocateDownloadPath("report.pdf")
	if p1 == p2 {
		t.Fatalf("expected unique paths, both %q", p1)
	}
}

func TestSanitizeDownloadName(t *testing.T) {
	got := sanitizeDownloadName(`weird/name?.pdf`)
	if got == "" || got == "weird/name?.pdf" {
		t.Fatalf("unexpected sanitize result: %q", got)
	}
}
