//go:build windows

package wailsapp

import (
	"sync"
	"syscall"
	"unsafe"

	"github.com/bafgion/scenaria-golang/internal/brand"
	"golang.org/x/sys/windows"
)

const (
	gwlStyle   = -16
	gwlExStyle = -20

	wsOverlappedWindow = 0x00CF0000
	wsVisible            = 0x10000000
	wsCaption            = 0x00C00000
	wsThickFrame         = 0x00040000
	wsSysmenu            = 0x00080000
	wsMinimizebox        = 0x00020000
	wsMaximizebox        = 0x00010000

	swpFramechanged  = 0x0020
	swpNomove        = 0x0002
	swpNosize        = 0x0001
	swpNozorder      = 0x0004
	swpNoownerzorder = 0x0200

	swShow = 5
)

var (
	chromeMu       sync.Mutex
	savedStyle     uintptr
	savedExStyle   uintptr
	appWindow      windows.Handle
	splashChromeOn bool
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procFindWindowW              = user32.NewProc("FindWindowW")
	procGetWindowLongPtrW        = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW        = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procShowWindow               = user32.NewProc("ShowWindow")
	procRedrawWindow             = user32.NewProc("RedrawWindow")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procMonitorFromWindow        = user32.NewProc("MonitorFromWindow")
	procGetMonitorInfoW          = user32.NewProc("GetMonitorInfoW")
)

const monitorDefaultToNearest = 2

type winRect struct {
	Left, Top, Right, Bottom int32
}

type monitorInfo struct {
	CbSize    uint32
	RcMonitor winRect
	RcWork    winRect
	DwFlags   uint32
}

func findWindow(className, windowName *uint16) windows.Handle {
	r, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)))
	return windows.Handle(r)
}

func getWindowLong(hwnd windows.Handle, index int32) uintptr {
	r, _, _ := procGetWindowLongPtrW.Call(uintptr(hwnd), uintptr(index))
	return r
}

func setWindowLong(hwnd windows.Handle, index int32, value uintptr) {
	procSetWindowLongPtrW.Call(uintptr(hwnd), uintptr(index), value)
}

func frameChanged(hwnd windows.Handle) {
	procSetWindowPos.Call(
		uintptr(hwnd),
		0,
		0,
		0,
		0,
		0,
		swpNomove|swpNosize|swpNozorder|swpNoownerzorder|swpFramechanged,
	)
}

func redrawWindow(hwnd windows.Handle) {
	const rdwInvalidate = 0x0001
	const rdwAllChildren = 0x0080
	const rdwFrame = 0x0400
	const rdwErase = 0x0004
	procRedrawWindow.Call(uintptr(hwnd), 0, 0, rdwInvalidate|rdwAllChildren|rdwFrame|rdwErase)
	procShowWindow.Call(uintptr(hwnd), swShow)
}

func findAppWindow() windows.Handle {
	if appWindow != 0 {
		return appWindow
	}

	className, err := windows.UTF16PtrFromString("wailsWindow")
	if err == nil {
		if hwnd := findWindow(className, nil); hwnd != 0 {
			return hwnd
		}
	}
	title, err := windows.UTF16PtrFromString(brand.Name)
	if err == nil {
		if hwnd := findWindow(nil, title); hwnd != 0 {
			return hwnd
		}
	}
	return findVisibleProcessWindow()
}

func findVisibleProcessWindow() windows.Handle {
	pid := uint32(windows.GetCurrentProcessId())
	var found windows.Handle
	cb := syscall.NewCallback(func(hwnd uintptr, _ uintptr) uintptr {
		var wpid uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&wpid)))
		if wpid != pid {
			return 1
		}
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}
		found = windows.Handle(hwnd)
		return 0
	})
	procEnumWindows.Call(cb, 0)
	return found
}

func applySplashChrome() {
	hwnd := findAppWindow()
	if hwnd == 0 {
		return
	}

	chromeMu.Lock()
	defer chromeMu.Unlock()

	if !splashChromeOn {
		savedStyle = getWindowLong(hwnd, gwlStyle)
		savedExStyle = getWindowLong(hwnd, gwlExStyle)
		appWindow = hwnd
	}

	undecorated := savedStyle &^ (wsCaption | wsThickFrame | wsSysmenu | wsMinimizebox | wsMaximizebox)
	if undecorated&wsVisible == 0 {
		undecorated |= wsVisible
	}
	setWindowLong(hwnd, gwlStyle, undecorated)
	frameChanged(hwnd)
	splashChromeOn = true
	centerAppWindow()
}

func applyMainChrome() {
	hwnd := appWindow
	if hwnd == 0 {
		hwnd = findAppWindow()
	}
	if hwnd == 0 {
		return
	}

	chromeMu.Lock()
	defer chromeMu.Unlock()

	current := getWindowLong(hwnd, gwlStyle)
	if current&wsCaption != 0 && current&wsSysmenu != 0 {
		splashChromeOn = false
		return
	}

	restoreStyle := savedStyle
	if restoreStyle == 0 || restoreStyle&(wsCaption|wsSysmenu) == 0 {
		restoreStyle = wsOverlappedWindow
	}
	if current&wsVisible != 0 {
		restoreStyle |= wsVisible
	}

	setWindowLong(hwnd, gwlStyle, restoreStyle)
	if savedExStyle != 0 {
		setWindowLong(hwnd, gwlExStyle, savedExStyle)
	}
	frameChanged(hwnd)
	redrawWindow(hwnd)
	splashChromeOn = false
	appWindow = hwnd
}

func centerAppWindow() {
	hwnd := appWindow
	if hwnd == 0 {
		hwnd = findAppWindow()
	}
	if hwnd == 0 {
		return
	}

	var windowRect winRect
	if r, _, _ := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&windowRect))); r == 0 {
		return
	}
	width := windowRect.Right - windowRect.Left
	height := windowRect.Bottom - windowRect.Top

	monitor, _, _ := procMonitorFromWindow.Call(uintptr(hwnd), monitorDefaultToNearest)
	if monitor == 0 {
		return
	}
	var info monitorInfo
	info.CbSize = uint32(unsafe.Sizeof(info))
	if r, _, _ := procGetMonitorInfoW.Call(monitor, uintptr(unsafe.Pointer(&info))); r == 0 {
		return
	}

	work := info.RcWork
	x := work.Left + (work.Right-work.Left-width)/2
	y := work.Top + (work.Bottom-work.Top-height)/2

	procSetWindowPos.Call(
		uintptr(hwnd),
		0,
		uintptr(x),
		uintptr(y),
		0,
		0,
		swpNosize|swpNozorder|swpNoownerzorder,
	)
}
