//go:build !windows

package winfocus

import "fmt"

func raiseNativeWindow(titleHint, urlHint string, processID uint32) error {
	_ = titleHint
	_ = urlHint
	_ = processID
	return fmt.Errorf("показ окна браузера поддерживается только в Windows")
}
