//go:build windows

package winfocus

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	asfwAny        = ^uintptr(0)
	vkMenu         = 0x12
	keyeventfKeyup = 0x0002
	swRestore      = 9
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetClassNameW            = user32.NewProc("GetClassNameW")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procAllowSetForegroundWindow = user32.NewProc("AllowSetForegroundWindow")
	procKeybdEvent               = user32.NewProc("keybd_event")
	procShowWindow               = user32.NewProc("ShowWindow")
	procBringWindowToTop         = user32.NewProc("BringWindowToTop")
)

func raiseNativeWindow(titleHint, urlHint string) error {
	hwnd := findChromiumWindow(titleHint, urlHint)
	if hwnd == 0 {
		return fmt.Errorf("окно браузера не найдено")
	}
	_, _, _ = procAllowSetForegroundWindow.Call(asfwAny)
	_, _, _ = procShowWindow.Call(uintptr(hwnd), swRestore)
	_, _, _ = procBringWindowToTop.Call(uintptr(hwnd))
	_, _, _ = procKeybdEvent.Call(vkMenu, 0, 0, 0)
	_, _, _ = procSetForegroundWindow.Call(uintptr(hwnd))
	_, _, _ = procKeybdEvent.Call(vkMenu, 0, keyeventfKeyup, 0)
	return nil
}

func findChromiumWindow(titleHint, urlHint string) windows.Handle {
	hint := strings.ToLower(strings.TrimSpace(titleHint))
	ownPID := uint32(os.Getpid())
	var best windows.Handle
	var bestScore int
	var fallback []windows.Handle
	cb := syscall.NewCallback(func(hwnd uintptr, _ uintptr) uintptr {
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}
		handle := windows.Handle(hwnd)
		if windowProcessID(handle) == ownPID {
			return 1
		}
		class := windowClass(handle)
		if class != "Chrome_WidgetWin_1" && class != "Chrome_WidgetWin_0" {
			return 1
		}
		if isWailsWindow(handle) {
			return 1
		}
		title := windowTitle(handle)
		score := scoreChromiumTitle(title, hint, urlHint)
		if score > bestScore {
			bestScore = score
			best = handle
		}
		if score >= 2 {
			fallback = append(fallback, handle)
		}
		return 1
	})
	_, _, _ = procEnumWindows.Call(cb, 0)
	if best != 0 {
		return best
	}
	if len(fallback) == 1 {
		return fallback[0]
	}
	return 0
}

func windowProcessID(hwnd windows.Handle) uint32 {
	var pid uint32
	_, _, _ = procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return pid
}

func isWailsWindow(hwnd windows.Handle) bool {
	return windowClass(hwnd) == "wailsWindow"
}

func windowClass(hwnd windows.Handle) string {
	buf := make([]uint16, 256)
	n, _, _ := procGetClassNameW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

func windowTitle(hwnd windows.Handle) string {
	buf := make([]uint16, 512)
	n, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if n == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}
