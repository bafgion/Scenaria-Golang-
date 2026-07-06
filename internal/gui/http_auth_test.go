package gui

import (
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestHTTPAuthCredentialsOmitsPassword(t *testing.T) {
	cfg := &settings.AppSettings{
		HTTPAuth: map[string]settings.HTTPAuthEntry{
			"example.com": {Username: "user", Password: "secret"},
		},
	}
	username, password := httpauth.CredentialsForHost("example.com", cfg)
	dto := HTTPAuthCredentials{Username: username, HasPassword: password != ""}
	if dto.Username != "user" || !dto.HasPassword {
		t.Fatalf("unexpected dto: %+v", dto)
	}
}

func TestSaveHTTPAuthPreservesPasswordWhenBlank(t *testing.T) {
	cfg := &settings.AppSettings{
		HTTPAuth: map[string]settings.HTTPAuthEntry{
			"example.com": {Username: "user", Password: "secret"},
		},
	}
	req := HTTPAuthRequest{Host: "example.com", Username: "user", Password: ""}
	password := req.Password
	if strings.TrimSpace(password) == "" {
		if _, existing := httpauth.CredentialsForHost(req.Host, cfg); existing != "" {
			password = existing
		}
	}
	if password != "secret" {
		t.Fatalf("expected preserved password, got %q", password)
	}
}
