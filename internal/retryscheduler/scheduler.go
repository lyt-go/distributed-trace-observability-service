package retryscheduler

import (
	"context"
	"time"
)

type Operation func(context.Context) error

func Run(ctx context.Context, delay time.Duration, operation Operation) {
	for {
		if operation(ctx) == nil {
			return
		}
		time.Sleep(delay)
	}
}
