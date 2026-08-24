package queryrunner

import (
	"context"
	"errors"
	"testing"
	"time"

	"tracing/internal/querysession"
	"tracing/internal/traceclient"
)

func TestDeadlinesCancelAndDoNotLeakToNextQuery(t *testing.T) {
	runner := New(&querysession.Session{})
	firstCtx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	started := time.Now()
	firstErr := runner.Fetch(firstCtx, &traceclient.Client{Delay: 120 * time.Millisecond})
	if !errors.Is(firstErr, context.DeadlineExceeded) {
		t.Errorf("timed out query returned %v, want deadline exceeded", firstErr)
	}
	if elapsed := time.Since(started); elapsed > 70*time.Millisecond {
		t.Errorf("timed out query kept downstream work alive for %v", elapsed)
	}

	secondCtx, stop := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer stop()
	if err := runner.Fetch(secondCtx, &traceclient.Client{Delay: time.Millisecond}); err != nil {
		t.Errorf("fresh query inherited the earlier deadline: %v", err)
	}
}
