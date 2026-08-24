package traceoverview

import (
	"sync"

	"tracing/internal/liveindex"
)

type Overview struct {
	mu   sync.RWMutex
	last map[string]int
}

func (o *Overview) Summarize(index *liveindex.Index, ready <-chan struct{}) map[string]int {
	snapshot := index.Snapshot()
	if ready != nil {
		<-ready
	}
	result := make(map[string]int, len(snapshot))
	for traceID, spans := range snapshot {
		result[traceID] = len(spans)
	}
	o.mu.Lock()
	o.last = result
	o.mu.Unlock()
	return result
}

func (o *Overview) Last() map[string]int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.last
}
