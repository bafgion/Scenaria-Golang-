package plugin

import "testing"

func TestInstallListUninstall(t *testing.T) {
	tmp := t.TempDir()

	if err := Install(tmp, "vanessa", "https://example.com/vanessa.zip"); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "vanessa" {
		t.Fatalf("unexpected plugins list: %+v", plugins)
	}

	removed, err := Uninstall(tmp, "vanessa")
	if err != nil {
		t.Fatalf("Uninstall returned error: %v", err)
	}
	if !removed {
		t.Fatal("expected plugin to be removed")
	}

	plugins, err = List(tmp)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plugins) != 0 {
		t.Fatalf("expected empty plugin list, got: %+v", plugins)
	}
}
