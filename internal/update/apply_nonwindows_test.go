//go:build !windows

package update

import (
	"strings"
	"testing"
)

func TestApplyUnsupportedOnNonWindows(t *testing.T) {
	_, err := Apply("0.1.0", "/tmp/scenaria", 1)
	if err == nil || !strings.Contains(err.Error(), "Windows") {
		t.Fatalf("unexpected error: %v", err)
	}
}
