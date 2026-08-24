package service

import (
	"hash/fnv"
	"sort"
	"strings"
	"time"

	"tracing/internal/model"
	"tracing/pkg/idgen"
)

// CreateSamplingRule 创建采样规则。
func (s *Service) CreateSamplingRule(input model.SamplingRule) (*model.SamplingRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateSamplingRule(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建采样规则 %s", input.ID)
	return &input, nil
}

// GetSamplingRule 按 ID 获取采样规则。
func (s *Service) GetSamplingRule(id string) (*model.SamplingRule, error) {
	return s.store.GetSamplingRule(id)
}

// ListSamplingRules 分页列出采样规则。
func (s *Service) ListSamplingRules(filter model.SamplingRuleFilter, page, size int) ([]*model.SamplingRule, int, error) {
	all := s.store.ListSamplingRules()
	matched := make([]*model.SamplingRule, 0, len(all))
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
		return []*model.SamplingRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateSamplingRule 更新采样规则。
func (s *Service) UpdateSamplingRule(id string, input model.SamplingRule) (*model.SamplingRule, error) {
	existing, err := s.store.GetSamplingRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.ServicePattern = input.ServicePattern
	existing.Rate = input.Rate
	if input.Status != "" {
		existing.Status = input.Status
	}
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateSamplingRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteSamplingRule 删除采样规则。
func (s *Service) DeleteSamplingRule(id string) error {
	if err := s.store.DeleteSamplingRule(id); err != nil {
		return err
	}
	s.log.Infof("删除采样规则 %s", id)
	return nil
}

// EvaluateSampling 判断某服务名的轨迹是否应被采样，返回命中的规则与决策。
func (s *Service) EvaluateSampling(serviceName string) (*model.SamplingRule, bool, error) {
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		return nil, false, model.NewValidationError("service_name", "服务名不能为空")
	}
	for _, r := range s.store.ListSamplingRules() {
		if r.Status != model.SamplingEnabled {
			continue
		}
		if matchServicePattern(r.ServicePattern, serviceName) {
			return r, r.ShouldSample(hashService(serviceName)), nil
		}
	}
	defaultSample := 100
	if s.cfg != nil {
		defaultSample = s.cfg.DefaultSample
	}
	return nil, hashService(serviceName)%100 < defaultSample, nil
}

// matchServicePattern 支持通配符 * 的简单模式匹配。
func matchServicePattern(pattern, name string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		return strings.Contains(name, strings.Trim(pattern, "*"))
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(name, strings.TrimPrefix(pattern, "*"))
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(name, strings.TrimSuffix(pattern, "*"))
	}
	return pattern == name
}

// hashService 计算服务名的稳定哈希值（0-99）。
func hashService(name string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return int(h.Sum32() % 100)
}
