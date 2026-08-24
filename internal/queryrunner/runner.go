package queryrunner

import (
	"context"

	"tracing/internal/querysession"
	"tracing/internal/traceclient"
)

type Runner struct{ session *querysession.Session }

func New(session *querysession.Session) *Runner { return &Runner{session: session} }

func (r *Runner) Fetch(ctx context.Context, client *traceclient.Client) error {
	bound := r.session.Bind(ctx)
	if err := bound.Err(); err != nil {
		return err
	}
	return client.Fetch(bound)
}
