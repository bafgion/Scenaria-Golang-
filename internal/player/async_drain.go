package player

import (
	"sync"
	"time"
)

var pendingAsync sync.WaitGroup

func drainAsync(ch <-chan error) {
	pendingAsync.Add(1)
	go func() {
		defer pendingAsync.Done()
		<-ch
	}()
}

func drainChan[T any](ch <-chan T) {
	pendingAsync.Add(1)
	go func() {
		defer pendingAsync.Done()
		<-ch
	}()
}

func drainPendingAsync(maxWait time.Duration) {
	if maxWait <= 0 {
		pendingAsync.Wait()
		return
	}
	done := make(chan struct{})
	go func() {
		pendingAsync.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(maxWait):
	}
}
