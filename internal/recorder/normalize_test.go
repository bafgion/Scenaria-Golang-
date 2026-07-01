package recorder

import (
	"strings"
	"testing"
)

func TestNormalizeStepsCoalescesFill(t *testing.T) {
	steps := []RecordedStep{
		{Action: "fill", Selector: "#email", Value: "a"},
		{Action: "fill", Selector: "#email", Value: "ab@x.com"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
	if out[0].Value != "ab@x.com" {
		t.Fatalf("unexpected value: %q", out[0].Value)
	}
}

func TestNormalizeStepsCoalescesSelect(t *testing.T) {
	steps := []RecordedStep{
		{Action: "select", Selector: "#lang", Value: "en"},
		{Action: "select", Selector: "#lang", Value: "ru"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
	if out[0].Value != "ru" {
		t.Fatalf("unexpected value: %q", out[0].Value)
	}
}

func TestNormalizeStepsDropsDuplicateGoto(t *testing.T) {
	steps := []RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "goto", Value: "https://example.com"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
}

func TestNormalizeStepsDropsDuplicateClick(t *testing.T) {
	steps := []RecordedStep{
		{Action: "click", Selector: `text="OK"`},
		{Action: "click", Selector: `text="OK"`},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %d", len(out))
	}
}

func TestNormalizeDropsClickBeforeFillAndUpgradesSelector(t *testing.T) {
	steps := []RecordedStep{
		{Action: "click", Selector: "div > label:nth-of-type(1) > div:nth-of-type(2) > input"},
		{
			Action:   "fill",
			Selector: "div > label:nth-of-type(1) > div:nth-of-type(2) > input",
			Value:    "Иван",
			Text:     "Имя",
		},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %+v", out)
	}
	if out[0].Action != "fill" {
		t.Fatalf("expected fill, got %q", out[0].Action)
	}
	if !strings.HasPrefix(out[0].Selector, `label:has-text("Имя`) {
		t.Fatalf("selector: %q", out[0].Selector)
	}
}

func TestNormalizeLegacyCheckboxNoise(t *testing.T) {
	steps := []RecordedStep{
		{Action: "click", Selector: "div > label:nth-of-type(3)"},
		{Action: "click", Selector: "div > label:nth-of-type(3) > input"},
		{
			Action:    "fill",
			Selector:  "div > label:nth-of-type(3) > input",
			Value:     "on",
			InputType: "checkbox",
		},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %+v", out)
	}
	if out[0].Action != "check" {
		t.Fatalf("expected check, got %q", out[0].Action)
	}
	if out[0].Selector != "div > label:nth-of-type(3) > input" {
		t.Fatalf("selector: %q", out[0].Selector)
	}
}

func TestNormalizeUpgradesFragileCheckboxSelectorFromText(t *testing.T) {
	steps := []RecordedStep{
		{
			Action:   "check",
			Selector: "div:nth-of-type(3) > div > div > label:nth-of-type(3) > input",
			Text:     "Подтверждаю согласие на обработку персональных данных",
		},
	}
	out := NormalizeSteps(steps)
	if !strings.HasPrefix(out[0].Selector, `label:has-text("`) {
		t.Fatalf("selector: %q", out[0].Selector)
	}
	if !strings.Contains(strings.ToLower(out[0].Selector), "согласие") {
		t.Fatalf("selector: %q", out[0].Selector)
	}
}

func TestNormalizeFragileCanvasClick(t *testing.T) {
	steps := []RecordedStep{
		{
			Action:   "click",
			Selector: "div:nth-of-type(3) > div > div > div > canvas",
			Text:     "Поставьте подпись",
		},
	}
	out := NormalizeSteps(steps)
	if out[0].Action != "draw-signature" {
		t.Fatalf("expected draw-signature, got %q", out[0].Action)
	}
	if !strings.HasSuffix(out[0].Selector, "canvas") {
		t.Fatalf("selector: %q", out[0].Selector)
	}
	if strings.Contains(out[0].Selector, "nth-of-type") {
		t.Fatalf("selector still fragile: %q", out[0].Selector)
	}
	if out[0].Text != "Поставьте подпись" {
		t.Fatalf("text: %q", out[0].Text)
	}
}

func TestNormalizeUpgradesGenericPlaceholderUsingFieldText(t *testing.T) {
	steps := []RecordedStep{
		{
			Action:   "fill",
			Selector: `input[placeholder="ДД.ММ.ГГГГ"]`,
			Value:    "01.01.1990",
			Text:     "Дата рождения",
		},
		{
			Action:   "fill",
			Selector: `input[placeholder="ДД.ММ.ГГГГ"]`,
			Value:    "02.02.2020",
			Text:     "Дата выдачи",
		},
	}
	out := NormalizeSteps(steps)
	if len(out) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(out))
	}
	if !strings.Contains(out[0].Selector, "рожд") {
		t.Fatalf("birth selector: %q", out[0].Selector)
	}
	if !strings.Contains(out[1].Selector, "выдач") {
		t.Fatalf("issue selector: %q", out[1].Selector)
	}
}

func TestNormalizeUpgradesFragileClickToTextSelector(t *testing.T) {
	steps := []RecordedStep{
		{
			Action:   "click",
			Selector: "div > button:nth-of-type(2)",
			Text:     "Сохранить",
		},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 {
		t.Fatalf("expected 1 step, got %+v", out)
	}
	if out[0].Selector != `text="Сохранить"` {
		t.Fatalf("selector: %q", out[0].Selector)
	}
}

func TestNormalizeKeepsUploadStep(t *testing.T) {
	steps := []RecordedStep{
		{Action: "upload", Selector: "#file", Value: "data.csv"},
	}
	out := NormalizeSteps(steps)
	if len(out) != 1 || out[0].Action != "upload" {
		t.Fatalf("unexpected: %+v", out)
	}
}
