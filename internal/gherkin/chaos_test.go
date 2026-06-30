package gherkin

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FuzzParseFeature(f *testing.F) {
	seeds := []string{
		"",
		"Функционал: X\nСценарий: Y\n\tКогда нажимаю \"#ok\"\n",
		"Функционал:\n",
		"Контекст:\n\tДано я подключаю TestClient \"U\"\nСценарий: S\n\tДопустим открыт \"https://x\"\n",
		string([]byte{0, 255, 127}),
	}
	examplesDir := filepath.Join("..", "..", "examples")
	if entries, err := os.ReadDir(examplesDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".feature" {
				continue
			}
			payload, readErr := os.ReadFile(filepath.Join(examplesDir, entry.Name()))
			if readErr == nil {
				seeds = append(seeds, string(payload))
			}
		}
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, content string) {
		t.Parallel()
		_ = NormalizeFeatureText(content)
		_ = CoalesceMixedIndents(content)
		feature, err := ParseFeature(content)
		if err != nil || feature == nil {
			return
		}
		_ = ValidateFeature(feature)
		for _, runnable := range ExpandFeature(feature) {
			for _, step := range FlattenSteps(runnable.Steps) {
				if step.Block != "" || IsTestClientStep(step) {
					continue
				}
				if strings.TrimSpace(step.Text) == "" {
					continue
				}
				_ = step.Text
			}
		}
	})
}

func TestChaosOutlineExpansionNeverPanics(t *testing.T) {
	const iterations = 40
	for seed := int64(0); seed < iterations; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			rows := 1 + rng.Intn(12)
			var b strings.Builder
			b.WriteString("Функционал: Chaos\nСтруктура сценария: Outline\n")
			b.WriteString("\tДопустим открыт \"<url>\"\n\tТогда вижу \"<tag>\"\nПримеры:\n\t| url | tag |\n")
			for i := 0; i < rows; i++ {
				fmt.Fprintf(&b, "\t| https://example.com/%d | h%d |\n", i, i)
			}
			feature, err := ParseFeature(b.String())
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			expanded := ExpandFeature(feature)
			if len(expanded) != rows {
				t.Fatalf("expected %d expanded scenarios, got %d", rows, len(expanded))
			}
			for _, runnable := range expanded {
				for _, step := range FlattenSteps(runnable.Steps) {
					if strings.Contains(step.Text, "<") {
						t.Fatalf("unexpanded placeholder in %q", step.Text)
					}
				}
			}
		})
	}
}

func TestChaosBlockIndentCombinations(t *testing.T) {
	headers := []string{
		`Если вижу "#x"`,
		`Повторяю 2 раза`,
		`Пока вижу "#y"`,
	}
	for i, header := range headers {
		t.Run(header, func(t *testing.T) {
			content := fmt.Sprintf("Функционал: B\nСценарий: S\n\tДопустим открыт \"https://example.com\"\n\t%s\n\t\tИ жду 1 сек\n\tТогда вижу \"h1\"\n", header)
			feature, err := ParseFeature(content)
			if err != nil {
				t.Fatalf("parse block %d: %v", i, err)
			}
			if len(feature.Scenarios) != 1 || len(feature.Scenarios[0].Steps) < 3 {
				t.Fatalf("unexpected tree: %+v", feature.Scenarios[0].Steps)
			}
			issues := ValidateFeature(feature)
			if len(issues) > 0 {
				t.Fatalf("validation: %+v", issues)
			}
		})
	}
}
