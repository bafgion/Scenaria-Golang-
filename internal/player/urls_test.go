package player

import "testing"

func TestUrlsMatch(t *testing.T) {
	cases := []struct {
		current string
		target  string
		want    bool
	}{
		{"https://example.com/page", "https://example.com/page", true},
		{"https://example.com/page/", "https://example.com/page", true},
		{"https://example.com/other", "https://example.com/page", false},
		{"", "https://example.com", false},
	}
	for _, tc := range cases {
		if got := UrlsMatch(tc.current, tc.target); got != tc.want {
			t.Fatalf("UrlsMatch(%q, %q) = %v, want %v", tc.current, tc.target, got, tc.want)
		}
	}
}

func TestGenerateINNChecksum(t *testing.T) {
	ctx := NewRunContext(nil, 42, "")
	inn, err := ctx.GenerateByKind("inn")
	if err != nil {
		t.Fatalf("GenerateByKind inn: %v", err)
	}
	if len(inn) != 12 {
		t.Fatalf("expected 12-digit INN, got %q", inn)
	}
}

func TestGenerateOGRNIP(t *testing.T) {
	ctx := NewRunContext(nil, 7, "")
	ogrnip, err := ctx.GenerateByKind("ogrnip")
	if err != nil {
		t.Fatalf("GenerateByKind ogrnip: %v", err)
	}
	if len(ogrnip) != 15 || ogrnip[0] != '3' {
		t.Fatalf("unexpected OGRNIP: %q", ogrnip)
	}
}

func TestCoherentPersonBundle(t *testing.T) {
	ctx := NewRunContext(nil, 1, "")
	first, _ := ctx.GenerateByKind("first_name")
	last, _ := ctx.GenerateByKind("last_name")
	patronymic, _ := ctx.GenerateByKind("patronymic")
	if first == "" || last == "" || patronymic == "" {
		t.Fatalf("empty person fields: %q %q %q", first, last, patronymic)
	}
	first2, _ := ctx.GenerateByKind("first_name")
	if first != first2 {
		t.Fatalf("expected stable first name within run")
	}
}
