package store

import (
	"tracing/internal/model"
)

func (s *MemoryStore) CreateSamplingRule(r *model.SamplingRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.samplingRules[r.ID]; ok {
		return ErrConflict
	}
	s.samplingRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetSamplingRule(id string) (*model.SamplingRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.samplingRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListSamplingRules() []*model.SamplingRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SamplingRule, 0, len(s.samplingRules))
	for _, r := range s.samplingRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateSamplingRule(r *model.SamplingRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.samplingRules[r.ID]; !ok {
		return ErrNotFound
	}
	s.samplingRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteSamplingRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.samplingRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.samplingRules, id)
	return nil
}
