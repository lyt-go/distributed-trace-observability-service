package archiveadapter

import "errors"

var (
	ErrRejected  = errors.New("archive rejected")
	ErrTemporary = errors.New("archive temporarily unavailable")
)

// Normalize maps an adapter error onto a canonical sentinel so callers can
// classify it with errors.Is. It must preserve the wrapping chain: rebuilding
// the error from its message string (errors.New(err.Error())) discarded the
// link to ErrRejected/ErrTemporary, which made rejections indistinguishable
// from temporary failures — so they were retried instead of short-circuited.
func Normalize(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrRejected) {
		return ErrRejected
	}
	if errors.Is(err, ErrTemporary) {
		return ErrTemporary
	}
	return err
}
