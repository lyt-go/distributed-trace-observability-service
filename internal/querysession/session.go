package querysession

import (
	"context"
	"sync"
)

type Session struct {
	mu  sync.Mutex
	ctx context.Context
}

func (s *Session) Bind(ctx context.Context) context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ctx == nil {
		s.ctx = ctx
	}
	return s.ctx
}
