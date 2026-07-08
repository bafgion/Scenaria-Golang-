package gui

import (
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
)

const wailsTimingThreshold = 50 * time.Millisecond

func logWailsTiming(op string, started time.Time) {
	if elapsed := time.Since(started); elapsed >= wailsTimingThreshold {
		logx.Debug("wails timing", "op", op, "elapsed_ms", elapsed.Milliseconds())
	}
}
