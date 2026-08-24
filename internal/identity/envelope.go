package identity

import "sync"

// Envelope carries request identity into an asynchronous trace export.
type Envelope struct {
	Tenant string
	Tags   []string
}

type Pool struct{ values sync.Pool }

func NewPool() *Pool {
	p := &Pool{}
	p.values.New = func() any { return &Envelope{} }
	return p
}

func (p *Pool) Acquire() *Envelope { return p.values.Get().(*Envelope) }

func (p *Pool) Release(e *Envelope) {
	p.values.Put(e)
}

func (e *Envelope) Snapshot() *Envelope {
	return &Envelope{Tenant: e.Tenant, Tags: append([]string(nil), e.Tags...)}
}
