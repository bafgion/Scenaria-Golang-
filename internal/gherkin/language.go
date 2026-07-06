package gherkin

import (
	"regexp"
	"strings"
)

// Language is the Gherkin dialect for keywords and step text (ru or en).
type Language string

const (
	LangRU Language = "ru"
	LangEN Language = "en"
)

var languageTagRE = regexp.MustCompile(`(?i)^#\s*language\s*:\s*(ru|en)\s*$`)

// ParseLanguageTag scans leading comment lines for `# language: ru|en`.
// Defaults to ru when absent or unknown.
func ParseLanguageTag(content string) Language {
	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			break
		}
		if groups := languageTagRE.FindStringSubmatch(line); groups != nil {
			return NormalizeLanguage(groups[1])
		}
	}
	return LangRU
}

func NormalizeLanguage(raw string) Language {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "en", "english":
		return LangEN
	default:
		return LangRU
	}
}

type headerKind int

const (
	headerNone headerKind = iota
	headerFeature
	headerBackground
	headerScenario
	headerScenarioOutline
	headerExamples
)

type dialect struct {
	feature         []string
	background      []string
	scenario        []string
	scenarioOutline []string
	examples        []string
	steps           []string
	defaultGiven    string
}

var dialects = map[Language]dialect{
	LangRU: {
		feature:         []string{"Функционал", "Функциональность", "Функция", "Feature"},
		background:      []string{"Контекст", "Background"},
		scenario:        []string{"Сценарий", "Scenario"},
		scenarioOutline: []string{"Структура сценария", "Scenario Outline"},
		examples:        []string{"Примеры", "Examples"},
		steps:           []string{"Допустим", "Дано", "Когда", "Тогда", "И", "Но", "Given", "When", "Then", "And", "But", "*"},
		defaultGiven:    "Допустим",
	},
	LangEN: {
		feature:         []string{"Feature", "Функционал", "Функциональность", "Функция"},
		background:      []string{"Background", "Контекст"},
		scenario:        []string{"Scenario", "Сценарий"},
		scenarioOutline: []string{"Scenario Outline", "Структура сценария"},
		examples:        []string{"Examples", "Примеры"},
		steps:           []string{"Given", "When", "Then", "And", "But", "Допустим", "Дано", "Когда", "Тогда", "И", "Но", "*"},
		defaultGiven:    "Given",
	},
}

func dialectFor(lang Language) dialect {
	if d, ok := dialects[lang]; ok {
		return d
	}
	return dialects[LangRU]
}

func stepKeywordsFor(lang Language) []string {
	return dialectFor(lang).steps
}

func matchHeader(line string, lang Language) (kind headerKind, title string, ok bool) {
	d := dialectFor(lang)
	lower := strings.ToLower(line)
	for _, kw := range d.feature {
		prefix := strings.ToLower(kw) + ":"
		if strings.HasPrefix(lower, prefix) {
			return headerFeature, strings.TrimSpace(line[len(kw)+1:]), true
		}
	}
	for _, kw := range d.background {
		prefix := strings.ToLower(kw) + ":"
		if strings.HasPrefix(lower, prefix) {
			return headerBackground, "", true
		}
	}
	for _, kw := range d.scenarioOutline {
		prefix := strings.ToLower(kw) + ":"
		if strings.HasPrefix(lower, prefix) {
			return headerScenarioOutline, strings.TrimSpace(line[len(kw)+1:]), true
		}
	}
	for _, kw := range d.scenario {
		prefix := strings.ToLower(kw) + ":"
		if strings.HasPrefix(lower, prefix) {
			return headerScenario, strings.TrimSpace(line[len(kw)+1:]), true
		}
	}
	for _, kw := range d.examples {
		prefix := strings.ToLower(kw) + ":"
		if strings.HasPrefix(lower, prefix) {
			return headerExamples, "", true
		}
	}
	return headerNone, "", false
}

func defaultGivenKeyword(lang Language) string {
	return dialectFor(lang).defaultGiven
}

func serializeFeatureHeader(lang Language) string {
	if lang == LangEN {
		return "Feature"
	}
	return "Функционал"
}

func serializeBackgroundHeader(lang Language) string {
	if lang == LangEN {
		return "Background"
	}
	return "Контекст"
}

func serializeScenarioHeader(lang Language, outline bool) string {
	if lang == LangEN {
		if outline {
			return "Scenario Outline"
		}
		return "Scenario"
	}
	if outline {
		return "Структура сценария"
	}
	return "Сценарий"
}

func serializeExamplesHeader(lang Language) string {
	if lang == LangEN {
		return "Examples"
	}
	return "Примеры"
}

func languageTagLine(lang Language) string {
	if lang == LangEN {
		return "# language: en\n"
	}
	return "# language: ru\n"
}
