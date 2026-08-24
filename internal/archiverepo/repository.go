package archiverepo

import "sync"

type Repository struct {
	mu      sync.Mutex
	records []string
}

// New returns an empty Repository.
func New() *Repository {
	return &Repository{}
}

type Attempt struct {
	repo    *Repository
	key     string
	commit  bool
	rolled  bool
}

func (r *Repository) Begin(key string) *Attempt {
	return &Attempt{repo: r, key: key}
}

// Commit finalizes the attempt and persists a record. Only committed
// attempts leave a trace, so a temporary failure that succeeds on retry
// yields exactly one record regardless of how many attempts ran.
func (a *Attempt) Commit() {
	if a == nil || a.commit || a.rolled {
		return
	}
	a.commit = true
	a.repo.append(a.key)
}

// Rollback discards the attempt. The matching no-op-on-nil guard keeps it
// safe to call from deferred cleanup regardless of how Begin fared.
func (a *Attempt) Rollback() {
	if a == nil || a.commit || a.rolled {
		return
	}
	a.rolled = true
}

func (r *Repository) append(key string) {
	r.mu.Lock()
	r.records = append(r.records, key)
	r.mu.Unlock()
}

func (r *Repository) Records() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.records...)
}
