package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTempFeatureResourcesBoundToRunTargets(t *testing.T) {
	svc := NewService()
	pathA, err := svc.WriteTempFeature("Feature: A")
	if err != nil {
		t.Fatalf("WriteTempFeature A failed: %v", err)
	}
	pathB, err := svc.WriteTempFeature("Feature: B")
	if err != nil {
		t.Fatalf("WriteTempFeature B failed: %v", err)
	}
	dirA := filepath.Dir(pathA)
	dirB := filepath.Dir(pathB)

	resources := svc.tempFeatureDirsForTargets([]string{pathA})
	if len(resources) != 1 || resources[0] != dirA {
		t.Fatalf("expected only dirA in resources, got %#v", resources)
	}

	svc.cleanupTempFeatureResources(resources)
	if _, err := os.Stat(dirA); !os.IsNotExist(err) {
		t.Fatalf("expected dirA removed, got err=%v", err)
	}
	if _, err := os.Stat(dirB); err != nil {
		t.Fatalf("expected dirB to remain, got err=%v", err)
	}

	svc.cleanupTempFeatureDirs()
}
