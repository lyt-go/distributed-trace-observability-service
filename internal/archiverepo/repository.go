package archiverepo

import "sync"

type Repository struct {
	mu      sync.Mutex
	records []string
}

type Attempt struct {
	repo *Repository
	key  string
}

func (r *Repository) Begin(key string) *Attempt {
	r.mu.Lock()
	r.records = append(r.records, key)
	r.mu.Unlock()
	return &Attempt{repo: r, key: key}
}

func (a *Attempt) Commit() {}
func (a *Attempt) Rollback() {}

func (r *Repository) Records() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.records...)
}
