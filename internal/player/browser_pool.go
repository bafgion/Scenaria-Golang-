package player

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
)

type browserPool struct {
	mu      sync.Mutex
	slots   chan *browserPoolSlot
	size    int
	retired int
	closed  bool
}

type browserPoolSlot struct {
	session *browserSession
	stop    func()
	index   int
}

func newPoolSlot(ctx context.Context, options PlaywrightExecutorOptions, index int) (*browserPoolSlot, error) {
	pw, stopPW, err := startPlaywright(ctx)
	if err != nil {
		return nil, fmt.Errorf("start playwright worker %d: %w", index+1, err)
	}
	session, err := newBrowserSession(pw, options)
	if err != nil {
		stopPW()
		return nil, err
	}
	stopWatch := session.watchContext(ctx)
	stop := func() {
		stopBrowserWorker(stopWatch, session, stopPW)
	}
	return &browserPoolSlot{session: session, stop: stop, index: index}, nil
}

func newBrowserPool(ctx context.Context, options PlaywrightExecutorOptions, size int) (*browserPool, error) {
	if size < 1 {
		size = 1
	}
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, size),
		size:  size,
	}
	for i := 0; i < size; i++ {
		if err := ctx.Err(); err != nil {
			pool.Close()
			return nil, err
		}
		slot, err := newPoolSlot(ctx, options, i)
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.slots <- slot
	}
	return pool, nil
}

func (p *browserPool) acquire(ctx context.Context) (*browserPoolSlot, error) {
	if p == nil {
		return nil, fmt.Errorf("browser pool is nil")
	}
	p.mu.Lock()
	closed := p.closed
	exhausted := p.retired >= p.size
	p.mu.Unlock()
	if closed {
		return nil, fmt.Errorf("browser pool is closed")
	}
	if exhausted {
		return nil, fmt.Errorf("browser pool has no healthy slots")
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

func (p *browserPool) release(
	ctx context.Context,
	slot *browserPoolSlot,
	options PlaywrightExecutorOptions,
	result ScenarioResult,
	runCase RunCase,
	runID string,
) {
	if p == nil || slot == nil {
		return
	}
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()
	if closed {
		if slot.stop != nil {
			slot.stop()
		}
		return
	}

	needsReplace := ScenarioBlocksSessionReuse(result, runCase, slot.session)
	if !needsReplace {
		if err := slot.session.resetForScenario(); err != nil {
			needsReplace = true
			logPoolResetFailure(slot, err, runCase, runID)
		}
	}
	if needsReplace {
		if err := slot.replaceSession(ctx, options); err != nil {
			p.retireSlot(slot, err)
			return
		}
		logx.Debug("pool slot replaced after browser close",
			"run_id", runID,
			"worker_id", slot.index,
			"case_id", runCase.CaseID,
			"scenario", runCase.Name,
		)
	}

	select {
	case p.slots <- slot:
	default:
		logx.Debug("pool release skipped", "error", "slot channel is full", "slot", slot.index)
	}
}

func (slot *browserPoolSlot) replaceSession(ctx context.Context, options PlaywrightExecutorOptions) error {
	if slot == nil {
		return fmt.Errorf("pool slot is nil")
	}
	if slot.stop != nil {
		slot.stop()
	}
	next, err := newPoolSlot(ctx, options, slot.index)
	if err != nil {
		return err
	}
	slot.session = next.session
	slot.stop = next.stop
	return nil
}

func (p *browserPool) retireSlot(slot *browserPoolSlot, reason error) {
	if p == nil || slot == nil {
		return
	}
	if slot.stop != nil {
		slot.stop()
	}
	p.mu.Lock()
	p.retired++
	retired := p.retired
	size := p.size
	if retired >= size && !p.closed {
		p.closed = true
		close(p.slots)
	}
	p.mu.Unlock()
	logx.Warn("pool slot retired",
		"slot", slot.index,
		"retired", retired,
		"size", size,
		"error", reason,
		"session", slot.session.poolDiagState(),
	)
}

func logPoolResetFailure(slot *browserPoolSlot, err error, runCase RunCase, runID string) {
	state := "nil"
	browserNil, contextNil, pageNil, sessionClosed := true, true, true, true
	if slot != nil && slot.session != nil {
		state = slot.session.poolDiagState()
		browserNil, contextNil, pageNil, sessionClosed, _ = slot.session.poolDiagFlags()
	}
	logx.Warn("pool reset failed",
		"run_id", runID,
		"worker_id", slot.index,
		"case_id", runCase.CaseID,
		"scenario", runCase.Name,
		"error", err,
		"browser_nil", browserNil,
		"context_nil", contextNil,
		"page_nil", pageNil,
		"session_closed", sessionClosed,
		"session", state,
	)
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
			if slot != nil {
				select {
				case p.slots <- slot:
				default:
				}
			}
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
	size := p.size
	retired := p.retired
	active := size - retired
	if active < 0 {
		active = 0
	}
	p.mu.Unlock()

	started := time.Now()
	logx.Debug("browser pool closing", "workers", size, "retired", retired, "active", active)
	for i := 0; i < active; i++ {
		slotStart := time.Now()
		select {
		case slot := <-p.slots:
			if slot != nil && slot.stop != nil {
				slot.stop()
			}
		case <-time.After(2 * time.Second):
			logx.Warn("browser pool close timed out waiting for idle slot",
				"index", i,
				"elapsed_ms", time.Since(slotStart).Milliseconds(),
			)
		}
	}
	logx.Debug("browser pool closed",
		"elapsed_ms", time.Since(started).Milliseconds(),
		"retired", retired,
	)
	p.size = 0
}

func poolEligible(options PlaywrightExecutorOptions) bool {
	return options.TraceDir == "" && options.VideoDir == ""
}

func poolEligibleForPlan(options PlaywrightExecutorOptions, plan ExecutionPlan) bool {
	if !poolEligible(options) {
		return false
	}
	if PlanContainsCloseBrowser(plan) {
		logx.Debug("browser pool disabled", "reason", "plan contains close-browser")
		return false
	}
	return true
}
