package tracejournal

import "errors"

var (
	ErrCommit   = errors.New("journal commit failed")
	ErrRollback = errors.New("journal rollback failed")
)

type Transaction struct {
	CommitError   error
	RollbackError error
}

func (t *Transaction) Commit() error   { return t.CommitError }
func (t *Transaction) Rollback() error { return t.RollbackError }

func (t *Transaction) Close(cause error) error {
	if cause == nil {
		return nil
	}
	return t.Rollback()
}
