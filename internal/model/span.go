package model

import (
	"strings"
	"time"
)

// Span 状态常量。
const (
	SpanRunning   = "running"
	SpanCompleted = "completed"
	SpanFailed    = "failed"
)

// spanTransitions 定义跨度合法状态流转。
var spanTransitions = map[string]map[string]bool{
	SpanRunning: {SpanCompleted: true, SpanFailed: true},
}

// CanTransitionSpan 判断跨度状态能否从 from 流转到 to。
func CanTransitionSpan(from, to string) bool {
	if m, ok := spanTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Span 表示轨迹中的单个操作片段。
type Span struct {
	ID           string    `json:"id"`
	TraceID      string    `json:"trace_id"`
	ParentSpanID string    `json:"parent_span_id"`
	ServiceID    string    `json:"service_id"`
	Operation    string    `json:"operation"`
	Status       string    `json:"status"`
	StartTime    time.Time `json:"start_time"`
	DurationMs   int64     `json:"duration_ms"`
	Tags         string    `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
}

// Validate 校验跨度字段。
func (s *Span) Validate() error {
	s.TraceID = strings.TrimSpace(s.TraceID)
	s.ParentSpanID = strings.TrimSpace(s.ParentSpanID)
	s.ServiceID = strings.TrimSpace(s.ServiceID)
	s.Operation = strings.TrimSpace(s.Operation)
	s.Tags = strings.TrimSpace(s.Tags)
	if s.TraceID == "" {
		return NewValidationError("trace_id", "轨迹 ID 不能为空")
	}
	if s.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if s.Operation == "" {
		return NewValidationError("operation", "操作名称不能为空")
	}
	if s.Status == "" {
		s.Status = SpanRunning
	}
	if s.Status != SpanRunning && s.Status != SpanCompleted && s.Status != SpanFailed {
		return NewValidationError("status", "跨度状态不合法")
	}
	if s.DurationMs < 0 {
		return NewValidationError("duration_ms", "耗时不能为负数")
	}
	return nil
}

// SpanFilter 跨度列表筛选条件。
type SpanFilter struct {
	TraceID   string
	ServiceID string
	Status    string
	Operation string
}

// Match 判断跨度是否匹配筛选条件。
func (f SpanFilter) Match(s *Span) bool {
	if f.TraceID != "" && s.TraceID != f.TraceID {
		return false
	}
	if f.ServiceID != "" && s.ServiceID != f.ServiceID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Operation != "" && !strings.Contains(s.Operation, f.Operation) {
		return false
	}
	return true
}
