package gui

import (
	"strings"
)

// LineIsReplaceable reports whether find/replace may touch a line (ported from Python text_replace.py).
func LineIsReplaceable(line string, stepsOnly bool) bool {
	if !stepsOnly {
		return true
	}
	stripped := strings.TrimSpace(line)
	if stripped == "" || strings.HasPrefix(stripped, "#") {
		return false
	}
	if refactorHeaderRE.MatchString(stripped) || refactorTagRE.MatchString(stripped) || refactorExamplesRE.MatchString(stripped) {
		return false
	}
	return true
}

// ReplaceAllInText replaces needle in text; when stepsOnly is true, headers/tags/comments are skipped.
func ReplaceAllInText(text, find, replace string, caseSensitive, stepsOnly bool) RefactorResult {
	find = strings.TrimSpace(find)
	if find == "" {
		return RefactorResult{Text: text}
	}
	if !stepsOnly {
		return ReplaceInText(text, find, replace, caseSensitive)
	}

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	endsWithNL := strings.HasSuffix(normalized, "\n")
	lines := strings.Split(normalized, "\n")
	count := 0
	for i, body := range lines {
		if !LineIsReplaceable(body, true) {
			continue
		}
		replaced, n := replaceInLine(body, find, replace, caseSensitive)
		if n > 0 {
			lines[i] = replaced
			count += n
		}
	}
	out := strings.Join(lines, "\n")
	if endsWithNL {
		out += "\n"
	}
	return RefactorResult{Text: out, Count: count}
}

func replaceInLine(line, find, replace string, caseSensitive bool) (string, int) {
	if caseSensitive {
		n := strings.Count(line, find)
		if n == 0 {
			return line, 0
		}
		return strings.ReplaceAll(line, find, replace), n
	}
	lower := strings.ToLower(line)
	needle := strings.ToLower(find)
	var b strings.Builder
	b.Grow(len(line))
	count := 0
	for i := 0; i < len(line); {
		if strings.HasPrefix(lower[i:], needle) {
			b.WriteString(replace)
			i += len(find)
			count++
			continue
		}
		b.WriteByte(line[i])
		i++
	}
	return b.String(), count
}
