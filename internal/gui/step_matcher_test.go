package gui

import "testing"

func TestMatchStepText_RecoverySmartQuotesAndYo(t *testing.T) {
	match := matchStepText(`ввожу случайный расчётный счёт в “input[name=account]”`, 1)
	if match.ParseErr != nil {
		t.Fatalf("expected recovered parse, got err: %v", match.ParseErr)
	}
	if !match.Recovered {
		t.Fatal("expected recovery mode to be used")
	}
	if match.Action.Kind != "fill-generated" {
		t.Fatalf("unexpected kind: %s", match.Action.Kind)
	}
}

func TestMatchStepText_TestClientSkipsDSL(t *testing.T) {
	match := matchStepText(`я подключаю TestClient "DemoUser"`, 1)
	if !match.TestClient {
		t.Fatal("expected test-client step")
	}
	if match.ParseErr != nil {
		t.Fatalf("unexpected parse error: %v", match.ParseErr)
	}
}
