package store

import (
	"tracing/internal/model"
)

func (s *MemoryStore) CreateService(svc *model.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.services {
		if exist.Name == svc.Name {
			return ErrConflict
		}
	}
	s.services[svc.ID] = svc
	return nil
}

func (s *MemoryStore) GetService(id string) (*model.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	svc, ok := s.services[id]
	if !ok {
		return nil, ErrNotFound
	}
	return svc, nil
}

func (s *MemoryStore) GetServiceByName(name string) (*model.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, svc := range s.services {
		if svc.Name == name {
			return svc, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListServices() []*model.Service {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Service, 0, len(s.services))
	for _, svc := range s.services {
		list = append(list, svc)
	}
	return list
}

func (s *MemoryStore) UpdateService(svc *model.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[svc.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.services {
		if exist.ID != svc.ID && exist.Name == svc.Name {
			return ErrConflict
		}
	}
	s.services[svc.ID] = svc
	return nil
}

func (s *MemoryStore) DeleteService(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.services[id]; !ok {
		return ErrNotFound
	}
	delete(s.services, id)
	return nil
}
