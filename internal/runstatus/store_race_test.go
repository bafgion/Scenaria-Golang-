package runstatus

import (
	"fmt"
	"sync"
	"testing"
)

func TestStoreParallelRecordPreservesEntries(t *testing.T) {
	tmp := t.TempDir()
	store, err := Open(tmp)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	const workers = 16
	const perWorker = 8
	var wg sync.WaitGroup
	errCh := make(chan error, workers*perWorker)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				path := fmt.Sprintf("feature-%d-%d.feature", worker, i)
				if err := store.Record(Entry{Path: path, Success: true, Runner: "playwright"}); err != nil {
					errCh <- err
				}
			}
		}(w)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("Record failed: %v", err)
	}

	entries, err := store.List(0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	want := workers * perWorker
	if len(entries) != want {
		t.Fatalf("expected %d entries, got %d", want, len(entries))
	}

	seen := make(map[string]int, want)
	for _, entry := range entries {
		seen[entry.Path]++
	}
	if len(seen) != want {
		t.Fatalf("expected %d unique paths, got %d (possible lost updates)", want, len(seen))
	}
}
