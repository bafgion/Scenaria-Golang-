package player

import (
	"context"
	"testing"
	"time"
)

func TestPressKeyRespectsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := pressKey(ctx, nil, "Enter")
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestExpectDownloadRespectsContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Millisecond)
	_, err := expectDownload(ctx, nil, func() error { return nil })
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestCtxAwareEmailPromptCancels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	block := make(chan struct{})
	prompt := ctxAwareEmailPrompt(ctx, func(email string) (string, error) {
		<-block
		return "1234", nil
	})
	cancel()
	if _, err := prompt("user@example.com"); err == nil {
		t.Fatal("expected context error")
	}
	close(block)
}
