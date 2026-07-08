package healthx

import (
	"context"
	"testing"
)

func TestCheckerSeparatesLivenessAndReadiness(t *testing.T) {
	checker := NewChecker()
	checker.AddLiveness("process", func(context.Context) error { return nil })
	checker.AddReadiness("database", func(context.Context) error { return ErrNotReady("migration running") })

	live := checker.Liveness(context.Background())
	if !live.OK {
		t.Fatalf("liveness should pass: %+v", live)
	}
	ready := checker.Readiness(context.Background())
	if ready.OK || ready.Checks[0].Error != "migration running" {
		t.Fatalf("readiness should fail with reason: %+v", ready)
	}
}
