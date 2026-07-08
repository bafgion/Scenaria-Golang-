package player

import "testing"

func TestParseNumberFromText(t *testing.T) {
	cases := map[string]string{
		"7 180 ₽":    "7180",
		"1 795":      "1795",
		"1\u00a0795": "1795",
		"Pay 1795":   "1795",
		"12 345.67":  "1234567",
	}
	for input, want := range cases {
		got, err := ParseNumberFromText(input)
		if err != nil {
			t.Fatalf("%q: %v", input, err)
		}
		if got != want {
			t.Fatalf("%q: got %q want %q", input, got, want)
		}
	}
	if _, err := ParseNumberFromText("no digits"); err == nil {
		t.Fatal("expected error for text without digits")
	}
}
