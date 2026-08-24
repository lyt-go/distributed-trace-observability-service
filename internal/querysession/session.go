package querysession

import "context"

// Session 提供查询上下文的透传与取消协调。
//
// 每次查询都应只受自身请求上下文的取消/超时影响，前一次查询结束后
// 不得污染后续查询。因此 Bind 不再缓存首个上下文，而是把当前请求的
// 上下文与上一次查询可能遗留的取消信号解耦：总是返回传入的 ctx，
// 保证每个查询独立。
type Session struct{}

// Bind 返回当前查询应使用的上下文。
//
// 它直接透传 ctx，使每次查询只受自己的取消时间影响；上一次查询的
// 已到期/已取消上下文不会被带入本次查询。
func (s *Session) Bind(ctx context.Context) context.Context {
	return ctx
}
