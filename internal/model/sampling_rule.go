package model

import (
	"strings"
	"time"
)

// SamplingRule 状态常量。
const (
	SamplingEnabled  = "enabled"
	SamplingDisabled = "disabled"
)

// SamplingRule 表示一条采样策略，用于决定哪些轨迹需要被记录。
type SamplingRule struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	ServicePattern string    `json:"service_pattern"`
	Rate           int       `json:"rate"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Validate 校验采样规则字段。
func (r *SamplingRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.ServicePattern = strings.TrimSpace(r.ServicePattern)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.ServicePattern == "" {
		return NewValidationError("service_pattern", "服务匹配模式不能为空")
	}
	if r.Rate < 0 || r.Rate > 100 {
		return NewValidationError("rate", "采样率必须在 0-100 之间")
	}
	if r.Status == "" {
		r.Status = SamplingEnabled
	}
	if r.Status != SamplingEnabled && r.Status != SamplingDisabled {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

// ShouldSample 按采样率判断是否采样。
func (r *SamplingRule) ShouldSample(hash int) bool {
	if r.Status != SamplingEnabled || r.Rate <= 0 {
		return false
	}
	if r.Rate >= 100 {
		return true
	}
	return hash%100 < r.Rate
}

// SamplingRuleFilter 采样规则筛选条件。
type SamplingRuleFilter struct {
	Status  string
	Keyword string
}

// Match 判断规则是否匹配筛选条件。
func (f SamplingRuleFilter) Match(r *SamplingRule) bool {
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) &&
			!strings.Contains(strings.ToLower(r.ServicePattern), k) {
			return false
		}
	}
	return true
}
