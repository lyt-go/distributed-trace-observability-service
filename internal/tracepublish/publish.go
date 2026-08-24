package tracepublish

import "tracing/internal/tracejournal"

type Publisher interface{ Publish(string) error }
type Audit interface{ Success(string) }

func Publish(tx *tracejournal.Transaction, publisher Publisher, audit Audit, traceID string) (err error) {
	defer func() { err = tx.Close(err) }()
	audit.Success(traceID)
	if err = publisher.Publish(traceID); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
