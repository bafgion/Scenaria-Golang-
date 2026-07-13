package wailsapp

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gui"
)

func TestBeforeCloseEmitsFrontendDialogInsteadOfNativePrompt(t *testing.T) {
	app := &App{svc: gui.NewService()}
	app.svc.UpdateDirtyTabsState(true)
	prevent := app.BeforeClose(nil)
	if !prevent {
		t.Fatal("expected close to be prevented for frontend dialog")
	}
	if app.closeConfirmed.Load() {
		t.Fatal("expected closeConfirmed to remain false until user confirms")
	}
}

func TestConfirmAppCloseBypassesGuard(t *testing.T) {
	app := &App{svc: gui.NewService()}
	app.svc.UpdateDirtyTabsState(true)
	if !app.BeforeClose(nil) {
		t.Fatal("expected first close attempt to be prevented")
	}
	app.ConfirmAppClose()
	if !app.closeConfirmed.Load() {
		t.Fatal("expected closeConfirmed after ConfirmAppClose")
	}
	if app.BeforeClose(nil) {
		t.Fatal("expected guarded close to be allowed after confirmation")
	}
}

func TestBeforeCloseBlocksEvenWithoutGuardReasons(t *testing.T) {
	app := &App{svc: gui.NewService()}
	if !app.BeforeClose(nil) {
		t.Fatal("expected close to be prevented so the frontend can flush session state")
	}
}
