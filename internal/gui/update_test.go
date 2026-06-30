package gui

import (
	"runtime"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/update"
	"github.com/bafgion/scenaria-golang/internal/version"
)

func TestUpdateInfoFromNil(t *testing.T) {
	dto := updateInfoFrom(nil)
	if dto.CurrentVersion != version.Version {
		t.Fatalf("current=%q", dto.CurrentVersion)
	}
	if dto.CanAutoApply != (runtime.GOOS == "windows") {
		t.Fatalf("canAutoApply=%v", dto.CanAutoApply)
	}
}

func TestUpdateInfoFromPopulated(t *testing.T) {
	info := &update.Info{
		CurrentVersion:  "0.16.0",
		LatestVersion:   "0.17.0",
		UpdateAvailable: true,
		HTMLURL:         "https://example.com/release",
		DownloadURL:     "https://example.com/setup.exe",
		DownloadName:    "Scenaria-Setup.exe",
		Message:         "update",
		InstallMode:     "setup",
		ApplyKind:       "setup",
	}
	dto := updateInfoFrom(info)
	if !dto.UpdateAvailable || dto.ApplyKind != "setup" || dto.InstallMode != "setup" {
		t.Fatalf("dto=%+v", dto)
	}
	if !dto.CanAutoApply && runtime.GOOS == "windows" {
		t.Fatal("expected canAutoApply on windows")
	}
}

func TestApplyUpdateNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows runs full apply flow in internal/update")
	}
	svc := &Service{}
	err := svc.ApplyUpdate()
	if err == nil || !strings.Contains(err.Error(), "Windows") {
		t.Fatalf("unexpected error: %v", err)
	}
}
