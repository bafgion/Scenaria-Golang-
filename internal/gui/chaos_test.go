package gui

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestChaosValidateNeverPanics(t *testing.T) {
	base := "Функционал: Chaos\nСценарий: S\n\tДопустим открыт \"https://example.com\"\n\tКогда нажимаю \"#ok\"\n"
	for seed := int64(0); seed < 60; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			text := corruptFeatureText(base, rng)
			_ = ValidateFeatureContent(text)
			_ = AnalyzeScenarioHints(text)
			_ = FormatFeature(text)
			_ = NormalizeStepIndents(text)
			_ = CollapseBlankLinesBetweenSteps(text)
		})
	}
}

func TestChaosHintFixRoundTripNeverPanics(t *testing.T) {
	features := []string{
		"Функционал: H\nСценарий: S\n\tДопустим открыт \"https://example.com\"\n\tКогда нажимаю \"#menu\"\n\tИ нажимаю \"#item\"\n",
		"Функционал: H\nСценарий: S\n\tДопустим открыт \"https://example.com\"\n\tКогда заполняю \"#email\" значением \"a@b.c\"\n",
	}
	for i, text := range features {
		t.Run(fmt.Sprintf("feature-%d", i), func(t *testing.T) {
			hints := AnalyzeScenarioHints(text)
			current := text
			for _, hint := range hints {
				if !hint.AutoFixable {
					continue
				}
				result := ApplyScenarioHintFix(ScenarioHintFixRequest{
					Text:      current,
					HintID:    hint.ID,
					StepIndex: hint.StepIndex,
				})
				current = result.Text
				_ = ValidateFeatureContent(current)
			}
		})
	}
}

func TestChaosRefactorReplaceStable(t *testing.T) {
	text := "Функционал: R\nСценарий: S\n\tДопустим открыт \"https://old.example\"\n\tКогда нажимаю \"#old\"\n"
	for _, find := range []string{"old", "OLD", "https://old.example", "#old", ""} {
		result := ReplaceInText(text, find, "new", false)
		if result.Text == "" && text != "" {
			t.Fatalf("empty text after replace find=%q", find)
		}
		_ = gherkin.NormalizeFeatureText(result.Text)
	}
}

func corruptFeatureText(base string, rng *rand.Rand) string {
	lines := strings.Split(base, "\n")
	if len(lines) == 0 {
		return base
	}
	switch rng.Intn(6) {
	case 0:
		return strings.Repeat("\t", rng.Intn(5)) + base
	case 1:
		idx := rng.Intn(len(lines))
		lines[idx] = lines[idx] + string(rune(rng.Intn(0x10ffff)))
		return strings.Join(lines, "\n")
	case 2:
		return base + "\n\t" + strings.Repeat("Когда ", rng.Intn(4)) + "нажимаю \"#x\"\n"
	case 3:
		return strings.ReplaceAll(base, "\"", strings.Repeat("\"", rng.Intn(3)))
	case 4:
		return base[:rng.Intn(len(base))]
	default:
		return base
	}
}
