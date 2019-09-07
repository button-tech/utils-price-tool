package storage

import (
	"sync"
	"time"
)


type Top10List struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
	} `json:"status"`
	Data []struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         int         `json:"max_supply"`
		CirculatingSupply int         `json:"circulating_supply"`
		TotalSupply       int         `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             struct {
			USD struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

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
