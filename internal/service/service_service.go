package service

import (
	"sort"
	"time"

	"tracing/internal/model"
	"tracing/pkg/idgen"
)

// CreateService 创建被观测服务。
func (s *Service) CreateService(input model.Service) (*model.Service, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateService(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建服务 %s (%s)", input.Name, input.ID)
	return &input, nil
}

// GetService 按 ID 获取服务。
func (s *Service) GetService(id string) (*model.Service, error) {
	return s.store.GetService(id)
}

// ListServices 分页列出服务。
func (s *Service) ListServices(filter model.ServiceFilter, page, size int) ([]*model.Service, int, error) {
	all := s.store.ListServices()
	matched := make([]*model.Service, 0, len(all))
	for _, svc := range all {
		if filter.Match(svc) {
			matched = append(matched, svc)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	return paginateServices(matched, page, size)
}

func paginateServices(items []*model.Service, page, size int) ([]*model.Service, int, error) {
	total := len(items)
	start := (page - 1) * size
	if start >= total {
		return []*model.Service{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return items[start:end], total, nil
}

// UpdateService 更新服务信息。
func (s *Service) UpdateService(id string, input model.Service) (*model.Service, error) {
	existing, err := s.store.GetService(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Host = input.Host
	existing.Port = input.Port
	existing.Description = input.Description
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateService(existing); err != nil {
		return nil, err
	}
	s.log.Infof("更新服务 %s", id)
	return existing, nil
}

// DeleteService 删除服务。
func (s *Service) DeleteService(id string) error {
	if err := s.store.DeleteService(id); err != nil {
		return err
	}
	s.log.Infof("删除服务 %s", id)
	return nil
}
