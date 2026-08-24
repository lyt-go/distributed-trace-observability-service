package tracecache

import "sync"

type Draft struct {
	ID     string
	Fields map[string]string
}

type Cache struct {
	mu     sync.RWMutex
	drafts map[string]*Draft
}

func New() *Cache { return &Cache{drafts: make(map[string]*Draft)} }

func (c *Cache) Put(draft *Draft) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.drafts[draft.ID] = draft
}

func (c *Cache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.drafts, id)
}

func (c *Cache) Get(id string) (*Draft, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	draft, ok := c.drafts[id]
	return draft, ok
}
