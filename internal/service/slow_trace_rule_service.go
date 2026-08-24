package service

import (
	"sort"
	"time"

	"tracing/internal/model"
	"tracing/pkg/idgen"
)

// CreateSlowTraceRule 创建慢轨迹检测规则。
func (s *Service) CreateSlowTraceRule(input model.SlowTraceRule) (*model.SlowTraceRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateSlowTraceRule(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建慢轨迹规则 %s", input.ID)
	return &input, nil
}

// GetSlowTraceRule 按 ID 获取规则。
func (s *Service) GetSlowTraceRule(id string) (*model.SlowTraceRule, error) {
	return s.store.GetSlowTraceRule(id)
}

// ListSlowTraceRules 分页列出规则。
func (s *Service) ListSlowTraceRules(filter model.SlowTraceRuleFilter, page, size int) ([]*model.SlowTraceRule, int, error) {
	all := s.store.ListSlowTraceRules()
	matched := make([]*model.SlowTraceRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.SlowTraceRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateSlowTraceRule 更新规则。
func (s *Service) UpdateSlowTraceRule(id string, input model.SlowTraceRule) (*model.SlowTraceRule, error) {
	existing, err := s.store.GetSlowTraceRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.ServiceID = input.ServiceID
	existing.ThresholdMs = input.ThresholdMs
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateSlowTraceRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteSlowTraceRule 删除规则。
func (s *Service) DeleteSlowTraceRule(id string) error {
	if err := s.store.DeleteSlowTraceRule(id); err != nil {
		return err
	}
	s.log.Infof("删除慢轨迹规则 %s", id)
	return nil
}

// DetectSlowTraces 依据全部启用规则检测慢轨迹，返回去重后的命中轨迹。
func (s *Service) DetectSlowTraces() ([]*model.Trace, error) {
	rules := s.store.ListSlowTraceRules()
	seen := make(map[string]bool)
	hit := make([]*model.Trace, 0)
	for _, t := range s.store.ListTraces() {
		for _, r := range rules {
			if r.Match(t) && !seen[t.ID] {
				seen[t.ID] = true
				hit = append(hit, t)
				break
			}
		}
	}
	sort.Slice(hit, func(i, j int) bool {
		return hit[i].DurationMs > hit[j].DurationMs
	})
	return hit, nil
}
