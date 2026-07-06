package playwrightrt

import (
	"context"
	"sync"
	"testing"

	playwright "github.com/mxschmitt/playwright-go"
)

func TestAcquireReleaseSingleDriver(t *testing.T) {
	Shutdown()
	t.Cleanup(Shutdown)

	pw1, release1, err := Acquire(context.Background())
	if err != nil {
		t.Skip("playwright not available:", err)
	}
	if pw1 == nil {
		t.Fatal("nil playwright")
	}

	pw2, release2, err := Acquire(context.Background())
	if err != nil {
		release1()
		t.Fatal(err)
	}
	if pw1 != pw2 {
		t.Fatal("expected same Playwright instance")
	}

	release2()
	release1()
}

func TestAcquireConcurrentStartsOnce(t *testing.T) {
	Shutdown()
	t.Cleanup(Shutdown)

	var wg sync.WaitGroup
	handles := make(chan acquiredHandle, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pw, release, err := Acquire(context.Background())
			if err != nil {
				return
			}
			handles <- acquiredHandle{pw: pw, release: release}
		}()
	}
	wg.Wait()
	close(handles)

	var first *playwright.Playwright
	count := 0
	for item := range handles {
		count++
		if first == nil {
			first = item.pw
		} else if first != item.pw {
			t.Fatal("concurrent acquire returned different drivers")
		}
		item.release()
	}
	if count == 0 {
		t.Skip("playwright not available")
	}
}

type acquiredHandle struct {
	pw      *playwright.Playwright
	release func()
}

func TestAcquireRespectsCancel(t *testing.T) {
	Shutdown()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err := Acquire(ctx)
	if err == nil {
		t.Fatal("expected cancel error")
	}
}
