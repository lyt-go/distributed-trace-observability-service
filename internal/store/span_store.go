package store

import (
	"tracing/internal/model"
)

func (s *MemoryStore) CreateSpan(sp *model.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.spans[sp.ID]; ok {
		return ErrConflict
	}
	s.spans[sp.ID] = sp
	return nil
}

func (s *MemoryStore) GetSpan(id string) (*model.Span, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sp, ok := s.spans[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sp, nil
}

func (s *MemoryStore) ListSpans() []*model.Span {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Span, 0, len(s.spans))
	for _, sp := range s.spans {
		list = append(list, sp)
	}
	return list
}

func (s *MemoryStore) UpdateSpan(sp *model.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.spans[sp.ID]; !ok {
		return ErrNotFound
	}
	s.spans[sp.ID] = sp
	return nil
}

func (s *MemoryStore) DeleteSpan(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.spans[id]; !ok {
		return ErrNotFound
	}
	delete(s.spans, id)
	return nil
}
