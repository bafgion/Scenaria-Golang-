package player

import (
	"context"
	"testing"
	"time"
)

type urlStubPage struct {
	urls []string
	idx  int
}

func (p *urlStubPage) URL() string {
	if p.idx >= len(p.urls) {
		return p.urls[len(p.urls)-1]
	}
	u := p.urls[p.idx]
	p.idx++
	return u
}

func TestWaitForURLMatch(t *testing.T) {
	page := &urlStubPage{urls: []string{"about:blank", "https://example.com/page"}}
	if err := waitForURL(context.Background(), page, "https://example.com/page", false); err != nil {
		t.Fatalf("waitForURL: %v", err)
	}
}

func TestWaitForURLContains(t *testing.T) {
	page := &urlStubPage{urls: []string{"about:blank", "https://example.com/dashboard"}}
	if err := waitForURL(context.Background(), page, "dashboard", true); err != nil {
		t.Fatalf("waitForURL contains: %v", err)
	}
}

func TestWaitForURLRespectsContext(t *testing.T) {
	page := &urlStubPage{urls: []string{"about:blank"}}
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	err := waitForURL(ctx, page, "https://example.com", false)
	if err == nil {
		t.Fatal("expected timeout/cancel error")
	}
}
