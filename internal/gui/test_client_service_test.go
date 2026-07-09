package gui

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

func TestTestClientServiceRequiresProject(t *testing.T) {
	svc := NewTestClientService(func() string { return "" })
	if _, err := svc.List(); err == nil {
		t.Fatal("expected error without project")
	}
}

func TestTestClientServiceListRoundTrip(t *testing.T) {
	root := t.TempDir()
	if err := settings.SaveTestClientFromJSON(root, "demo", `{"name":"demo","base_url":"https://example.com"}`); err != nil {
		t.Fatalf("seed test client: %v", err)
	}
	svc := NewTestClientService(func() string { return root })
	names, err := svc.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(names) != 1 || names[0] != "demo" {
		t.Fatalf("unexpected names: %#v", names)
	}
	details, err := svc.Details("demo")
	if err != nil {
		t.Fatalf("details: %v", err)
	}
	if details == "" {
		t.Fatal("expected details summary")
	}
	raw, err := svc.ReadJSON("demo")
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	if raw == "" {
		t.Fatal("expected json payload")
	}
	if err := svc.Delete("demo"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	names, err = svc.List()
	if err != nil {
		t.Fatalf("list after delete: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected empty list, got %#v", names)
	}
}
