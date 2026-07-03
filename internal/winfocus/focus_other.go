//go:build !windows

package winfocus

import "fmt"

func raiseNativeWindow(titleHint, urlHint string) error {
	_ = titleHint
	_ = urlHint
	return fmt.Errorf("показ окна браузера поддерживается только в Windows")
}
