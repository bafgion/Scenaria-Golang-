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
