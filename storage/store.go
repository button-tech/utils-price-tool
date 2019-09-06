package storage

import "sync"


type Prices struct {
	mu sync.Mutex
	Currency string               `json:"currency"`
	Rates    []*map[string]string `json:"rates"`
}

type resultPrices struct {
	mu     sync.Mutex
	Prices *[]Prices
}

type Storage interface {
	Update(res *[]Prices)
	Get() *[]Prices
}

func NewInMemoryStore() Storage {
	return &resultPrices{
		mu:     sync.Mutex{},
		Prices: new([]Prices),
	}
}

func(r *resultPrices) Update(res *[]Prices) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Prices = res
}

func(r *resultPrices) Get() *[]Prices {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Prices
}

//func(st *Prices)
