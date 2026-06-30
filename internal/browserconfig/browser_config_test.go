package browserconfig

import "testing"

func TestLaunchOptionsHeadedChromiumMaximized(t *testing.T) {
	opts := LaunchOptions("chromium", false, 0)
	if opts.Headless == nil || *opts.Headless {
		t.Fatal("expected headed launch")
	}
	found := false
	for _, arg := range opts.Args {
		if arg == "--start-maximized" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected --start-maximized in %v", opts.Args)
	}
}

func TestNewContextOptionsHeadedNoViewport(t *testing.T) {
	opts := NewContextOptions(false, nil)
	if opts.NoViewport == nil || !*opts.NoViewport {
		t.Fatal("headed context should disable fixed viewport")
	}
}
