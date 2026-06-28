package gherkin

import (
	"regexp"
	"strings"
)

var (
	fancyQuotes = strings.NewReplacer(
		"\u201c", `"`, "\u201d", `"`, "\u201e", `"`, "\u00ab", `"`, "\u00bb", `"`,
		"\u2018", `'`, "\u2019", `'`,
	)
	zeroWidthRE = regexp.MustCompile(`[\x{200B}-\x{200D}\x{FEFF}]`)
	legacyHasTextRE = regexp.MustCompile(`:has-text\("((?:[^"\\]|\\.)*)"\)`)
)

func NormalizeFeatureText(content string) string {
	content = fancyQuotes.Replace(content)
	content = zeroWidthRE.ReplaceAllString(content, "")
	return content
}

func NormalizeLegacyHasTextEscapes(selector string) string {
	return legacyHasTextRE.ReplaceAllStringFunc(selector, func(match string) string {
		groups := legacyHasTextRE.FindStringSubmatch(match)
		if groups == nil {
			return match
		}
		inner := strings.ReplaceAll(groups[1], `\"`, `"`)
		escaped := strings.ReplaceAll(inner, `"`, `\"`)
		return `:has-text("` + escaped + `")`
	})
}

func CoalesceMixedIndents(content string) string {
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "    ") && !strings.HasPrefix(line, "\t") {
			trimmed := strings.TrimLeft(line, " ")
			spaces := len(line) - len(trimmed)
			tabs := spaces / 4
			if spaces%4 == 2 {
				tabs++
			}
			line = strings.Repeat("\t", tabs) + trimmed
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
