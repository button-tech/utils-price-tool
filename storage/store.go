package storage

import (
	"sync"
)


type CryptoCurrencies struct {
	BTC   FiatCurrency `json:"BTC"`   //
	ETH   FiatCurrency `json:"ETH"`   //
	XRP   FiatCurrency `json:"XRP"`   //
	BCH   FiatCurrency `json:"BCH"`   //
	LTC   FiatCurrency `json:"LTC"`   //
	BNB   FiatCurrency `json:"BNB"`   //
	WAVES FiatCurrency `json:"WAVES"` //
	XLM   FiatCurrency `json:"XLM"`   //
	EOS   FiatCurrency `json:"EOS"`   //
	ETC   FiatCurrency `json:"ETC"`   //
}

type FiatCurrency struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	RUB float64 `json:"RUB"`
}

type storedPrices struct {
	mu     sync.Mutex
	Stored *[]GotPrices
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
	Update(res *[]GotPrices)
	Get() []GotPrices
}

func NewInMemoryStore() Storage {
	return &storedPrices{
		mu:     sync.Mutex{},
		Stored: new([]GotPrices),
	}
}

func (r *storedPrices) Update(res *[]GotPrices) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Stored = res
}

func (r *storedPrices) Get() []GotPrices {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *r.Stored
}
