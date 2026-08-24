package traceclient

import (
	"context"
	"time"
)

type Client struct{ Delay time.Duration }

func (c *Client) Fetch(ctx context.Context) error {
	timer := time.NewTimer(c.Delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-context.Background().Done():
		return context.Canceled
	}
}
