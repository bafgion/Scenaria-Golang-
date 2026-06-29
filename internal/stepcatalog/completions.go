package stepcatalog

import (
	"regexp"
	"strings"
)

// CompletionSnippet is one autocomplete row for the Gherkin editor.
type CompletionSnippet struct {
	Label       string `json:"label"`
	Insert      string `json:"insert"`
	Description string `json:"description"`
}

// CompletionsResult is the replace range and matching snippets for one editor line.
type CompletionsResult struct {
	Start int                 `json:"start"` // 0-based column in line
	End   int                 `json:"end"`
	Items []CompletionSnippet `json:"items"`
}

var (
	completionKeywords = []string{"Допустим", "Дано", "Когда", "Тогда", "И", "Но"}
	headerLineRE       = regexp.MustCompile(`(?i)^\s*(функционал|сценарий|функция)\s*:?`)
	stepKeywordRE      = regexp.MustCompile(`(?i)^(?:(Допустим|Дано|Когда|Тогда|И|Но)\s+)?(.*)$`)
)

var headerSnippets = []CompletionSnippet{
	{Label: "Функционал:", Insert: "Функционал: UI сценарий", Description: "Заголовок feature-файла"},
	{Label: "Контекст:", Insert: "Контекст:\n\tДано я подключаю TestClient \"ИмяКлиента\"", Description: "Блок контекста: именованный TestClient перед сценарием"},
	{Label: "Сценарий:", Insert: "Сценарий: Имя сценария", Description: "Название сценария"},
}

func leadingIndent(line string) string {
	trimmed := strings.TrimLeft(line, " \t")
	return line[:len(line)-len(trimmed)]
}

func isStepIndented(line string) bool {
	indent := leadingIndent(line)
	if strings.HasPrefix(indent, "\t") {
		return true
	}
	return len(indent) >= 2 && strings.TrimSpace(indent) == ""
}

func isStepLine(line string) bool {
	if isStepIndented(line) {
		return true
	}
	stripped := strings.TrimLeft(line, " \t")
	lower := strings.ToLower(stripped)
	for _, kw := range completionKeywords {
		lkw := strings.ToLower(kw)
		if strings.HasPrefix(lower, lkw+" ") || lower == lkw {
			return true
		}
	}
	return false
}

func matchPrefix(word, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(word), strings.ToLower(prefix))
}

func keywordCandidates(prefix string) []CompletionSnippet {
	prefix = strings.TrimSpace(prefix)
	out := make([]CompletionSnippet, 0, len(completionKeywords))
	for _, word := range completionKeywords {
		if matchPrefix(word, prefix) {
			out = append(out, CompletionSnippet{
				Label:       word,
				Insert:      word,
				Description: "Ключевое слово «" + word + "»",
			})
		}
	}
	return out
}

func headerCandidates(prefix string) []CompletionSnippet {
	stripped := strings.TrimLeft(prefix, " \t")
	out := make([]CompletionSnippet, 0)
	for _, snip := range headerSnippets {
		if matchPrefix(snip.Label, stripped) || matchPrefix(snip.Insert, stripped) {
			out = append(out, snip)
		}
	}
	return out
}

func stepSnippetsForCompletion() []CompletionSnippet {
	out := make([]CompletionSnippet, len(stepSnippets))
	for i, snip := range stepSnippets {
		out[i] = CompletionSnippet{
			Label:       snip.label,
			Insert:      snip.insert,
			Description: plainDescription(snip.description),
		}
	}
	return out
}

func stepCandidates(prefix string) []CompletionSnippet {
	all := stepSnippetsForCompletion()
	out := make([]CompletionSnippet, 0)
	for _, snip := range all {
		if matchPrefix(snip.Label, prefix) || matchPrefix(snip.Insert, prefix) {
			out = append(out, snip)
		}
	}
	return out
}

func keywordEndColumn(indentLen int, keyword string) int {
	return indentLen + len(keyword)
}

// CompletionsForLine returns replace columns and snippets (Python completions_for_line parity).
func CompletionsForLine(line string, column int) CompletionsResult {
	empty := CompletionsResult{Start: column, End: column}
	if column < 0 {
		column = len(line)
	}
	if column > len(line) {
		column = len(line)
	}

	stripped := strings.TrimLeft(line, " \t")
	indentLen := len(line) - len(stripped)
	stepLine := isStepLine(line)

	if stripped == "" || strings.HasPrefix(stripped, "#") {
		if strings.TrimSpace(line) == "" && column >= indentLen {
			if stepLine || indentLen == 0 {
				items := append(stepSnippetsForCompletion(), keywordCandidates("")...)
				return CompletionsResult{Start: indentLen, End: column, Items: items}
			}
			return CompletionsResult{Start: indentLen, End: column, Items: keywordCandidates("")}
		}
		return empty
	}

	if headerLineRE.MatchString(stripped) && !isStepIndented(line) {
		return CompletionsResult{
			Start: indentLen,
			End:   column,
			Items: headerCandidates(stripped),
		}
	}

	if !stepLine {
		return empty
	}

	match := stepKeywordRE.FindStringSubmatch(stripped)
	if match == nil {
		return empty
	}
	keyword := match[1]
	body := match[2]
	bodyOffset := indentLen + len(stripped) - len(body)

	if keyword != "" {
		keywordEnd := keywordEndColumn(indentLen, keyword)
		if column <= keywordEnd {
			prefix := stripped[:max(0, column-indentLen)]
			return CompletionsResult{
				Start: indentLen,
				End:   column,
				Items: keywordCandidates(strings.TrimSpace(prefix)),
			}
		}
	} else {
		linePrefix := stripped[:max(0, column-indentLen)]
		if !strings.Contains(strings.TrimRight(linePrefix, " "), " ") {
			if matches := keywordCandidates(strings.TrimSpace(linePrefix)); len(matches) > 0 {
				return CompletionsResult{Start: indentLen, End: column, Items: matches}
			}
		}
	}

	bodyPrefix := strings.TrimLeft(line[bodyOffset:column], " \t")
	if bodyPrefix == "" {
		if keyword != "" {
			return CompletionsResult{Start: bodyOffset, End: column, Items: stepSnippetsForCompletion()}
		}
		return CompletionsResult{Start: indentLen, End: column, Items: keywordCandidates("")}
	}

	matches := stepCandidates(bodyPrefix)
	start := bodyOffset + len(body[:column-bodyOffset]) - len(bodyPrefix)
	return CompletionsResult{Start: start, End: column, Items: matches}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
