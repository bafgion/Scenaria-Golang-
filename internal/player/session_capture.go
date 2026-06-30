package player

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/settings"
	playwright "github.com/mxschmitt/playwright-go"
)

// CaptureTestClientFromPage exports cookies and localStorage from the active browser page.
func CaptureTestClientFromPage(page playwright.Page, name string) (*settings.TestClient, error) {
	if page == nil {
		return nil, fmt.Errorf("browser page is required")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("test client name is required")
	}

	ctx := page.Context()
	if ctx == nil {
		return nil, fmt.Errorf("browser context is required")
	}

	pageURL := strings.TrimSpace(page.URL())
	cookies, err := ctx.Cookies()
	if err != nil {
		return nil, fmt.Errorf("read browser cookies: %w", err)
	}

	localStorage, err := readPageLocalStorage(page)
	if err != nil {
		return nil, err
	}

	return &settings.TestClient{
		Name:         name,
		BaseURL:      originFromURL(pageURL),
		Cookies:      cookiesToSettings(cookies),
		LocalStorage: localStorage,
	}, nil
}

func readPageLocalStorage(page playwright.Page) (map[string]string, error) {
	raw, err := page.Evaluate(`() => {
		const out = {};
		try {
			for (let i = 0; i < localStorage.length; i++) {
				const key = localStorage.key(i);
				if (key) out[key] = localStorage.getItem(key) ?? '';
			}
		} catch (e) {
			return { __error: String(e) };
		}
		return out;
	}`)
	if err != nil {
		return nil, fmt.Errorf("read local storage: %w", err)
	}
	payload, err := decodeStringMap(raw)
	if err != nil {
		return nil, err
	}
	if msg := strings.TrimSpace(payload["__error"]); msg != "" {
		return nil, fmt.Errorf("read local storage: %s", msg)
	}
	delete(payload, "__error")
	if payload == nil {
		payload = map[string]string{}
	}
	return payload, nil
}

func decodeStringMap(raw any) (map[string]string, error) {
	if raw == nil {
		return map[string]string{}, nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("decode local storage: %w", err)
	}
	out := map[string]string{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("decode local storage: %w", err)
	}
	return out, nil
}

func cookiesToSettings(cookies []playwright.Cookie) []settings.Cookie {
	out := make([]settings.Cookie, 0, len(cookies))
	for _, cookie := range cookies {
		out = append(out, settings.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			HTTPOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
		})
	}
	return out
}

func originFromURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "about:blank" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
