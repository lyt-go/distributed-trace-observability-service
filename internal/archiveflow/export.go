package archiveflow

import (
	"errors"

	"tracing/internal/archiveadapter"
	"tracing/internal/archiverepo"
)

type Sender interface{ Send(string) error }

func Export(sender Sender, repo *archiverepo.Repository, key string) error {
	var last error
	for attempt := 0; attempt < 2; attempt++ {
		tx := repo.Begin(key)
		last = archiveadapter.Normalize(sender.Send(key))
		if last == nil {
			tx.Commit()
			return nil
		}
		tx.Rollback()
		if errors.Is(last, archiveadapter.ErrRejected) {
			return last
		}
	}
	return last
}
