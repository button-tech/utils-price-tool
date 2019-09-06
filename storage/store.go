package storage

import "sync"

type ResultPrices struct {
	Prices []*Prices
}

type Prices struct {
	Currency string               `json:"currency"`
	Rates    []*map[string]string `json:"rates"`
}

type resultPrices struct {
	mu     *sync.Mutex
	prices *ResultPrices
}

type Storage interface {
	Update(res *ResultPrices)
	Get() *ResultPrices
}

func NewInMemoryStore() Storage {
	return &resultPrices{
		mu:     new(sync.Mutex),
		prices: new(ResultPrices),
	}
}

func (r *resultPrices) Update(res *ResultPrices) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prices = res
}

func (r *resultPrices) Get() *ResultPrices {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.prices
}
