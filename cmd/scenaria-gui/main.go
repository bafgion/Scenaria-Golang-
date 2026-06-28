//go:build desktop

// Legacy Fyne desktop entrypoint. Not included in portable releases; use Wails (main.go).
package main

import "github.com/bafgion/scenaria-golang/ui/desktop"

func main() {
	desktop.Run()
}
