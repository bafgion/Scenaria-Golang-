package player

import (
	"context"
	"testing"
)

func TestContinueOnFail(t *testing.T) {
	ctx := WithContinueOnFail(context.Background(), true)
	if !ContinueOnFail(ctx) {
		t.Fatal("expected continue on fail")
	}
	if ContinueOnFail(context.Background()) {
		t.Fatal("expected default false")
	}
}
