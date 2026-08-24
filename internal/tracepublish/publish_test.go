package tracepublish

import (
	"errors"
	"testing"

	"tracing/internal/tracejournal"
)

type publisherSpy struct{ calls int }
func (p *publisherSpy) Publish(string) error { p.calls++; return nil }

type auditSpy struct{ traces []string }
func (a *auditSpy) Success(traceID string) { a.traces = append(a.traces, traceID) }

func TestCommitFailureKeepsEventsPrivate(t *testing.T) {
	failedPublisher := &publisherSpy{}
	failedAudit := &auditSpy{}
	tx := &tracejournal.Transaction{CommitError: tracejournal.ErrCommit, RollbackError: tracejournal.ErrRollback}
	err := Publish(tx, failedPublisher, failedAudit, "trace-failed")
	if !errors.Is(err, tracejournal.ErrCommit) {
		t.Errorf("commit failure returned %v, want original commit error", err)
	}
	if failedPublisher.calls != 0 {
		t.Errorf("publisher called %d times before failed commit", failedPublisher.calls)
	}
	if len(failedAudit.traces) != 0 {
		t.Errorf("failed commit emitted success audit: %v", failedAudit.traces)
	}

	successPublisher := &publisherSpy{}
	successAudit := &auditSpy{}
	if err := Publish(&tracejournal.Transaction{}, successPublisher, successAudit, "trace-ok"); err != nil {
		t.Errorf("successful publish returned %v", err)
	}
	if successPublisher.calls != 1 || len(successAudit.traces) != 1 {
		t.Errorf("successful publish calls=%d audit=%v", successPublisher.calls, successAudit.traces)
	}
}
