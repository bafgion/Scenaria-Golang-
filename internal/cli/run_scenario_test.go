package cli

import "testing"

func TestParseRunOptions_Scenario(t *testing.T) {
	opts, err := parseRunOptions([]string{"./features", "--scenario", "Успешный вход", "--dry-run"})
	if err != nil {
		t.Fatalf("parseRunOptions: %v", err)
	}
	if opts.scenario != "Успешный вход" || !opts.dryRun {
		t.Fatalf("unexpected opts: %+v", opts)
	}
}
