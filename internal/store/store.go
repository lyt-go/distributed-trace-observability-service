// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"tracing/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Service
	CreateService(s *model.Service) error
	GetService(id string) (*model.Service, error)
	GetServiceByName(name string) (*model.Service, error)
	ListServices() []*model.Service
	UpdateService(s *model.Service) error
	DeleteService(id string) error

	// Trace
	CreateTrace(t *model.Trace) error
	GetTrace(id string) (*model.Trace, error)
	ListTraces() []*model.Trace
	UpdateTrace(t *model.Trace) error
	DeleteTrace(id string) error

	// Span
	CreateSpan(s *model.Span) error
	GetSpan(id string) (*model.Span, error)
	ListSpans() []*model.Span
	UpdateSpan(s *model.Span) error
	DeleteSpan(id string) error

	// Annotation
	CreateAnnotation(a *model.Annotation) error
	GetAnnotation(id string) (*model.Annotation, error)
	ListAnnotations() []*model.Annotation
	UpdateAnnotation(a *model.Annotation) error
	DeleteAnnotation(id string) error

	// SamplingRule
	CreateSamplingRule(r *model.SamplingRule) error
	GetSamplingRule(id string) (*model.SamplingRule, error)
	ListSamplingRules() []*model.SamplingRule
	UpdateSamplingRule(r *model.SamplingRule) error
	DeleteSamplingRule(id string) error

	// SlowTraceRule
	CreateSlowTraceRule(r *model.SlowTraceRule) error
	GetSlowTraceRule(id string) (*model.SlowTraceRule, error)
	ListSlowTraceRules() []*model.SlowTraceRule
	UpdateSlowTraceRule(r *model.SlowTraceRule) error
	DeleteSlowTraceRule(id string) error
}
