package model

import (
	"strings"
	"time"
)

// Service 状态常量。
const (
	ServiceActive   = "active"
	ServiceInactive = "inactive"
)

// Service 表示一个接入链路追踪的被观测服务节点。
type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化服务字段。
func (s *Service) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Host = strings.TrimSpace(s.Host)
	s.Description = strings.TrimSpace(s.Description)
	if s.Name == "" {
		return NewValidationError("name", "服务名称不能为空")
	}
	if s.Host == "" {
		return NewValidationError("host", "主机地址不能为空")
	}
	if s.Port < 1 || s.Port > 65535 {
		return NewValidationError("port", "端口必须在 1-65535 之间")
	}
	if s.Status == "" {
		s.Status = ServiceActive
	}
	if s.Status != ServiceActive && s.Status != ServiceInactive {
		return NewValidationError("status", "服务状态不合法")
	}
	return nil
}

// ServiceFilter 服务列表筛选条件。
type ServiceFilter struct {
	Status  string
	Keyword string
}

// Match 判断服务是否匹配筛选条件。
func (f ServiceFilter) Match(s *Service) bool {
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Host), k) {
			return false
		}
	}
	return true
}
