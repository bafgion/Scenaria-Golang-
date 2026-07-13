package secret

import (
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strings"
)

const Replacement = "[REDACTED]"

var (
	keyValuePattern = regexp.MustCompile(`(?i)\b(password|passwd|pwd|secret|token|authorization|api[_-]?key|cookie|credential|scenaria_email_code|email_code|http_auth)\b(\s*[:=]\s*)("[^"]*"|'[^']*'|[^\s,;]+)`)
	bearerPattern   = regexp.MustCompile(`(?i)\b(Bearer|Basic)\s+[A-Za-z0-9._~+/=-]+`)
	urlPattern      = regexp.MustCompile(`(?i)\b[a-z][a-z0-9+.-]*://[^\s"'<>]+`)
)

// RedactString removes accidental credentials from diagnostic text. It is not a
// storage transform; callers must still decide where secrets are allowed.
func RedactString(text string) string {
	if text == "" {
		return ""
	}
	redacted := urlPattern.ReplaceAllStringFunc(text, stripURLCredentials)
	redacted = bearerPattern.ReplaceAllString(redacted, `${1} `+Replacement)
	redacted = keyValuePattern.ReplaceAllString(redacted, `${1}${2}`+Replacement)
	return redacted
}

func RedactKeyValues(args []any) []any {
	if len(args) == 0 {
		return nil
	}
	out := make([]any, len(args))
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case slog.Attr:
			out[i] = redactAttr(v)
		case string:
			out[i] = RedactString(v)
			if i+1 < len(args) && IsSensitiveKey(v) {
				i++
				out[i] = Replacement
			}
		default:
			out[i] = RedactAny(v)
		}
	}
	return out
}

func RedactAny(value any) any {
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		return RedactString(v)
	case error:
		return RedactString(v.Error())
	case fmt.Stringer:
		return RedactString(v.String())
	default:
		return value
	}
}

func IsSensitiveKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" {
		return false
	}
	for _, marker := range []string{
		"password",
		"passwd",
		"pwd",
		"secret",
		"token",
		"authorization",
		"api_key",
		"apikey",
		"cookie",
		"credential",
		"scenaria_email_code",
		"email_code",
		"http_auth",
	} {
		if strings.Contains(key, marker) {
			return true
		}
	}
	return false
}

func redactAttr(attr slog.Attr) slog.Attr {
	if IsSensitiveKey(attr.Key) {
		return slog.Any(attr.Key, Replacement)
	}
	if attr.Value.Kind() == slog.KindString {
		return slog.String(attr.Key, RedactString(attr.Value.String()))
	}
	if attr.Value.Kind() == slog.KindAny {
		return slog.Any(attr.Key, RedactAny(attr.Value.Any()))
	}
	return attr
}

func stripURLCredentials(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil || parsed == nil || parsed.User == nil {
		return raw
	}
	parsed.User = nil
	return parsed.String()
}
