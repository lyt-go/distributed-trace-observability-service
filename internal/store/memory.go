package store

import (
	"sync"

	"tracing/internal/model"
)

// MemoryStore 基于内存 map 的线程安全存储实现。
type MemoryStore struct {
	mu            sync.RWMutex
	services      map[string]*model.Service
	traces        map[string]*model.Trace
	spans         map[string]*model.Span
	annotations   map[string]*model.Annotation
	samplingRules map[string]*model.SamplingRule
	slowRules     map[string]*model.SlowTraceRule
}

// NewMemoryStore 创建空的 MemoryStore。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		services:      make(map[string]*model.Service),
		traces:        make(map[string]*model.Trace),
		spans:         make(map[string]*model.Span),
		annotations:   make(map[string]*model.Annotation),
		samplingRules: make(map[string]*model.SamplingRule),
		slowRules:     make(map[string]*model.SlowTraceRule),
	}
}

var _ Store = (*MemoryStore)(nil)
