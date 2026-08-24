package store

import (
	"tracing/internal/model"
)

func (s *MemoryStore) CreateAnnotation(a *model.Annotation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.annotations[a.ID]; ok {
		return ErrConflict
	}
	s.annotations[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAnnotation(id string) (*model.Annotation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.annotations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAnnotations() []*model.Annotation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Annotation, 0, len(s.annotations))
	for _, a := range s.annotations {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAnnotation(a *model.Annotation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.annotations[a.ID]; !ok {
		return ErrNotFound
	}
	s.annotations[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAnnotation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.annotations[id]; !ok {
		return ErrNotFound
	}
	delete(s.annotations, id)
	return nil
}
