package service

import (
	"sort"
	"time"

	"tracing/internal/model"
	"tracing/pkg/idgen"
)

// CreateTrace 创建一条调用轨迹。
func (s *Service) CreateTrace(input model.Trace) (*model.Trace, error) {
	if err := input.Validate(); err != nil {
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
		input.Status = model.TraceRunning
	}
	if err := s.store.CreateTrace(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建轨迹 %s", input.ID)
	return &input, nil
}

// GetTrace 按 ID 获取轨迹。
func (s *Service) GetTrace(id string) (*model.Trace, error) {
	return s.store.GetTrace(id)
}

// ListTraces 分页列出轨迹。
func (s *Service) ListTraces(filter model.TraceFilter, page, size int) ([]*model.Trace, int, error) {
	all := s.store.ListTraces()
	matched := make([]*model.Trace, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].StartTime.After(matched[j].StartTime)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Trace{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateTrace 更新轨迹的可编辑字段。
func (s *Service) UpdateTrace(id string, input model.Trace) (*model.Trace, error) {
	existing, err := s.store.GetTrace(id)
	if err != nil {
		return nil, err
	}
	existing.Operation = input.Operation
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTrace(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// FinishTrace 结束轨迹（completed/failed），带状态机校验。
func (s *Service) FinishTrace(id, status string, durationMs int64) (*model.Trace, error) {
	existing, err := s.store.GetTrace(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionTrace(existing.Status, status) {
		return nil, model.NewValidationError("status", "非法状态流转: "+existing.Status+" -> "+status)
	}
	existing.Status = status
	existing.DurationMs = durationMs
	existing.SpanCount = s.countSpans(id)
	if err := s.store.UpdateTrace(existing); err != nil {
		return nil, err
	}
	s.log.Infof("轨迹 %s 结束，状态 %s", id, status)
	return existing, nil
}

func (s *Service) countSpans(traceID string) int {
	n := 0
	for _, sp := range s.store.ListSpans() {
		if sp.TraceID == traceID {
			n++
		}
	}
	return n
}

// DeleteTrace 删除轨迹。
func (s *Service) DeleteTrace(id string) error {
	if err := s.store.DeleteTrace(id); err != nil {
		return err
	}
	s.log.Infof("删除轨迹 %s", id)
	return nil
}
