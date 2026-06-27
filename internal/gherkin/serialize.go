package gherkin

import (
	"fmt"
	"os"
	"strings"
)

func SerializeFeature(feature *Feature) (string, error) {
	if feature == nil {
		return "", fmt.Errorf("feature is nil")
	}
	if strings.TrimSpace(feature.Title) == "" {
		return "", fmt.Errorf("feature title is empty")
	}

	var b strings.Builder

	writeTags(&b, feature.Tags)
	b.WriteString("Функционал: ")
	b.WriteString(feature.Title)
	b.WriteString("\n")

	if len(feature.Background) > 0 {
		b.WriteString("\n")
		b.WriteString("Контекст:\n")
		for _, step := range feature.Background {
			writeStep(&b, step)
		}
	}

	for i, scenario := range feature.Scenarios {
		if i == 0 && len(feature.Background) == 0 {
			b.WriteString("\n")
		} else {
			b.WriteString("\n")
		}

		writeTags(&b, scenario.Tags)

		if scenario.IsOutline {
			b.WriteString("Структура сценария: ")
		} else {
			b.WriteString("Сценарий: ")
		}
		b.WriteString(scenario.Title)
		b.WriteString("\n")

		for _, step := range scenario.Steps {
			writeStep(&b, step)
		}

		if scenario.IsOutline {
			for _, example := range scenario.Examples {
				b.WriteString("\n")
				b.WriteString("Примеры:\n")
				for _, row := range example.Rows {
					writeTableRow(&b, row, "  ")
				}
			}
		}
	}

	return strings.TrimRight(b.String(), "\n") + "\n", nil
}

func SaveFeatureFile(path string, feature *Feature) error {
	content, err := SerializeFeature(feature)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write feature file: %w", err)
	}
	return nil
}

func writeTags(b *strings.Builder, tags []string) {
	if len(tags) == 0 {
		return
	}
	b.WriteString(strings.Join(tags, " "))
	b.WriteString("\n")
}

func writeStep(b *strings.Builder, step Step) {
	b.WriteString("  ")
	b.WriteString(step.Keyword)
	if step.Text != "" {
		b.WriteString(" ")
		b.WriteString(step.Text)
	}
	b.WriteString("\n")

	if step.DocString != "" {
		b.WriteString("  \"\"\"\n")
		for _, line := range strings.Split(step.DocString, "\n") {
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("  \"\"\"\n")
	}
	for _, row := range step.Table {
		writeTableRow(b, row, "  ")
	}
}

func writeTableRow(b *strings.Builder, row []string, indent string) {
	b.WriteString(indent)
	b.WriteString("|")
	for _, cell := range row {
		b.WriteString(" ")
		b.WriteString(cell)
		b.WriteString(" |")
	}
	b.WriteString("\n")
}
