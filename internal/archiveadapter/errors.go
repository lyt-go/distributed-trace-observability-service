package archiveadapter

import "errors"

var (
	ErrRejected  = errors.New("archive rejected")
	ErrTemporary = errors.New("archive temporarily unavailable")
)

func Normalize(err error) error {
	if err == nil {
		return nil
	}
	return errors.New(err.Error())
}
