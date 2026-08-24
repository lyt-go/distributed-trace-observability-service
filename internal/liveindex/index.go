package liveindex

import "sync"

type Index struct {
	mu     sync.RWMutex
	spans  map[string]map[string]string
}

func New() *Index { return &Index{spans: make(map[string]map[string]string)} }

func (i *Index) Put(traceID, spanID, operation string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	byTrace := i.spans[traceID]
	if byTrace == nil {
		byTrace = make(map[string]string)
		i.spans[traceID] = byTrace
	}
	byTrace[spanID] = operation
}

func (i *Index) Snapshot() map[string]map[string]string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.spans
}
