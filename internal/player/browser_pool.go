package player

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
)

type browserPool struct {
	mu     sync.Mutex
	slots  chan *browserPoolSlot
	stops  []func()
	size   int
	closed bool
}

type browserPoolSlot struct {
	session *browserSession
}

func newBrowserPool(ctx context.Context, options PlaywrightExecutorOptions, size int) (*browserPool, error) {
	if size < 1 {
		size = 1
	}
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, size),
		stops: make([]func(), 0, size),
		size:  size,
	}
	for i := 0; i < size; i++ {
		if err := ctx.Err(); err != nil {
			pool.Close()
			return nil, err
		}
		pw, stopPW, err := startPlaywright(ctx)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("start playwright worker %d: %w", i+1, err)
		}
		session, err := newBrowserSession(pw, options)
		if err != nil {
			stopPW()
			pool.Close()
			return nil, err
		}
		stopWatch := session.watchContext(ctx)
		pool.stops = append(pool.stops, func() {
			stopBrowserWorker(stopWatch, session, stopPW)
		})
		pool.slots <- &browserPoolSlot{session: session}
	}
	return pool, nil
}

func (p *browserPool) acquire(ctx context.Context) (*browserPoolSlot, error) {
	if p == nil {
		return nil, fmt.Errorf("browser pool is nil")
	}
	select {
	case slot, ok := <-p.slots:
		if !ok {
			return nil, fmt.Errorf("browser pool is closed")
		}
		return slot, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *browserPool) release(slot *browserPoolSlot) {
	if p == nil || slot == nil {
		return
	}
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()
	if closed {
		return
	}
	if err := slot.session.resetForScenario(); err != nil {
		logx.Debug("pool reset failed", "error", err)
		slot.session.abortRun()
	}
	select {
	case p.slots <- slot:
	default:
		logx.Debug("pool release skipped", "error", "slot channel is full")
	}
}

// abortActiveSessions marks idle pool workers cancelled so in-flight Playwright work stops promptly.
func (p *browserPool) abortActiveSessions() {
	if p == nil || p.size <= 0 {
		return
	}
	for i := 0; i < p.size; i++ {
		select {
		case slot := <-p.slots:
			if slot != nil && slot.session != nil {
				slot.session.abortRun()
			}
			p.slots <- slot
		default:
			return
		}
	}
}

func (p *browserPool) Close() {
	if p == nil {
		return
	}
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()
	for i := 0; i < len(p.stops); i++ {
		select {
		case <-p.slots:
		case <-time.After(2 * time.Second):
			logx.Debug("browser pool close timed out waiting for idle slot", "index", i)
		}
	}
	for _, stop := range p.stops {
		stop()
	}
	p.stops = nil
	p.size = 0
}

func poolEligible(options PlaywrightExecutorOptions) bool {
	return options.TraceDir == "" && options.VideoDir == ""
}
