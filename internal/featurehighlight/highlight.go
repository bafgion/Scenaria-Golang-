package featurehighlight

import (
	"regexp"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

type Kind int

const (
	KindDefault Kind = iota
	KindComment
	KindTag
	KindGherkinKeyword
	KindStepKeyword
	KindBlockKeyword
	KindString
	KindTestClient
	KindError
)

type Span struct {
	Text string
	Kind Kind
}

var (
	gherkinLineRe = regexp.MustCompile(`(?i)^(функциональность|функционал|функция|feature|сценарий|scenario|структура сценария|scenario outline|примеры|examples|контекст|background)\s*:`)
	blockLineRe   = regexp.MustCompile(`(?i)^(если|повторяю|пока|для каждого|иначе|конец если|конец)(?:\s|$)`)
	stepKeywordRe = regexp.MustCompile(`(?i)^\s*(Допустим|Дано|Когда|Тогда|И|Но|Given|When|Then|And|But)\s+`)
	tagRe         = regexp.MustCompile(`@\w+`)
	stringRe      = regexp.MustCompile(`"([^"\\]|\\.)*"`)
)

func Highlight(text string) []Span {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	out := make([]Span, 0, len(text)/8+len(lines))
	for i, line := range lines {
		if i > 0 {
			out = append(out, Span{Text: "\n", Kind: KindDefault})
		}
		out = append(out, highlightLine(line)...)
	}
	return out
}

func highlightLine(line string) []Span {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return []Span{{Text: line, Kind: KindDefault}}
	}
	if strings.HasPrefix(strings.TrimLeft(line, " \t"), "#") {
		return []Span{{Text: line, Kind: KindComment}}
	}
	if strings.HasPrefix(trimmed, "@") {
		return highlightTags(line)
	}
	if gherkinLineRe.MatchString(trimmed) {
		return highlightGherkinHeader(line)
	}
	if blockLineRe.MatchString(trimmed) {
		return highlightBlock(line)
	}
	if loc := stepKeywordRe.FindStringSubmatchIndex(line); loc != nil {
		prefix := line[:loc[1]]
		rest := line[loc[1]:]
		stepText := strings.TrimSpace(rest)
		kind := KindDefault
		if stepText != "" {
			if _, err := stepdsl.Parse(gherkin.Step{Line: 1, Text: stepText}); err != nil {
				kind = KindError
			}
		}
		spans := []Span{{Text: prefix, Kind: KindStepKeyword}}
		spans = append(spans, highlightStrings(rest, kind)...)
		return spans
	}
	if strings.Contains(line, "TestClient") {
		return highlightTestClient(line)
	}
	return highlightStrings(line, KindDefault)
}

func highlightTestClient(line string) []Span {
	idx := strings.Index(line, "TestClient")
	if idx < 0 {
		return highlightStrings(line, KindDefault)
	}
	out := []Span{{Text: line[:idx], Kind: KindDefault}}
	out = append(out, Span{Text: "TestClient", Kind: KindTestClient})
	out = append(out, highlightStrings(line[idx+len("TestClient"):], KindDefault)...)
	return out
}

func highlightGherkinHeader(line string) []Span {
	idx := strings.Index(line, ":")
	if idx < 0 {
		return []Span{{Text: line, Kind: KindGherkinKeyword}}
	}
	spans := []Span{
		{Text: line[:idx+1], Kind: KindGherkinKeyword},
	}
	spans = append(spans, highlightStrings(line[idx+1:], KindDefault)...)
	return spans
}

func highlightBlock(line string) []Span {
	loc := blockLineRe.FindStringIndex(strings.TrimSpace(line))
	if loc == nil {
		return highlightStrings(line, KindDefault)
	}
	lead := len(line) - len(strings.TrimLeft(line, " \t"))
	trim := strings.TrimLeft(line, " \t")
	word := blockLineRe.FindString(trim)
	if word == "" {
		return highlightStrings(line, KindDefault)
	}
	prefixLen := lead + strings.Index(trim, word) + len(word)
	spans := []Span{
		{Text: line[:prefixLen], Kind: KindBlockKeyword},
	}
	spans = append(spans, highlightStrings(line[prefixLen:], KindDefault)...)
	return spans
}

func highlightTags(line string) []Span {
	out := make([]Span, 0, 4)
	cursor := 0
	for _, loc := range tagRe.FindAllStringIndex(line, -1) {
		if loc[0] > cursor {
			out = append(out, Span{Text: line[cursor:loc[0]], Kind: KindDefault})
		}
		out = append(out, Span{Text: line[loc[0]:loc[1]], Kind: KindTag})
		cursor = loc[1]
	}
	if cursor < len(line) {
		out = append(out, Span{Text: line[cursor:], Kind: KindDefault})
	}
	return out
}

func highlightStrings(line string, defaultKind Kind, prefix ...Span) []Span {
	out := append([]Span{}, prefix...)
	cursor := 0
	matches := stringRe.FindAllStringIndex(line, -1)
	if len(matches) == 0 {
		if len(out) == 0 {
			return []Span{{Text: line, Kind: defaultKind}}
		}
		out = append(out, Span{Text: line, Kind: defaultKind})
		return out
	}
	for _, loc := range matches {
		if loc[0] > cursor {
			out = append(out, Span{Text: line[cursor:loc[0]], Kind: defaultKind})
		}
		out = append(out, Span{Text: line[loc[0]:loc[1]], Kind: KindString})
		cursor = loc[1]
	}
	if cursor < len(line) {
		out = append(out, Span{Text: line[cursor:], Kind: defaultKind})
	}
	return out
}
