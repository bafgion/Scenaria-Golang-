package stepcatalog

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// CompletionSnippet is one autocomplete row for the Gherkin editor.
type CompletionSnippet struct {
	Label       string `json:"label"`
	Insert      string `json:"insert"`
	Description string `json:"description"`
}

// CompletionsResult is the replace range and matching snippets for one editor line.
// Start and End are 0-based rune indices in line (not UTF-8 byte offsets).
type CompletionsResult struct {
	Start int                 `json:"start"`
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

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

func leadingIndentRunes(line []rune) int {
	n := 0
	for n < len(line) && (line[n] == ' ' || line[n] == '\t') {
		n++
	}
	return n
}

func isStepIndentedRunes(line []rune, indentLen int) bool {
	if indentLen > 0 && line[0] == '\t' {
		return true
	}
	if indentLen >= 2 {
		spaces := true
		for i := 0; i < indentLen && i < len(line); i++ {
			if line[i] != ' ' {
				spaces = false
				break
			}
		}
		return spaces
	}
	return false
}

func isStepLineRunes(line []rune, indentLen int) bool {
	if isStepIndentedRunes(line, indentLen) {
		return true
	}
	if indentLen >= len(line) {
		return false
	}
	stripped := strings.ToLower(string(line[indentLen:]))
	for _, kw := range completionKeywords {
		lkw := strings.ToLower(kw)
		if strings.HasPrefix(stripped, lkw+" ") || stripped == lkw {
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

func keywordEndRune(indentLen int, keyword string) int {
	return indentLen + runeLen(keyword)
}

// CompletionsForLine returns replace columns and snippets (Python completions_for_line parity).
// column is a 0-based rune index in line.
func CompletionsForLine(line string, column int) CompletionsResult {
	runes := []rune(line)
	if column < 0 {
		column = len(runes)
	}
	if column > len(runes) {
		column = len(runes)
	}

	empty := CompletionsResult{Start: column, End: column}
	indentLen := leadingIndentRunes(runes)
	stepLine := isStepLineRunes(runes, indentLen)

	if indentLen >= len(runes) || (len(runes) > indentLen && runes[indentLen] == '#') {
		if len(runes) == indentLen && column >= indentLen {
			if stepLine || indentLen == 0 {
				items := append(stepSnippetsForCompletion(), keywordCandidates("")...)
				return CompletionsResult{Start: indentLen, End: column, Items: items}
			}
			return CompletionsResult{Start: indentLen, End: column, Items: keywordCandidates("")}
		}
		return empty
	}

	stripped := string(runes[indentLen:])
	if headerLineRE.MatchString(stripped) && !isStepIndentedRunes(runes, indentLen) {
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
	bodyRunes := []rune(body)
	bodyOffset := indentLen + runeLen(stripped) - len(bodyRunes)

	if keyword != "" {
		keywordEnd := keywordEndRune(indentLen, keyword)
		if column <= keywordEnd {
			prefixRunes := runes[indentLen:column]
			return CompletionsResult{
				Start: indentLen,
				End:   column,
				Items: keywordCandidates(strings.TrimSpace(string(prefixRunes))),
			}
		}
	} else {
		linePrefixRunes := runes[indentLen:column]
		linePrefix := string(linePrefixRunes)
		if !strings.Contains(strings.TrimRight(linePrefix, " "), " ") {
			if matches := keywordCandidates(strings.TrimSpace(linePrefix)); len(matches) > 0 {
				return CompletionsResult{Start: indentLen, End: column, Items: matches}
			}
		}
	}

	bodyPrefix := strings.TrimLeft(string(runes[bodyOffset:column]), " \t")
	if bodyPrefix == "" {
		if keyword != "" {
			return CompletionsResult{Start: bodyOffset, End: column, Items: stepSnippetsForCompletion()}
		}
		return CompletionsResult{Start: indentLen, End: column, Items: keywordCandidates("")}
	}

	matches := stepCandidates(bodyPrefix)
	start := column - runeLen(bodyPrefix)
	return CompletionsResult{Start: start, End: column, Items: matches}
}
