package brand

import (
	"os"
	"strings"
	"testing"
)

func TestBrandingConstants(t *testing.T) {
	if Name != "Scenaria" || Tagline == "" || Description == "" {
		t.Fatalf("unexpected branding constants: name=%q tagline=%q", Name, Tagline)
	}
	if !strings.Contains(AboutText(), Tagline) {
		t.Fatalf("AboutText should include tagline: %q", AboutText())
	}
	if !strings.Contains(strings.ToLower(AboutText()), "автотестов") {
		t.Fatalf("AboutText should include description: %q", AboutText())
	}
}

func TestBrandingAssetsPresent(t *testing.T) {
	dir := Dir()
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Fatalf("branding dir missing: %s (%v)", dir, err)
	}
	for _, name := range RequiredAssets {
		path := AssetPath(name)
		if st, err := os.Stat(path); err != nil || st.IsDir() {
			t.Fatalf("missing branding asset %s (%v)", path, err)
		}
	}
}

func TestUserAgent(t *testing.T) {
	if got := UserAgent("0.15.0"); got != "Scenaria/0.15.0" {
		t.Fatalf("UserAgent = %q", got)
	}
}
