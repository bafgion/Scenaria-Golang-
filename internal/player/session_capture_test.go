package player

import (
	"testing"

	"github.com/mxschmitt/playwright-go"
)

func TestOriginFromURL(t *testing.T) {
	tests := map[string]string{
		"https://app.example.com/path?q=1": "https://app.example.com",
		"http://localhost:8080/":         "http://localhost:8080",
		"about:blank":                      "",
		"":                                 "",
		"not-a-url":                        "",
	}
	for input, want := range tests {
		if got := originFromURL(input); got != want {
			t.Fatalf("originFromURL(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestCookiesToSettings(t *testing.T) {
	got := cookiesToSettings([]playwright.Cookie{
		{Name: "sid", Value: "abc", Domain: "example.com", Path: "/", HttpOnly: true, Secure: true},
	})
	if len(got) != 1 || got[0].Name != "sid" || !got[0].HTTPOnly || !got[0].Secure {
		t.Fatalf("unexpected cookies: %+v", got)
	}
}

func TestDecodeStringMap(t *testing.T) {
	got, err := decodeStringMap(map[string]any{"token": "1", "lang": "ru"})
	if err != nil {
		t.Fatal(err)
	}
	if got["token"] != "1" || got["lang"] != "ru" {
		t.Fatalf("unexpected map: %+v", got)
	}
}
