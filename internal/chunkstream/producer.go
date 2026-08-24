package chunkstream

import (
	"errors"
	"sync"
)

var ErrDecode = errors.New("trace chunk decode failed")

type Stream struct {
	Items <-chan string
	Errors <-chan error
	Done <-chan struct{}
}

func Start(chunks [][]string, failAt int) Stream {
	items := make(chan string)
	errs := make(chan error, 1)
	done := make(chan struct{})
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		for i, chunk := range chunks {
			if i == failAt {
				errs <- ErrDecode
				return
			}
			for _, item := range chunk {
				items <- item
			}
		}
		close(items)
	}()
	go func() {
		wg.Wait()
		close(done)
	}()
	return Stream{Items: items, Errors: errs, Done: done}
}
