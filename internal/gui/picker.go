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

type PickSelectorResult struct {
	Selector string `json:"selector"`
	Error    string `json:"error"`
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

func (s *Service) PickSelector() PickSelectorResult {
	s.mu.RLock()
	session := s.liveSession
	ctx := s.recordCtx
	s.mu.RUnlock()
	return s.recorderOps().PickSelector(session, ctx)
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
