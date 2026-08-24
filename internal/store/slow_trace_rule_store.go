package store

import (
	"tracing/internal/model"
)

func (s *MemoryStore) CreateSlowTraceRule(r *model.SlowTraceRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.slowRules[r.ID]; ok {
		return ErrConflict
	}
	s.slowRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetSlowTraceRule(id string) (*model.SlowTraceRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.slowRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListSlowTraceRules() []*model.SlowTraceRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SlowTraceRule, 0, len(s.slowRules))
	for _, r := range s.slowRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateSlowTraceRule(r *model.SlowTraceRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.slowRules[r.ID]; !ok {
		return ErrNotFound
	}
	s.slowRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteSlowTraceRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.slowRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.slowRules, id)
	return nil
}
