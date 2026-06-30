package update

import (
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

func TestPickAssetForModeSetup(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: brand.PortableZip, BrowserDownloadURL: "https://example.com/portable.zip"},
		{Name: brand.SetupExe, BrowserDownloadURL: "https://example.com/setup.exe"},
	}
	got := pickAssetForMode(InstallModeSetup, assets)
	if got == nil || got.Name != brand.SetupExe {
		t.Fatalf("unexpected asset: %+v", got)
	}
	if applyKindForAsset(got.Name) != ApplyKindSetup {
		t.Fatalf("expected setup kind")
	}
}

func TestPickAssetForModePortable(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: brand.PortableZip, BrowserDownloadURL: "https://example.com/portable.zip"},
		{Name: brand.SetupExe, BrowserDownloadURL: "https://example.com/setup.exe"},
	}
	got := pickAssetForMode(InstallModePortable, assets)
	if got == nil || got.Name != brand.PortableZip {
		t.Fatalf("unexpected asset: %+v", got)
	}
	if applyKindForAsset(got.Name) != ApplyKindPortable {
		t.Fatalf("expected portable kind")
	}
}

func TestPortableUpdateScriptContainsRobocopy(t *testing.T) {
	script := portableUpdateScript(`C:\Scenaria`, `C:\Scenaria\_update_staging`, 12345)
	if !strings.Contains(script, "robocopy") {
		t.Fatalf("expected robocopy in script: %s", script)
	}
	if !strings.Contains(script, "scenaria-gui.exe") {
		t.Fatalf("expected gui restart in script")
	}
}
