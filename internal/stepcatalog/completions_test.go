package stepcatalog_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
)

func runeCol(line string) int {
	return utf8.RuneCountInString(line)
}

func labels(items []stepcatalog.CompletionSnippet) map[string]bool {
	out := make(map[string]bool, len(items))
	for _, item := range items {
		out[item.Label] = true
	}
	return out
}

func TestCompletionsFilterByTypedPrefix(t *testing.T) {
	line := "\tКогда наж"
	result := stepcatalog.CompletionsForLine(line, runeCol(line))
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
	matches := stepcatalog.CompletionsForLine("\tИ переключаюсь", runeCol("\tИ переключаюсь"))
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

	closeMatches := stepcatalog.CompletionsForLine("\tИ закрываю", runeCol("\tИ закрываю"))
	if !labels(closeMatches.Items)["закрываю текущую вкладку"] {
		t.Fatalf("expected close tab snippet, got %v", labels(closeMatches.Items))
	}
}

func TestCompletionsEmptyBodyAfterKeywordListsAllSteps(t *testing.T) {
	line := "\tКогда "
	result := stepcatalog.CompletionsForLine(line, runeCol(line))
	if len(result.Items) < 50 {
		t.Fatalf("expected full step catalog after keyword, got %d items", len(result.Items))
	}
}

func TestCompletionsHeaderLines(t *testing.T) {
	result := stepcatalog.CompletionsForLine("Функционал", runeCol("Функционал"))
	if len(result.Items) == 0 {
		t.Fatal("expected header completions for Функционал")
	}
}

func TestCompletionsCyrillicTypedPrefix(t *testing.T) {
	line := "\tИ в"
	result := stepcatalog.CompletionsForLine(line, runeCol(line))
	got := labels(result.Items)
	if !got["ввожу"] {
		t.Fatalf("expected ввожу for prefix в, got %v", got)
	}
	for _, item := range result.Items {
		if item.Label == "И" {
			t.Fatalf("keyword-only completion should not appear when body prefix is typed")
		}
	}
}

func TestStepSnippetsHaveActionTags(t *testing.T) {
	result := stepcatalog.CompletionsForLine("\tКогда ", runeCol("\tКогда "))
	for _, item := range result.Items {
		if item.Label == "" || item.Insert == "" {
			t.Fatalf("empty snippet: %+v", item)
		}
	}
}
