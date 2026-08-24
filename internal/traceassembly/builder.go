package traceassembly

import (
	"fmt"
	"strings"

	"tracing/internal/tracecache"
)

type Builder struct{ cache *tracecache.Cache }

func New(cache *tracecache.Cache) *Builder { return &Builder{cache: cache} }

func (b *Builder) Build(id string, segments []string) (draft *tracecache.Draft, err error) {
	draft = &tracecache.Draft{ID: id, Fields: make(map[string]string)}
	b.cache.Put(draft)
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("trace segment panic: %v", recovered)
		}
	}()
	for _, segment := range segments {
		parts := strings.SplitN(segment, "=", 2)
		draft.Fields[parts[0]] = parts[1]
	}
	return draft, nil
}
