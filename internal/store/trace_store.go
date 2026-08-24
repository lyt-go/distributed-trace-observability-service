package store

import (
	"tracing/internal/model"
)

func (s *MemoryStore) CreateTrace(t *model.Trace) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.traces[t.ID]; ok {
		return ErrConflict
	}
	s.traces[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTrace(id string) (*model.Trace, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.traces[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTraces() []*model.Trace {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Trace, 0, len(s.traces))
	for _, t := range s.traces {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTrace(t *model.Trace) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.traces[t.ID]; !ok {
		return ErrNotFound
	}
	s.traces[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTrace(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.traces[id]; !ok {
		return ErrNotFound
	}
	delete(s.traces, id)
	return nil
}
