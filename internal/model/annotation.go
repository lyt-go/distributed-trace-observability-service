package model

import (
	"strings"
	"time"
)

// Annotation 表示挂在跨度上的关键事件标注。
type Annotation struct {
	ID        string    `json:"id"`
	SpanID    string    `json:"span_id"`
	TraceID   string    `json:"trace_id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验标注字段。
func (a *Annotation) Validate() error {
	a.SpanID = strings.TrimSpace(a.SpanID)
	a.TraceID = strings.TrimSpace(a.TraceID)
	a.Name = strings.TrimSpace(a.Name)
	a.Message = strings.TrimSpace(a.Message)
	if a.SpanID == "" {
		return NewValidationError("span_id", "跨度 ID 不能为空")
	}
	if a.TraceID == "" {
		return NewValidationError("trace_id", "轨迹 ID 不能为空")
	}
	if a.Name == "" {
		return NewValidationError("name", "事件名称不能为空")
	}
	if a.Message == "" {
		return NewValidationError("message", "事件内容不能为空")
	}
	return nil
}

// AnnotationFilter 标注列表筛选条件。
type AnnotationFilter struct {
	SpanID  string
	TraceID string
	Keyword string
}

// Match 判断标注是否匹配筛选条件。
func (f AnnotationFilter) Match(a *Annotation) bool {
	if f.SpanID != "" && a.SpanID != f.SpanID {
		return false
	}
	if f.TraceID != "" && a.TraceID != f.TraceID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Name), k) &&
			!strings.Contains(strings.ToLower(a.Message), k) {
			return false
		}
	}
	return true
}
