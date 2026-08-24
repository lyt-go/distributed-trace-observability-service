package service

import (
	"tracing/internal/model"
)

// TraceStats 轨迹统计概览。
type TraceStats struct {
	ServiceCount    int                    `json:"service_count"`
	TraceCount      int                    `json:"trace_count"`
	SpanCount       int                    `json:"span_count"`
	AnnotationCount int                    `json:"annotation_count"`
	RunningCount    int                    `json:"running_count"`
	CompletedCount  int                    `json:"completed_count"`
	FailedCount     int                    `json:"failed_count"`
	AvgDurationMs   int64                  `json:"avg_duration_ms"`
	SlowTraceCount  int                    `json:"slow_trace_count"`
	ByStatus        map[string]int         `json:"by_status"`
	ByService       map[string]ServiceStat `json:"by_service"`
}

// ServiceStat 单个服务的轨迹统计。
type ServiceStat struct {
	Name         string `json:"name"`
	TraceCount   int    `json:"trace_count"`
	SpanCount    int    `json:"span_count"`
	AvgDuration  int64  `json:"avg_duration_ms"`
	FailedCount  int    `json:"failed_count"`
}

// Overview 汇总全链路观测统计。
func (s *Service) Overview() (*TraceStats, error) {
	stats := &TraceStats{
		ServiceCount:    len(s.store.ListServices()),
		ByStatus:        make(map[string]int),
		ByService:       make(map[string]ServiceStat),
	}
	spanCountByTrace := make(map[string]int)
	for _, sp := range s.store.ListSpans() {
		stats.SpanCount++
		spanCountByTrace[sp.TraceID]++
		if svc, err := s.store.GetService(sp.ServiceID); err == nil {
			st := stats.ByService[sp.ServiceID]
			st.Name = svc.Name
			st.SpanCount++
			stats.ByService[sp.ServiceID] = st
		}
	}
	var totalDur int64
	var completed int
	for _, t := range s.store.ListTraces() {
		stats.TraceCount++
		stats.ByStatus[t.Status]++
		switch t.Status {
		case model.TraceRunning:
			stats.RunningCount++
		case model.TraceCompleted:
			stats.CompletedCount++
			totalDur += t.DurationMs
			completed++
		case model.TraceFailed:
			stats.FailedCount++
		}
		if svc, err := s.store.GetService(t.ServiceID); err == nil {
			st := stats.ByService[t.ServiceID]
			st.Name = svc.Name
			st.TraceCount++
			if t.Status == model.TraceFailed {
				st.FailedCount++
			}
			if t.Status == model.TraceCompleted {
				st.AvgDuration += t.DurationMs
			}
			stats.ByService[t.ServiceID] = st
		}
		_ = spanCountByTrace
	}
	stats.AnnotationCount = len(s.store.ListAnnotations())
	if completed > 0 {
		stats.AvgDurationMs = totalDur / int64(completed)
	}
	slow, _ := s.DetectSlowTraces()
	stats.SlowTraceCount = len(slow)
	return stats, nil
}
