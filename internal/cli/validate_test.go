package cli

import "testing"

func TestParseValidateOptionsFlagsBeforeTarget(t *testing.T) {
	opts, err := parseValidateOptions([]string{"--no-browser", "./features"})
	if err != nil {
		t.Fatalf("parseValidateOptions: %v", err)
	}
	if !opts.noBrowser || len(opts.targets) != 1 || opts.targets[0] != "./features" {
		t.Fatalf("unexpected opts: %+v", opts)
	}
}
