package storetoplist

import (
	"sync"
)

type storedList struct {
	mu     sync.Mutex
	Stored *[]Top10List
}

type Storage interface {
	Update(res *[]Top10List)
	Get() []Top10List
}

func NewInMemoryListStore() Storage {
	return &storedList{
		mu:     sync.Mutex{},
		Stored: new([]Top10List),
	}
}

func (r *storedList) Update(res *[]Top10List) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Stored = res
}

func (r *storedList) Get() []Top10List {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *r.Stored
}
