package player

import (
	"github.com/bafgion/scenaria-golang/internal/logx"
)

func closeBrowserResource(resource string, closeFn func() error) {
	if closeFn == nil {
		return
	}
	if err := closeFn(); err != nil {
		logx.Debug("browser cleanup", "resource", resource, "error", err)
	}
}
