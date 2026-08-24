package model

import (
	"strings"
	"time"
)

// Trace 状态常量。
const (
	TraceRunning   = "running"
	TraceCompleted = "completed"
	TraceFailed    = "failed"
)

// traceTransitions 定义轨迹合法状态流转。
var traceTransitions = map[string]map[string]bool{
	TraceRunning: {TraceCompleted: true, TraceFailed: true},
}

// CanTransitionTrace 判断轨迹状态能否从 from 流转到 to。
func CanTransitionTrace(from, to string) bool {
	if m, ok := traceTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Trace 表示一次请求的完整调用轨迹。
type Trace struct {
	ID         string    `json:"id"`
	ServiceID  string    `json:"service_id"`
	Operation  string    `json:"operation"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	DurationMs int64     `json:"duration_ms"`
	SpanCount  int       `json:"span_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// Validate 校验轨迹字段。
func (t *Trace) Validate() error {
	t.ServiceID = strings.TrimSpace(t.ServiceID)
	t.Operation = strings.TrimSpace(t.Operation)
	if t.ServiceID == "" {
		return NewValidationError("service_id", "服务 ID 不能为空")
	}
	if t.Operation == "" {
		return NewValidationError("operation", "根操作名称不能为空")
	}
	if t.Status == "" {
		t.Status = TraceRunning
	}
	if t.Status != TraceRunning && t.Status != TraceCompleted && t.Status != TraceFailed {
		return NewValidationError("status", "轨迹状态不合法")
	}
	if t.DurationMs < 0 {
		return NewValidationError("duration_ms", "耗时不能为负数")
	}
	return nil
}

// TraceFilter 轨迹列表筛选条件。
type TraceFilter struct {
	ServiceID    string
	Status       string
	MinDuration  int64
	Keyword      string
}

// Match 判断轨迹是否匹配筛选条件。
func (f TraceFilter) Match(t *Trace) bool {
	if f.ServiceID != "" && t.ServiceID != f.ServiceID {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.MinDuration > 0 && t.DurationMs < f.MinDuration {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Operation), k) &&
			!strings.Contains(strings.ToLower(t.ID), k) {
			return false
		}
	}
	return true
}
