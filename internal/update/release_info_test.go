package update

import (
	"strings"
	"testing"
)

func TestPickUpdateAsset(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "source.zip", BrowserDownloadURL: "https://example.com/source.zip"},
		{Name: "Scenaria-Portable.zip", BrowserDownloadURL: "https://example.com/portable.zip"},
	}
	got := pickUpdateAsset(assets)
	if got == nil || got.Name != "Scenaria-Portable.zip" {
		t.Fatalf("unexpected asset: %+v", got)
	}
}

func TestPickUpdateAssetPrefersInstallerOnWindows(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "Scenaria-Portable.zip", BrowserDownloadURL: "https://example.com/portable.zip"},
		{Name: "Scenaria-Setup.exe", BrowserDownloadURL: "https://example.com/setup.exe"},
		{Name: "Scenaria-Setup.msi", BrowserDownloadURL: "https://example.com/setup.msi"},
	}
	var got *ReleaseAsset
	for _, pref := range updateAssetPreferences("windows") {
		if asset := matchAsset(assets, pref); asset != nil {
			got = asset
			break
		}
	}
	if got == nil || got.Name != "Scenaria-Setup.exe" {
		t.Fatalf("unexpected asset: %+v", got)
	}
}

func TestPickUpdateAssetPrefersMSIWhenNoSetupExe(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "Scenaria-Portable.zip", BrowserDownloadURL: "https://example.com/portable.zip"},
		{Name: "Scenaria-Setup.msi", BrowserDownloadURL: "https://example.com/setup.msi"},
	}
	got := matchAsset(assets, func(name string) bool { return strings.HasSuffix(name, ".msi") })
	if got == nil || got.Name != "Scenaria-Setup.msi" {
		t.Fatalf("unexpected asset: %+v", got)
	}
}

func TestCheckMessageWhenUpToDate(t *testing.T) {
	info := &Info{CurrentVersion: "v1.0.0", LatestVersion: "v1.0.0"}
	info.UpdateAvailable = IsNewer(info.CurrentVersion, info.LatestVersion)
	if info.UpdateAvailable {
		t.Fatal("expected no update")
	}
}
