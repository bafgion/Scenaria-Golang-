package update

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

func portableUpdateScript(installDir, stagingDir string, parentPID int) string {
	target := batchQuote(filepath.Clean(installDir))
	staging := batchQuote(filepath.Clean(stagingDir))
	gui := batchQuote(filepath.Join(filepath.Clean(installDir), brand.GUIExeName))
	pid := strconv.Itoa(parentPID)
	lines := []string{
		"@echo off",
		"setlocal EnableExtensions EnableDelayedExpansion",
		"ping -n 3 127.0.0.1 >nul",
		"set WAIT_SEC=0",
		"set MAX_WAIT=120",
		":wait_exit",
		`for /f %%C in ('tasklist /FI "PID eq ` + pid + `" /NH 2^>nul ^| find /C "` + pid + `"') do set "FOUND=%%C"`,
		`if not defined FOUND set "FOUND=0"`,
		`if "!FOUND!"=="0" goto do_copy`,
		"set /a WAIT_SEC+=1",
		"if !WAIT_SEC! GEQ !MAX_WAIT! goto force_kill",
		"ping -n 2 127.0.0.1 >nul",
		"goto wait_exit",
		":force_kill",
		"taskkill /PID " + pid + " /T /F >nul 2>&1",
		"taskkill /IM " + brand.GUIExeName + " /T /F >nul 2>&1",
		"taskkill /IM scenaria.exe /T /F >nul 2>&1",
		"ping -n 3 127.0.0.1 >nul",
		":do_copy",
		`robocopy ` + staging + ` ` + target + ` /E /XD browsers .scenaria ` + updateStagingDir + ` /R:3 /W:2 /NFL /NDL /NJH /NJS /NC /NS /NP`,
		"set RC=%ERRORLEVEL%",
		"if %RC% GEQ 8 exit /b %RC%",
		`if exist ` + staging + ` rmdir /S /Q ` + staging,
		`cd /d ` + target,
		`start "" /D ` + target + ` ` + gui,
		"endlocal",
		"del \"%~f0\" >nul 2>&1",
		"exit /b 0",
	}
	return strings.Join(lines, "\r\n") + "\r\n"
}

func batchQuote(path string) string {
	return `"` + strings.ReplaceAll(path, `"`, `""`) + `"`
}
