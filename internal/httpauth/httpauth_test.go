package httpauth

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestParseURLCredentials(t *testing.T) {
	user, pass, clean, ok := ParseURLCredentials(`https://user:secret@example.com/path`)
	if !ok || user != "user" || pass != "secret" {
		t.Fatalf("got %q %q ok=%v", user, pass, ok)
	}
	if clean != "https://example.com/path" {
		t.Fatalf("clean=%q", clean)
	}
}

func TestStoreAndResolveHostCredentials(t *testing.T) {
	cfg := &settings.AppSettings{}
	StoreHostCredentials("Example.COM", "alice", "pw", cfg)
	creds := ResolveCredentials("https://example.com/page", cfg)
	if creds == nil || creds.Username != "alice" || creds.Password != "pw" {
		t.Fatalf("got %+v", creds)
	}
}

func TestApplyURLCredentials(t *testing.T) {
	cfg := &settings.AppSettings{}
	clean := ApplyURLCredentials(`https://bob:123@site.test/`, cfg)
	if clean != "https://site.test/" {
		t.Fatalf("clean=%q", clean)
	}
	user, pass := CredentialsForHost("site.test", cfg)
	if user != "bob" || pass != "123" {
		t.Fatalf("stored %q %q", user, pass)
	}
}
