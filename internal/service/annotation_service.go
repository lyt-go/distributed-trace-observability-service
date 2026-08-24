package service

import (
	"sort"
	"time"

	"tracing/internal/model"
	"tracing/pkg/idgen"
)

// CreateAnnotation 创建标注，校验所属跨度存在且轨迹一致。
func (s *Service) CreateAnnotation(input model.Annotation) (*model.Annotation, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	span, err := s.store.GetSpan(input.SpanID)
	if err != nil {
		return nil, err
	}
	if span.TraceID != input.TraceID {
		return nil, model.NewValidationError("trace_id", "标注轨迹 ID 与跨度不一致")
	}
	input.ID = idgen.HexN(8)
	now := time.Now()
	input.CreatedAt = now
	if input.Timestamp.IsZero() {
		input.Timestamp = now
	}
	if err := s.store.CreateAnnotation(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建标注 %s", input.ID)
	return &input, nil
}

// GetAnnotation 按 ID 获取标注。
func (s *Service) GetAnnotation(id string) (*model.Annotation, error) {
	return s.store.GetAnnotation(id)
}

// ListAnnotations 分页列出标注。
func (s *Service) ListAnnotations(filter model.AnnotationFilter, page, size int) ([]*model.Annotation, int, error) {
	all := s.store.ListAnnotations()
	matched := make([]*model.Annotation, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.After(matched[j].Timestamp)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Annotation{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateAnnotation 更新标注内容。
func (s *Service) UpdateAnnotation(id string, input model.Annotation) (*model.Annotation, error) {
	existing, err := s.store.GetAnnotation(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Message = input.Message
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAnnotation(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteAnnotation 删除标注。
func (s *Service) DeleteAnnotation(id string) error {
	if err := s.store.DeleteAnnotation(id); err != nil {
		return err
	}
	s.log.Infof("删除标注 %s", id)
	return nil
}
