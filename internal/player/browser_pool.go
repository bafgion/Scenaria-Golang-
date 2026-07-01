package player

import (
	"context"
	"fmt"
)

type browserPool struct {
	slots chan *browserPoolSlot
	stops []func()
	size  int
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
			stopWatch()
			session.close()
			stopPW()
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
	if err := slot.session.resetForScenario(); err != nil {
		return
	}
	p.slots <- slot
}

func (p *browserPool) Close() {
	if p == nil {
		return
	}
	for i := 0; i < len(p.stops); i++ {
		<-p.slots
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
