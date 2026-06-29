package gui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

func parseStepActions(text string) ([]string, error) {
	feature, err := gherkin.ParseFeature(text)
	if err != nil {
		return nil, err
	}
	if len(feature.Scenarios) == 0 {
		return nil, nil
	}
	actions := make([]string, 0, len(feature.Scenarios[0].Steps))
	for _, step := range feature.Scenarios[0].Steps {
		action, err := stepdsl.Parse(step)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action.Kind)
	}
	return actions, nil
}

func TestSaveFeaturePreservesCommentsAndBlankLines(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "annotated.feature")
	tab := "\t"
	editorText := "Функционал: UI\n" +
		"Сценарий: Demo\n" +
		"\n" +
		"# Подготовка пользователя\n" +
		tab + "Допустим открыт \"https://example.com\"\n" +
		"\n" +
		tab + "# TODO: добавить логин\n" +
		tab + "И нажимаю \"button\"\n"

	initial := "Функционал: UI\nСценарий: Demo\n" + tab + "Допустим открыт \"https://example.com\"\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if err := svc.SaveFeature(path, editorText); err != nil {
		t.Fatal(err)
	}

	disk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(disk)
	if !strings.Contains(text, "# Подготовка пользователя") {
		t.Fatalf("missing header comment:\n%s", text)
	}
	if !strings.Contains(text, "# TODO: добавить логин") {
		t.Fatalf("missing step comment:\n%s", text)
	}
	if !strings.HasSuffix(strings.TrimSpace(text), `И нажимаю "button"`) {
		t.Fatalf("unexpected ending:\n%s", text)
	}
	actions, err := parseStepActions(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 || actions[0] != "goto" || actions[1] != "click" {
		t.Fatalf("actions: %v", actions)
	}
}

func TestSaveFeatureWritesRemovedStepToDisk(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "signature_demo.feature")
	tab := "\t"
	initial := "Функционал: signature_demo\n" +
		"Сценарий: signature_demo\n" +
		tab + "Допустим открыт \"https://example.com\"\n" +
		tab + "И нажимаю \"button\"\n" +
		tab + "И закрываю браузер\n"
	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	updated := "Функционал: signature_demo\n" +
		"Сценарий: signature_demo\n" +
		tab + "Допустим открыт \"https://example.com\"\n" +
		tab + "И нажимаю \"button\"\n"

	svc := NewService()
	if err := svc.SaveFeature(path, updated); err != nil {
		t.Fatal(err)
	}

	disk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(disk)
	if strings.Contains(text, "закрываю браузер") {
		t.Fatalf("close step should be removed:\n%s", text)
	}
	actions, err := parseStepActions(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 || actions[0] != "goto" || actions[1] != "click" {
		t.Fatalf("actions: %v", actions)
	}
}

func TestWriteTempFeatureRoundTrip(t *testing.T) {
	tab := "\t"
	content := "Функционал: temp\n" +
		"Сценарий: temp\n" +
		tab + "Допустим открыт \"https://example.com\"\n" +
		tab + "И обновляю страницу\n"

	svc := NewService()
	path, err := svc.WriteTempFeature(content)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(filepath.Dir(path)) })

	got, err := svc.ReadFeature(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != content {
		t.Fatalf("content mismatch:\nwant %q\ngot  %q", content, got)
	}
	actions, err := parseStepActions(got)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 || actions[0] != "goto" || actions[1] != "reload" {
		t.Fatalf("actions: %v", actions)
	}
}

func TestSaveFeaturePreservesCommentsOnRewrite(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "demo.feature")
	tab := "\t"
	if err := os.WriteFile(path, []byte("Функционал: UI\nСценарий: Demo\n"+tab+"Допустим открыт \"https://example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	editorText := "Функционал: UI\n" +
		"Сценарий: Demo\n" +
		"\n" +
		"# обновлено\n" +
		tab + "Допустим открыт \"https://example.com\"\n" +
		tab + "И нажимаю \"go\"\n"

	svc := NewService()
	if err := svc.SaveFeature(path, editorText); err != nil {
		t.Fatal(err)
	}

	disk, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(disk)
	if !strings.Contains(text, "# обновлено") {
		t.Fatalf("missing comment:\n%s", text)
	}
	actions, err := parseStepActions(text)
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 2 || actions[0] != "goto" || actions[1] != "click" {
		t.Fatalf("actions: %v", actions)
	}
}
