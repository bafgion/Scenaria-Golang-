package settings

import "testing"

func TestCheckUpdatesOnStartupEnabled(t *testing.T) {
	if !CheckUpdatesOnStartupEnabled(nil) {
		t.Fatal("nil cfg should default to true")
	}
	if !CheckUpdatesOnStartupEnabled(&AppSettings{}) {
		t.Fatal("missing field should default to true")
	}
	off := false
	if CheckUpdatesOnStartupEnabled(&AppSettings{CheckUpdatesOnStartup: &off}) {
		t.Fatal("explicit false should disable")
	}
	on := true
	if !CheckUpdatesOnStartupEnabled(&AppSettings{CheckUpdatesOnStartup: &on}) {
		t.Fatal("explicit true should enable")
	}
}
