package httpauth

import (
	"net/url"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

type Credentials struct {
	Username string
	Password string
	Origin   string
}

func HostFromURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Hostname())
}

func ParseURLCredentials(raw string) (username, password, clean string, ok bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || !strings.Contains(raw, "://") {
		return "", "", raw, false
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.User == nil {
		return "", "", raw, false
	}
	username = parsed.User.Username()
	password, _ = parsed.User.Password()
	if username == "" {
		return "", "", raw, false
	}
	parsed.User = nil
	return username, password, parsed.String(), true
}

func StripURLCredentials(raw string) string {
	_, _, clean, ok := ParseURLCredentials(raw)
	if !ok {
		return raw
	}
	return clean
}

func ResolveCredentials(raw string, cfg *settings.AppSettings) *Credentials {
	if username, password, _, ok := ParseURLCredentials(raw); ok && username != "" {
		return &Credentials{Username: username, Password: password}
	}
	host := HostFromURL(raw)
	if host == "" || cfg == nil || len(cfg.HTTPAuth) == 0 {
		return nil
	}
	entry, exists := cfg.HTTPAuth[host]
	if !exists || strings.TrimSpace(entry.Username) == "" {
		return nil
	}
	return &Credentials{Username: entry.Username, Password: entry.Password}
}

func OriginForURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	host := parsed.Hostname()
	if host == "" {
		return ""
	}
	scheme := parsed.Scheme
	if scheme == "" {
		scheme = "https"
	}
	if parsed.Port() != "" {
		port := parsed.Port()
		if (scheme == "https" && port == "443") || (scheme == "http" && port == "80") {
			return scheme + "://" + host
		}
		return scheme + "://" + host + ":" + port
	}
	return scheme + "://" + host
}

func PlaywrightHTTPCredentials(raw string, cfg *settings.AppSettings) *playwright.HttpCredentials {
	creds := ResolveCredentials(raw, cfg)
	if creds == nil {
		return nil
	}
	origin := OriginForURL(raw)
	if origin == "" {
		origin = OriginForURL("https://" + HostFromURL(raw))
	}
	out := &playwright.HttpCredentials{
		Username: creds.Username,
		Password: creds.Password,
	}
	if origin != "" {
		out.Origin = playwright.String(origin)
	}
	return out
}

func StoreHostCredentials(host, username, password string, cfg *settings.AppSettings) {
	if cfg == nil {
		return
	}
	host = strings.ToLower(strings.TrimSpace(host))
	if cfg.HTTPAuth == nil {
		cfg.HTTPAuth = map[string]settings.HTTPAuthEntry{}
	}
	if host == "" || strings.TrimSpace(username) == "" {
		delete(cfg.HTTPAuth, host)
		return
	}
	cfg.HTTPAuth[host] = settings.HTTPAuthEntry{
		Username: strings.TrimSpace(username),
		Password: password,
	}
}

func CredentialsForHost(host string, cfg *settings.AppSettings) (username, password string) {
	if cfg == nil || len(cfg.HTTPAuth) == 0 {
		return "", ""
	}
	entry, ok := cfg.HTTPAuth[strings.ToLower(strings.TrimSpace(host))]
	if !ok {
		return "", ""
	}
	return entry.Username, entry.Password
}

func ListHosts(cfg *settings.AppSettings) []string {
	if cfg == nil || len(cfg.HTTPAuth) == 0 {
		return nil
	}
	out := make([]string, 0, len(cfg.HTTPAuth))
	for host, entry := range cfg.HTTPAuth {
		if strings.TrimSpace(entry.Username) != "" {
			out = append(out, host)
		}
	}
	sortStrings(out)
	return out
}

func RemoveHostCredentials(host string, cfg *settings.AppSettings) {
	if cfg == nil || len(cfg.HTTPAuth) == 0 {
		return
	}
	delete(cfg.HTTPAuth, strings.ToLower(strings.TrimSpace(host)))
}

func ApplyURLCredentials(raw string, cfg *settings.AppSettings) string {
	username, password, clean, ok := ParseURLCredentials(raw)
	if !ok {
		return raw
	}
	host := HostFromURL(clean)
	if host != "" && username != "" {
		StoreHostCredentials(host, username, password, cfg)
	}
	return clean
}

func sortStrings(values []string) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}
