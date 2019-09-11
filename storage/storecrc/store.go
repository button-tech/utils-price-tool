package storecrc

import "sync"

type Res struct {
	Res map[string]Cr
}

type Cr struct {
	USD Currency `json:"USD"`
	EUR Currency `json:"EUR"`
	RUB Currency `json:"RUB"`
}

type Currency struct {
	FROMSYMBOL      string  `json:"FROMSYMBOL"`
	PRICE           float64 `json:"PRICE"`
	CHANGE24HOUR    float64 `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	CHANGEPCTDAY    float64 `json:"CHANGEPCTDAY"`
	CHANGEPCTHOUR   float64 `json:"CHANGEPCTHOUR"`
}

//
//type Result struct {
//	CryptoCurr string
//	Curr       Currencies
//}
//
//type Currencies struct {
//	USD float64 `json:"USD"`
//	EUR float64 `json:"EUR"`
//	RUB float64 `json:"RUB"`
//}

type storedList struct {
	mu     *sync.Mutex
	Stored map[string]Cr
}

type Storage interface {
	Update(res map[string]Cr)
	Get() map[string]Cr
}

func NewInMemoryCRCStore() Storage {
	return &storedList{
		mu:     new(sync.Mutex),
		Stored: make(map[string]Cr),
	}
}

func (r *storedList) Update(res map[string]Cr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Stored = res
}

func (r *storedList) Get() map[string]Cr {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Stored
}
