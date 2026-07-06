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
	return CompletionsForLineLang(line, column, "ru")
}

// CompletionsForLineLang uses Gherkin dialect ru or en for keywords and step snippets.
func CompletionsForLineLang(line string, column int, lang string) CompletionsResult {
	if lang == "en" {
		return completionsForLineDialect(line, column, completionKeywordsEN, headerLineEN, stepKeywordEN, headerSnippetsEN, stepSnippetsEN)
	}
	return completionsForLineDialect(line, column, completionKeywords, headerLineRE, stepKeywordRE, headerSnippets, stepSnippets)
}

func completionsForLineDialect(
	line string,
	column int,
	keywords []string,
	headerRE *regexp.Regexp,
	stepRE *regexp.Regexp,
	headers []CompletionSnippet,
	snippets []snippetDef,
) CompletionsResult {
	runes := []rune(line)
	if column < 0 {
		column = len(runes)
	}
	if column > len(runes) {
		column = len(runes)
	}

	empty := CompletionsResult{Start: column, End: column}
	indentLen := leadingIndentRunes(runes)
	stepLine := isStepLineRunesWithKeywords(runes, indentLen, keywords)

	if indentLen >= len(runes) || (len(runes) > indentLen && runes[indentLen] == '#') {
		if len(runes) == indentLen && column >= indentLen {
			if stepLine || indentLen == 0 {
				items := append(stepSnippetsForDialect(snippets), keywordCandidatesWith(keywords, "")...)
				return CompletionsResult{Start: indentLen, End: column, Items: items}
			}
			return CompletionsResult{Start: indentLen, End: column, Items: keywordCandidatesWith(keywords, "")}
		}
		return empty
	}

	stripped := string(runes[indentLen:])
	if headerRE.MatchString(stripped) && !isStepIndentedRunes(runes, indentLen) {
		return CompletionsResult{
			Start: indentLen,
			End:   column,
			Items: headerCandidatesWith(headers, stripped),
		}
	}

	if !stepLine {
		return empty
	}

	match := stepRE.FindStringSubmatch(stripped)
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
				Items: keywordCandidatesWith(keywords, strings.TrimSpace(string(prefixRunes))),
			}
		}
	} else {
		linePrefixRunes := runes[indentLen:column]
		linePrefix := string(linePrefixRunes)
		if !strings.Contains(strings.TrimRight(linePrefix, " "), " ") {
			if matches := keywordCandidatesWith(keywords, strings.TrimSpace(linePrefix)); len(matches) > 0 {
				return CompletionsResult{Start: indentLen, End: column, Items: matches}
			}
		}
	}

	bodyPrefix := strings.TrimLeft(string(runes[bodyOffset:column]), " \t")
	if bodyPrefix == "" {
		if keyword != "" {
			return CompletionsResult{Start: bodyOffset, End: column, Items: stepSnippetsForDialect(snippets)}
		}
		return CompletionsResult{Start: indentLen, End: column, Items: keywordCandidatesWith(keywords, "")}
	}

	matches := stepCandidatesWith(snippets, bodyPrefix)
	start := column - runeLen(bodyPrefix)
	return CompletionsResult{Start: start, End: column, Items: matches}
}

func isStepLineRunesWithKeywords(line []rune, indentLen int, keywords []string) bool {
	if isStepIndentedRunes(line, indentLen) {
		return true
	}
	if indentLen >= len(line) {
		return false
	}
	stripped := strings.ToLower(string(line[indentLen:]))
	for _, kw := range keywords {
		lkw := strings.ToLower(kw)
		if strings.HasPrefix(stripped, lkw+" ") || stripped == lkw {
			return true
		}
	}
	return false
}

func keywordCandidatesWith(keywords []string, prefix string) []CompletionSnippet {
	prefix = strings.TrimSpace(prefix)
	out := make([]CompletionSnippet, 0, len(keywords))
	for _, word := range keywords {
		if matchPrefix(word, prefix) {
			out = append(out, CompletionSnippet{
				Label:       word,
				Insert:      word,
				Description: word,
			})
		}
	}
	return out
}

func headerCandidatesWith(headers []CompletionSnippet, prefix string) []CompletionSnippet {
	stripped := strings.TrimLeft(prefix, " \t")
	out := make([]CompletionSnippet, 0)
	for _, snip := range headers {
		if matchPrefix(snip.Label, stripped) || matchPrefix(snip.Insert, stripped) {
			out = append(out, snip)
		}
	}
	return out
}

func stepSnippetsForDialect(snippets []snippetDef) []CompletionSnippet {
	out := make([]CompletionSnippet, len(snippets))
	for i, snip := range snippets {
		out[i] = CompletionSnippet{
			Label:       snip.label,
			Insert:      snip.insert,
			Description: plainDescription(snip.description),
		}
	}
	return out
}

func stepCandidatesWith(snippets []snippetDef, prefix string) []CompletionSnippet {
	out := make([]CompletionSnippet, 0)
	for _, snip := range snippets {
		if matchPrefix(snip.label, prefix) || matchPrefix(snip.insert, prefix) {
			out = append(out, CompletionSnippet{
				Label:       snip.label,
				Insert:      snip.insert,
				Description: plainDescription(snip.description),
			})
		}
	}
	return out
}

