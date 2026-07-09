package gui

import (
	"fmt"
	"strings"
)

type PickerStepChoice struct {
	Label       string `json:"label"`
	StepBody    string `json:"stepBody"`
	Description string `json:"description"`
	Preview     string `json:"preview"`
}

type SelectorCandidate struct {
	Selector      string   `json:"selector"`
	Strategy      string   `json:"strategy"`
	Score         int      `json:"score"`
	MatchesCount  int      `json:"matches_count"`
	Unique        bool     `json:"unique"`
	Visible       bool     `json:"visible"`
	MatchesPicked bool     `json:"matches_picked"`
	Warnings      []string `json:"warnings,omitempty"`
}

type PickSelectorResult struct {
	Selector        string              `json:"selector"`
	Error           string              `json:"error"`
	SuggestedAction string              `json:"suggested_action,omitempty"`
	SuggestedChoice int                 `json:"suggested_choice,omitempty"`
	Warnings        []string            `json:"warnings,omitempty"`
	Candidates      []SelectorCandidate `json:"candidates,omitempty"`
}

func PickerStepChoices(selector, keyword string) []PickerStepChoice {
	selector = strings.TrimSpace(selector)
	if keyword == "" {
		keyword = "Допустим"
	}
	quoted := quotePickerSelector(selector)
	templates := []struct {
		label        string
		body         string
		description  string
		selectorOnly bool
	}{
		{"Клик", fmt.Sprintf(`нажимаю %s`, quoted), "Клик по элементу", false},
		{"Заполнить", fmt.Sprintf(`ввожу "" в %s`, quoted), "Ввод текста в поле", false},
		{"Выбрать", fmt.Sprintf(`выбираю "" в %s`, quoted), "Выбор значения в списке", false},
		{"Двойной клик", fmt.Sprintf(`дважды нажимаю %s`, quoted), "Двойной клик", false},
		{"Наведение", fmt.Sprintf(`навожу %s`, quoted), "Наведение курсора", false},
		{"Видимость", fmt.Sprintf(`вижу %s`, quoted), "Элемент виден", false},
		{"Скрыт", fmt.Sprintf(`не вижу %s`, quoted), "Элемент скрыт", false},
		{"Очистка поля", fmt.Sprintf(`очищаю %s`, quoted), "Очистить ввод", false},
		{"Галочка", fmt.Sprintf(`отмечаю %s`, quoted), "Установить галочку", false},
		{"Снять галочку", fmt.Sprintf(`снимаю отметку с %s`, quoted), "Снять галочку", false},
		{"Скролл", fmt.Sprintf(`скроллю к %s`, quoted), "Прокрутка к элементу", false},
		{"Жду появления", fmt.Sprintf(`жду появления %s`, quoted), "Ожидание элемента", false},
		{"Жду исчезновения", fmt.Sprintf(`жду исчезновения %s`, quoted), "Ожидание скрытия", false},
		{"Только селектор", quoted, "Вставить селектор без шага Gherkin", true},
	}
	out := make([]PickerStepChoice, 0, len(templates))
	for _, item := range templates {
		preview := quoted
		if !item.selectorOnly {
			preview = formatPickerStep(keyword, item.body)
		}
		out = append(out, PickerStepChoice{
			Label:       item.label,
			StepBody:    item.body,
			Description: item.description,
			Preview:     preview,
		})
	}
	return out
}

// PickerSuggestedChoiceIndex maps picker suggested_action to PickerStepChoices label index.
func PickerSuggestedChoiceIndex(suggestedAction string) int {
	labelByAction := map[string]string{
		"click":   "Клик",
		"fill":    "Заполнить",
		"select":  "Выбрать",
		"check":   "Галочка",
		"uncheck": "Снять галочку",
		"hover":   "Наведение",
	}
	label, ok := labelByAction[strings.TrimSpace(suggestedAction)]
	if !ok {
		return 0
	}
	choices := PickerStepChoices("x", "Допустим")
	for i, choice := range choices {
		if choice.Label == label {
			return i
		}
	}
	return 0
}

func (s *Service) PickSelector() PickSelectorResult {
	snap := s.recorderOps().Session().Snapshot()
	return s.recorderOps().PickSelector(snap.Session, snap.RecordCtx)
}

func quotePickerSelector(value string) string {
	return `"` + escapePickerArg(strings.TrimSpace(value)) + `"`
}

func escapePickerArg(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}

func formatPickerStep(keyword, body string) string {
	return "  " + keyword + " " + body
}
