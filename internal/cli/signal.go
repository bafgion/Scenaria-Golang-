package cli

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// InterruptContext returns a context cancelled on SIGINT or SIGTERM (CLI runs).
func InterruptContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	if runtime.GOOS == "windows" {
		signal.Notify(sigCh, os.Interrupt)
	} else {
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	}
	go func() {
		select {
		case <-sigCh:
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
	}()
	return ctx, cancel
}
