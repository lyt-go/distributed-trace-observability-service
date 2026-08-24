package service

import (
	"sort"
	"time"

	"tracing/internal/model"
	"tracing/pkg/idgen"
)

// CreateSpan 创建跨度，校验所属轨迹存在。
func (s *Service) CreateSpan(input model.Span) (*model.Span, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTrace(input.TraceID); err != nil {
		return nil, err
	}
	if _, err := s.store.GetService(input.ServiceID); err != nil {
		return nil, err
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	if input.StartTime.IsZero() {
		input.StartTime = now
	}
	if input.Status == "" {
		input.Status = model.SpanRunning
	}
	if err := s.store.CreateSpan(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建跨度 %s", input.ID)
	return &input, nil
}

// GetSpan 按 ID 获取跨度。
func (s *Service) GetSpan(id string) (*model.Span, error) {
	return s.store.GetSpan(id)
}

// ListSpans 分页列出跨度。
func (s *Service) ListSpans(filter model.SpanFilter, page, size int) ([]*model.Span, int, error) {
	all := s.store.ListSpans()
	matched := make([]*model.Span, 0, len(all))
	for _, sp := range all {
		if filter.Match(sp) {
			matched = append(matched, sp)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].StartTime.After(matched[j].StartTime)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Span{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateSpan 更新跨度可编辑字段。
func (s *Service) UpdateSpan(id string, input model.Span) (*model.Span, error) {
	existing, err := s.store.GetSpan(id)
	if err != nil {
		return nil, err
	}
	existing.Operation = input.Operation
	existing.Tags = input.Tags
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSpan(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// FinishSpan 结束跨度，带状态机校验。
func (s *Service) FinishSpan(id, status string, durationMs int64) (*model.Span, error) {
	existing, err := s.store.GetSpan(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSpan(existing.Status, status) {
		return nil, model.NewValidationError("status", "非法状态流转: "+existing.Status+" -> "+status)
	}
	existing.Status = status
	existing.DurationMs = durationMs
	if err := s.store.UpdateSpan(existing); err != nil {
		return nil, err
	}
	s.log.Infof("跨度 %s 结束，状态 %s", id, status)
	return existing, nil
}

// DeleteSpan 删除跨度。
func (s *Service) DeleteSpan(id string) error {
	if err := s.store.DeleteSpan(id); err != nil {
		return err
	}
	s.log.Infof("删除跨度 %s", id)
	return nil
}
