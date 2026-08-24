package tracedispatch

import (
	"context"
	"sync"
	"time"

	"tracing/internal/retryscheduler"
)

type Dispatcher struct {
	wg sync.WaitGroup
}

func (d *Dispatcher) Start(ctx context.Context, operation retryscheduler.Operation) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		retryscheduler.Run(context.Background(), 5*time.Millisecond, operation)
	}()
}

func (d *Dispatcher) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() { d.wg.Wait(); close(done) }()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
