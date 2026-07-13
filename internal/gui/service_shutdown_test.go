package gui

import (
	"context"
	"testing"
	"time"
)

func TestShutdownWaitsForActivePlaywright(t *testing.T) {
	svc := NewService()
	svc.activePlaywright.Add(1)
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		time.Sleep(50 * time.Millisecond)
		svc.activePlaywright.Done()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	svc.Shutdown(ctx)

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("shutdown returned before active playwright work finished")
	}
}

func TestShutdownRespectsTimeout(t *testing.T) {
	svc := NewService()
	svc.activePlaywright.Add(1)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	svc.Shutdown(ctx)
	if time.Since(start) < 15*time.Millisecond {
		t.Fatal("expected shutdown to wait for timeout")
	}
	svc.activePlaywright.Done()
}

func TestRunBoundedCleanupReturnsOnDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	started := make(chan struct{})
	release := make(chan struct{})
	defer close(release)

	start := time.Now()
	ok := runBoundedCleanup(ctx, "hung-test-cleanup", func() {
		close(started)
		<-release
	})
	if ok {
		t.Fatal("expected cleanup timeout")
	}
	if time.Since(start) > time.Second {
		t.Fatal("bounded cleanup waited too long")
	}
	select {
	case <-started:
	default:
		t.Fatal("cleanup should have started before deadline")
	}
}

func TestRunBoundedCleanupSkipsAfterDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	ok := runBoundedCleanup(ctx, "expired-test-cleanup", func() {
		called = true
	})
	if ok {
		t.Fatal("expected cleanup to be skipped")
	}
	if called {
		t.Fatal("cleanup should not start after deadline")
	}
}

func TestShutdownRepeatedWithExpiredContextDoesNotPanic(t *testing.T) {
	svc := NewService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	svc.Shutdown(ctx)
	svc.Shutdown(ctx)
}
