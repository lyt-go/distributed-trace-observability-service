package traceclient

import (
	"context"
	"time"
)

type Client struct{ Delay time.Duration }

// Fetch 模拟一次下游查询。
//
// 它必须在传入的 ctx 上 select，使该查询自身的超时/取消能及时生效；
// 不得对 context.Background() 监听，否则下游会忽略请求的取消而一直
// 跑到自然结束，并让被取消的上下文泄漏到后续查询。
func (c *Client) Fetch(ctx context.Context) error {
	timer := time.NewTimer(c.Delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
