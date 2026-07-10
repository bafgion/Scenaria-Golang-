package player

import "testing"

func TestResultStatusIsSuccessful(t *testing.T) {
	if !ResultStatusIsSuccessful("passed") {
		t.Fatal("passed should be successful")
	}
	if !ResultStatusIsSuccessful("dry-run") {
		t.Fatal("dry-run should be successful")
	}
	if ResultStatusIsSuccessful("skipped") {
		t.Fatal("skipped should not be successful")
	}
	if ResultStatusIsSuccessful("canceled") {
		t.Fatal("canceled should not be successful")
	}
	if ResultStatusIsSuccessful("not-started") {
		t.Fatal("not-started should not be successful")
	}
	if ResultStatusIsSuccessful("failed") {
		t.Fatal("failed should not be successful")
	}
}
