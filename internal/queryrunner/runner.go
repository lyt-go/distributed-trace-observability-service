package queryrunner

import (
	"context"

	"tracing/internal/querysession"
	"tracing/internal/traceclient"
)

type Runner struct{ session *querysession.Session }

func New(session *querysession.Session) *Runner { return &Runner{session: session} }

// Fetch 执行一次下游查询。
//
// 每次查询都通过 session 透传自己的请求上下文，因此只受自身的取消/
// 超时影响；上一次查询结束后其上下文不会带入本次查询。bound 永远
// 等价于传入的 ctx，保留 Bind 调用以维持 session 在调用链中的参与，
// 并在进入下游前对请求上下文做一次防御性的已取消检查。
func (r *Runner) Fetch(ctx context.Context, client *traceclient.Client) error {
	bound := r.session.Bind(ctx)
	if err := bound.Err(); err != nil {
		return err
	}
	return client.Fetch(bound)
}
