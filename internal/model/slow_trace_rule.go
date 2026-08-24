package model

import (
	"strings"
	"time"
)

// SlowTraceRule 状态常量。
const (
	SlowRuleEnabled  = "enabled"
	SlowRuleDisabled = "disabled"
)

// SlowTraceRule 表示一条慢轨迹检测规则，耗时超过阈值即触发。
type SlowTraceRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ServiceID   string    `json:"service_id"`
	ThresholdMs int64     `json:"threshold_ms"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验慢轨迹规则字段。
func (r *SlowTraceRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.ServiceID = strings.TrimSpace(r.ServiceID)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.ThresholdMs <= 0 {
		return NewValidationError("threshold_ms", "耗时阈值必须大于 0")
	}
	if r.Status == "" {
		r.Status = SlowRuleEnabled
	}
	if r.Status != SlowRuleEnabled && r.Status != SlowRuleDisabled {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

// Match 判断轨迹是否命中该慢轨迹规则。
func (r *SlowTraceRule) Match(t *Trace) bool {
	if r.Status != SlowRuleEnabled {
		return false
	}
	if r.ServiceID != "" && t.ServiceID != r.ServiceID {
		return false
	}
	return t.DurationMs >= r.ThresholdMs
}

// SlowTraceRuleFilter 慢轨迹规则筛选条件。
type SlowTraceRuleFilter struct {
	ServiceID string
	Status    string
}

// Match 判断规则是否匹配筛选条件。
func (f SlowTraceRuleFilter) Match(r *SlowTraceRule) bool {
	if f.ServiceID != "" && r.ServiceID != f.ServiceID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
