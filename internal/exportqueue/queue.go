package exportqueue

import "tracing/internal/identity"

type Sink interface {
	Submit(*identity.Envelope)
}

type Queue struct {
	pool *identity.Pool
	sink Sink
}

func New(pool *identity.Pool, sink Sink) *Queue { return &Queue{pool: pool, sink: sink} }

func (q *Queue) Enqueue(tenant string, tags []string) {
	e := q.pool.Acquire()
	e.Tenant = tenant
	e.Tags = append(e.Tags[:0], tags...)
	q.sink.Submit(e)
	q.pool.Release(e)
}
