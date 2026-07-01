package player

import (
	"fmt"
	"os"
	"strings"
)

const maxPlaceholderPasses = 16

func (c *RunContext) ResolveText(text string) (string, error) {
	if text == "" || !strings.Contains(text, "{{") {
		return text, nil
	}
	out := text
	for pass := 0; pass < maxPlaceholderPasses; pass++ {
		resolved, changed, err := resolvePlaceholdersOnce(c, out)
		if err != nil {
			return "", err
		}
		if !changed {
			if strings.Contains(resolved, "{{") && strings.Contains(resolved, "}}") {
				return "", fmt.Errorf("unresolved placeholder in %q", text)
			}
			return resolved, nil
		}
		out = resolved
	}
	return "", fmt.Errorf("placeholder resolution exceeded %d passes", maxPlaceholderPasses)
}

func resolvePlaceholdersOnce(c *RunContext, text string) (string, bool, error) {
	var b strings.Builder
	b.Grow(len(text))
	changed := false
	for i := 0; i < len(text); {
		open := strings.Index(text[i:], "{{")
		if open < 0 {
			b.WriteString(text[i:])
			break
		}
		open += i
		b.WriteString(text[i:open])
		closeRel := strings.Index(text[open+2:], "}}")
		if closeRel < 0 {
			b.WriteString(text[open:])
			break
		}
		close := open + 2 + closeRel
		key := strings.TrimSpace(text[open+2 : close])
		if key == "" {
			return "", false, fmt.Errorf("empty placeholder in %q", text)
		}
		value, err := c.resolvePlaceholderKey(key)
		if err != nil {
			return "", false, err
		}
		b.WriteString(value)
		changed = true
		i = close + 2
	}
	return b.String(), changed, nil
}

func (c *RunContext) resolvePlaceholderKey(key string) (string, error) {
	if strings.HasPrefix(strings.ToLower(key), "env:") {
		envName := key[4:]
		value := os.Getenv(envName)
		if value == "" {
			return "", fmt.Errorf("environment variable %q is empty", envName)
		}
		return value, nil
	}
	if value, ok := c.Variables[key]; ok {
		if strings.Contains(value, "{{") {
			return c.resolveNestedValue(key, value)
		}
		return value, nil
	}
	if value, ok := c.values[key]; ok {
		if strings.Contains(value, "{{") {
			return c.resolveNestedValue(key, value)
		}
		return value, nil
	}
	generated, err := c.generate(key)
	if err != nil {
		return "", err
	}
	c.values[key] = generated
	return generated, nil
}

func (c *RunContext) resolveNestedValue(key, value string) (string, error) {
	if c.resolving[key] {
		return "", fmt.Errorf("circular placeholder reference for %q", key)
	}
	c.resolving[key] = true
	defer delete(c.resolving, key)
	return c.ResolveText(value)
}
