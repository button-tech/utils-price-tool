package storage

import "sync"


type Prices struct {
	Currency string               `json:"currency"`
	Rates    []*map[string]string `json:"rates"`
}

type ResultPrices struct {
	mu     sync.Mutex
	Prices []Prices
}

type Storage interface {
	Update(res *[]Prices)
	Get() *[]Prices
}

//func NewInMemoryStore() Storage {
//	return &resultPrices{
//		mu:     new(sync.Mutex),
//		prices: new([]Prices),
//	}
//}

func (r *ResultPrices) Update(res Prices) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Prices = append(r.Prices, res)
}

//func (r *ResultPrices) Get() []*Prices {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//	return r.prices
//}
