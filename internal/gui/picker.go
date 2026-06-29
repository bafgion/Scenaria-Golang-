package gui

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/recorder"
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
		label       string
		body        string
		description string
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
	s.mu.RUnlock()
	if session == nil {
		return PickSelectorResult{Error: "браузер не открыт"}
	}
	selector, err := session.PickSelector(context.Background())
	if err != nil {
		if errors.Is(err, recorder.ErrPickerCancelled) {
			return PickSelectorResult{Error: "отменено"}
		}
		return PickSelectorResult{Error: err.Error()}
	}
	return PickSelectorResult{Selector: selector}
}

func quotePickerSelector(value string) string {
	value = strings.TrimSpace(value)
	if strings.Contains(value, `"`) {
		return "'" + value + "'"
	}
	return `"` + value + `"`
}

func formatPickerStep(keyword, body string) string {
	return "  " + keyword + " " + body
}
