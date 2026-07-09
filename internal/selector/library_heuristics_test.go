package selector

import (
	"strings"
	"testing"
)

func TestLibraryHeuristicsJS(t *testing.T) {
	js := LibraryHeuristicsJS(true, false)
	if !strings.Contains(js, `mui: true`) || !strings.Contains(js, `ant: false`) {
		t.Fatalf("unexpected js: %s", js)
	}
}
