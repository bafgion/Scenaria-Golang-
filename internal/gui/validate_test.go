package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFeatureContent(t *testing.T) {
	issues := ValidateFeatureContent("Функционал: X\nСценарий: Y\n  Когда нажимаю \"#ok\"\n  Когда битый шаг\n")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %#v", len(issues), issues)
	}
	if issues[0].Line != 4 {
		t.Fatalf("unexpected line: %d", issues[0].Line)
	}
}

func TestValidateFeatureContent_TagsAndOutline(t *testing.T) {
	root := filepath.Join("..", "..", "examples")
	for _, name := range []string{"01-pervaya-proverka.feature", "04-tablica-primerov.feature", "05-testclient-kontekst.feature"} {
		path := filepath.Join(root, name)
		payload, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		issues := ValidateFeatureContent(string(payload))
		if len(issues) != 0 {
			t.Fatalf("%s: expected no issues, got %d: %#v", name, len(issues), issues)
		}
	}
}

func TestValidateFeatureContent_TagLineOnly(t *testing.T) {
	content := `@smoke
Функционал: X
Сценарий: Y
  Допустим открыт "https://example.com"
  И закрываю браузер
`
	if issues := ValidateFeatureContent(content); len(issues) != 0 {
		t.Fatalf("expected no issues for tagged feature, got %#v", issues)
	}
}

func TestValidateFeatureContent_TestClientContext(t *testing.T) {
	content := `Функционал: Примеры
Контекст:
	Дано я подключаю TestClient "DemoUser"
Сценарий: S
	Допустим открыт "https://example.com"
	И закрываю браузер
`
	if issues := ValidateFeatureContent(content); len(issues) != 0 {
		t.Fatalf("expected no issues for TestClient context, got %#v", issues)
	}
}
