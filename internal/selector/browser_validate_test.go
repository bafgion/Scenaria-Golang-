package selector

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestFirstGotoURLFromFeatureStep(t *testing.T) {
	feature, err := gherkin.ParseFeature(`Функционал: smoke
  Сценарий: test
    Допустим открыт "https://example.com/page"
    И нажимаю "#ok"
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := firstGotoURL(feature, "", "")
	if got != "https://example.com/page" {
		t.Fatalf("got %q", got)
	}
}

func TestFirstGotoURLFallsBackToBaseURL(t *testing.T) {
	feature, err := gherkin.ParseFeature(`Функционал: smoke
  Сценарий: test
    Допустим нажимаю "#ok"
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := firstGotoURL(feature, "https://base.test/", "")
	if got != "https://base.test/" {
		t.Fatalf("got %q", got)
	}
}

func TestFirstGotoURLEmptyWithoutGotoOrBase(t *testing.T) {
	feature, err := gherkin.ParseFeature(`Функционал: smoke
  Сценарий: test
    Допустим нажимаю "#ok"
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got := firstGotoURL(feature, "", ""); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestFirstGotoURLResolvesRelativeURL(t *testing.T) {
	feature, err := gherkin.ParseFeature(`Функционал: smoke
  Сценарий: test
    Допустим открыт "/login"
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := firstGotoURL(feature, "https://app.test", "")
	if got != "https://app.test/login" {
		t.Fatalf("got %q", got)
	}
}
