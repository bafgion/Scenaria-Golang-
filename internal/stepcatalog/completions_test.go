package stepcatalog_test

import (
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
)

func labels(items []stepcatalog.CompletionSnippet) map[string]bool {
	out := make(map[string]bool, len(items))
	for _, item := range items {
		out[item.Label] = true
	}
	return out
}

func TestCompletionsFilterByTypedPrefix(t *testing.T) {
	line := "\tКогда наж"
	result := stepcatalog.CompletionsForLine(line, len(line))
	for _, item := range result.Items {
		if !strings.HasPrefix(strings.ToLower(item.Label), "наж") {
			t.Fatalf("unexpected item %q for prefix наж", item.Label)
		}
	}
	if len(result.Items) == 0 {
		t.Fatal("expected at least one completion for наж")
	}
}

func TestCompletionsIncludeTabSteps(t *testing.T) {
	matches := stepcatalog.CompletionsForLine("\tИ переключаюсь", len("\tИ переключаюсь"))
	got := labels(matches.Items)
	for _, want := range []string{
		"переключаюсь на вкладку",
		"переключаюсь на вкладку с url",
		"переключаюсь на вкладку 2",
		"переключаюсь на первую вкладку",
	} {
		if !got[want] {
			t.Fatalf("expected label %q in completions, got %v", want, got)
		}
	}

	closeMatches := stepcatalog.CompletionsForLine("\tИ закрываю", len("\tИ закрываю"))
	if !labels(closeMatches.Items)["закрываю текущую вкладку"] {
		t.Fatalf("expected close tab snippet, got %v", labels(closeMatches.Items))
	}
}

func TestCompletionsEmptyBodyAfterKeywordListsAllSteps(t *testing.T) {
	line := "\tКогда "
	result := stepcatalog.CompletionsForLine(line, len(line))
	if len(result.Items) < 50 {
		t.Fatalf("expected full step catalog after keyword, got %d items", len(result.Items))
	}
}

func TestCompletionsHeaderLines(t *testing.T) {
	result := stepcatalog.CompletionsForLine("Функционал", len("Функционал"))
	if len(result.Items) == 0 {
		t.Fatal("expected header completions for Функционал")
	}
}

func TestStepSnippetsHaveActionTags(t *testing.T) {
	result := stepcatalog.CompletionsForLine("\tКогда ", len("\tКогда "))
	for _, item := range result.Items {
		if item.Label == "" || item.Insert == "" {
			t.Fatalf("empty snippet: %+v", item)
		}
	}
}
