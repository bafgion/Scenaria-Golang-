package stepcatalog

import (
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func TestStepDSLGoldenKindsHaveCatalogEntries(t *testing.T) {
	cases, err := stepdsl.LoadGoldenCases()
	if err != nil {
		t.Fatalf("load golden cases: %v", err)
	}
	seenKinds := map[string]struct{}{}
	for _, tc := range cases {
		seenKinds[tc.Kind] = struct{}{}
	}
	for kind := range seenKinds {
		if _, ok := LookupByAction(kind); !ok {
			t.Fatalf("missing catalog entry for stepdsl kind %q", kind)
		}
	}
}

func TestCompletionSnippetsAreExecutableOrSnippetOnly(t *testing.T) {
	for _, snip := range append([]snippetDef{}, stepSnippets...) {
		firstLine := strings.Split(snip.insert, "\n")[0]
		if _, err := stepdsl.Parse(gherkin.Step{Line: 1, Text: firstLine}); err == nil {
			continue
		}
		action := actionFromDescription(snip.description)
		if action == "if" || action == "repeat" || action == "while" || action == "for_each" {
			continue
		}
		t.Fatalf("completion %q is not executable and not snippet-only", snip.label)
	}
	for _, snip := range append([]snippetDef{}, stepSnippetsEN...) {
		firstLine := strings.Split(snip.insert, "\n")[0]
		if _, err := stepdsl.Parse(gherkin.Step{Line: 1, Text: firstLine}); err == nil {
			continue
		}
		action := actionFromDescription(snip.description)
		if action == "if" || action == "repeat" || action == "while" || action == "for_each" {
			continue
		}
		t.Fatalf("english completion %q is not executable and not snippet-only", snip.label)
	}
}
