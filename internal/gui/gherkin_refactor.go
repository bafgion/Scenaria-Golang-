package gui

import (
	"regexp"
	"strings"
)

const gherkinStepIndent = "\t"

var (
	refactorKeywordRE = regexp.MustCompile(`(?i)^(Допустим|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+(.*)$`)
	refactorGotoBodyRE = regexp.MustCompile(`(?i)^открыт[а]?\s+"((?:\\.|[^"])*)"$`)
	refactorHeaderRE  = regexp.MustCompile(`(?i)^(функционал|функциональность|функция|feature|сценарий|scenario|структура\s+сценария|scenario\s+outline)\s*:`)
	refactorTagRE     = regexp.MustCompile(`^@\S+$`)
	refactorExamplesRE = regexp.MustCompile(`(?i)^примеры\s*:|^examples\s*:`)
	refactorBlockRE   = regexp.MustCompile(`(?i)^(если|повторяю|пока|для каждого|иначе|конец если|конец)\b`)
)

type RefactorResult struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
}

func UpdateStartURLs(text, newURL string) RefactorResult {
	newURL = strings.TrimSpace(newURL)
	if newURL == "" {
		return RefactorResult{Text: text}
	}
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	changed := 0
	for i, raw := range lines {
		stripped := strings.TrimSpace(raw)
		if !isRefactorStepBlockLine(stripped, raw) {
			continue
		}
		match := refactorKeywordRE.FindStringSubmatch(stripped)
		if match == nil {
			continue
		}
		body := strings.TrimSpace(match[2])
		if !refactorGotoBodyRE.MatchString(body) {
			continue
		}
		keyword := strings.TrimSpace(stripped[:len(stripped)-len(body)])
		prefix := keyword + " "
		if keyword == "" {
			prefix = ""
		}
		indent := leadingIndent(raw)
		lines[i] = indent + prefix + `открыт "` + newURL + `"`
		changed++
	}
	return RefactorResult{Text: joinLines(lines, text), Count: changed}
}

func ReplaceInText(text, find, replace string, caseSensitive bool) RefactorResult {
	find = strings.TrimSpace(find)
	if find == "" {
		return RefactorResult{Text: text}
	}
	var count int
	var out string
	if caseSensitive {
		count = strings.Count(text, find)
		out = strings.ReplaceAll(text, find, replace)
	} else {
		lower := strings.ToLower(text)
		needle := strings.ToLower(find)
		var b strings.Builder
		b.Grow(len(text))
		for i := 0; i < len(text); {
			if strings.HasPrefix(lower[i:], needle) {
				b.WriteString(replace)
				i += len(find)
				count++
				continue
			}
			b.WriteByte(text[i])
			i++
		}
		out = b.String()
	}
	return RefactorResult{Text: out, Count: count}
}

func NormalizeStepIndents(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	hasTabSteps := false
	for _, raw := range lines {
		stripped := strings.TrimSpace(raw)
		if isRefactorStepBlockLine(stripped, raw) && strings.HasPrefix(leadingIndent(raw), "\t") {
			hasTabSteps = true
			break
		}
	}
	out := make([]string, 0, len(lines))
	for _, raw := range lines {
		stripped := strings.TrimSpace(raw)
		if !isRefactorStepBlockLine(stripped, raw) {
			out = append(out, raw)
			continue
		}
		indent := leadingIndent(raw)
		level := 1
		if hasTabSteps {
			if strings.Contains(indent, "\t") {
				level = maxInt(1, strings.Count(indent, "\t"))
			} else if indent == "" || strings.TrimSpace(indent) == "" {
				if n := len(indent); n == 2 || n == 4 {
					level = 1
				} else if n >= 2 {
					level = maxInt(1, n/2)
				}
			}
		} else if strings.Contains(indent, "\t") {
			level = maxInt(1, strings.Count(indent, "\t"))
		} else if indent == "" || strings.TrimSpace(indent) == "" {
			if n := len(indent); n >= 2 {
				level = maxInt(1, n/2)
			}
		}
		out = append(out, strings.Repeat(gherkinStepIndent, level)+stripped)
	}
	return joinLines(out, text)
}

func CollapseBlankLinesBetweenSteps(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return text
	}
	result := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		raw := lines[i]
		stripped := strings.TrimSpace(raw)
		result = append(result, raw)
		if !isRefactorStepBlockLine(stripped, raw) {
			continue
		}
		next := i + 1
		for next < len(lines) && strings.TrimSpace(lines[next]) == "" {
			next++
		}
		if next < len(lines) && isRefactorStepBlockLine(strings.TrimSpace(lines[next]), lines[next]) {
			i = next - 1
		}
	}
	return joinLines(result, text)
}

// FormatFeature normalizes step indents and removes blank lines between consecutive steps.
func FormatFeature(text string) string {
	return CollapseBlankLinesBetweenSteps(NormalizeStepIndents(text))
}

func isRefactorStepBlockLine(stripped, raw string) bool {
	if stripped == "" || strings.HasPrefix(stripped, "#") {
		return false
	}
	if refactorHeaderRE.MatchString(stripped) || refactorTagRE.MatchString(stripped) || refactorExamplesRE.MatchString(stripped) {
		return false
	}
	if strings.HasPrefix(stripped, "|") {
		return false
	}
	if refactorKeywordRE.MatchString(stripped) {
		return true
	}
	return refactorBlockRE.MatchString(stripped) || isScenarioStructureLine(stripped)
}

func leadingIndent(raw string) string {
	return raw[:len(raw)-len(strings.TrimLeft(raw, " \t"))]
}

func joinLines(lines []string, original string) string {
	endsWithNL := strings.HasSuffix(original, "\n")
	if endsWithNL && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	payload := strings.Join(lines, "\n")
	if endsWithNL {
		payload += "\n"
	}
	return payload
}
