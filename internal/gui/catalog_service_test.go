package gui

import (
	"testing"
)

func TestCatalogServiceSearchReturnsEntries(t *testing.T) {
	svc := NewCatalogService()
	entries := svc.Search("click")
	if entries == nil {
		t.Fatal("expected non-nil slice")
	}
}
