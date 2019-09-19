package storetoplist

import (
	"sync"
	"time"
)

type TopList struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
	} `json:"status"`
	Data []Data `json:"data"`
}

type Data struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
}

type storedList struct {
	mu     *sync.Mutex
	Stored *TopList
}

type Storage interface {
	Update(res *TopList)
	Get() TopList
}

func NewInMemoryListStore() Storage {
	return &storedList{
		mu:     new(sync.Mutex),
		Stored: new(TopList),
	}
}

func (r *storedList) Update(res *TopList) {
	r.mu.Lock()
	r.Stored = res
	r.mu.Unlock()
}

func (r *storedList) Get() TopList {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *r.Stored
}
