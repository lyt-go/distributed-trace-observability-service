package archiveflow

import (
	"errors"
	"testing"

	"tracing/internal/archiveadapter"
	"tracing/internal/archiverepo"
)

type scriptedSender struct {
	errs  []error
	calls int
}

func (s *scriptedSender) Send(string) error {
	s.calls++
	if len(s.errs) == 0 {
		return nil
	}
	err := s.errs[0]
	s.errs = s.errs[1:]
	return err
}

func TestArchiveRetryKeepsErrorTypeAndCommitBoundary(t *testing.T) {
	rejectedRepo := &archiverepo.Repository{}
	rejected := &scriptedSender{errs: []error{archiveadapter.ErrRejected}}
	err := Export(rejected, rejectedRepo, "trace-rejected")
	if !errors.Is(err, archiveadapter.ErrRejected) {
		t.Errorf("rejected archive returned %v, want typed rejection", err)
	}
	if rejected.calls != 1 {
		t.Errorf("rejected archive calls = %d, want 1", rejected.calls)
	}
	if records := rejectedRepo.Records(); len(records) != 0 {
		t.Errorf("rejected archive left records: %v", records)
	}

	temporaryRepo := &archiverepo.Repository{}
	temporary := &scriptedSender{errs: []error{archiveadapter.ErrTemporary, nil}}
	if err := Export(temporary, temporaryRepo, "trace-retry"); err != nil {
		t.Errorf("temporary archive did not recover: %v", err)
	}
	if temporary.calls != 2 {
		t.Errorf("temporary archive calls = %d, want 2", temporary.calls)
	}
	if records := temporaryRepo.Records(); len(records) != 1 || records[0] != "trace-retry" {
		t.Errorf("temporary archive records = %v, want one committed trace", records)
	}
}
