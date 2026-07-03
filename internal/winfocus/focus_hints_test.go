package winfocus

import "testing"

func TestScoreChromiumTitle(t *testing.T) {
	score := scoreChromiumTitle("2moodstore — интернет-магазин", "2moodstore", "https://www.2moodstore.com/")
	if score < 100 {
		t.Fatalf("expected title match, got %d", score)
	}
	score = scoreChromiumTitle("example.com - Home", "", "https://example.com/page")
	if score < 90 {
		t.Fatalf("expected host match, got %d", score)
	}
}

func TestHostFromURL(t *testing.T) {
	if got := hostFromURL("https://www.2moodstore.com/"); got != "www.2moodstore.com" {
		t.Fatalf("unexpected host: %q", got)
	}
}
