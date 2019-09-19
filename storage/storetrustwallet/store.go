package storetrustwallet

import (
	"sync"
)

type storedPrices struct {
	mu     *sync.Mutex
	Stored []*GotPrices
}

type GotPrices struct {
	Status   bool         `json:"status"`
	Docs     []DocsPrices `json:"docs"`
	Currency string       `json:"currency"`
}

type DocsPrices struct {
	Price            string `json:"price"`
	Contract         string `json:"contract"`
	PercentChange24H string `json:"percent_change_24h"`
}

type Storage interface {
	Update(res []*GotPrices)
	Get() []*GotPrices
}

func NewInMemoryCMCStore() Storage {
	return &storedPrices{
		mu:     new(sync.Mutex),
		Stored: make([]*GotPrices, 0),
	}
}

func (r *storedPrices) Update(res []*GotPrices) {
	r.mu.Lock()
	r.Stored = res
	r.mu.Unlock()
}

func (r *storedPrices) Get() []*GotPrices {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Stored
}
