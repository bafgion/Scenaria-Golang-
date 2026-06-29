package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestFirstNavigationURL(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{{
			Steps: []gherkin.Step{{Line: 1, Text: `открыт "https://protected.example.com"`}},
		}},
	}
	if got := FirstNavigationURL(plan); got != "https://protected.example.com" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveAuthURLPrefersBaseURL(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{{
			Steps: []gherkin.Step{{Line: 1, Text: `открыт "https://other.example.com"`}},
		}},
	}
	if got := ResolveAuthURL("https://base.example.com", plan); got != "https://base.example.com" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveRunHTTPCredentialsFromSettings(t *testing.T) {
	cfg := &settings.AppSettings{
		HTTPAuth: map[string]settings.HTTPAuthEntry{
			"protected.example.com": {Username: "alice", Password: "secret"},
		},
	}
	plan := ExecutionPlan{
		Cases: []RunCase{{
			Steps: []gherkin.Step{{Line: 1, Text: `открыт "https://protected.example.com/login"`}},
		}},
	}
	creds := ResolveRunHTTPCredentials("", plan, cfg)
	if creds == nil || creds.Username != "alice" || creds.Password != "secret" {
		t.Fatalf("got %+v", creds)
	}
}
