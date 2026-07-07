package update

import "testing"

func TestSkipRemoteReleaseCheckDevInstall(t *testing.T) {
	if !SkipRemoteReleaseCheck(`C:\Users\dev\scenaria_go\build\bin`) {
		t.Fatal("expected skip for build/bin install")
	}
	if SkipRemoteReleaseCheck(`C:\Program Files\Scenaria`) {
		t.Fatal("expected check for Program Files install")
	}
}

func TestCheckInstallDirSkipsRemoteForDevBuild(t *testing.T) {
	info, err := CheckInstallDir("0.28.1", `C:\repo\scenaria_go\build\bin`)
	if err != nil {
		t.Fatal(err)
	}
	if info.UpdateAvailable {
		t.Fatal("dev build should be up to date without GitHub")
	}
	if info.LatestVersion != "0.28.1" {
		t.Fatalf("latest=%q", info.LatestVersion)
	}
}
