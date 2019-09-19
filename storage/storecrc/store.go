package storecrc

import "sync"

type Fiat struct {
	USD Currency `json:"USD"`
	EUR Currency `json:"EUR"`
	RUB Currency `json:"RUB"`
}

type Currency struct {
	TOSYMBOL        string  `json:"TOSYMBOL"`
	PRICE           float64 `json:"PRICE"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	CHANGEPCTHOUR   float64 `json:"CHANGEPCTHOUR"`
}

type storedList struct {
	mu     *sync.Mutex
	Stored map[string]Fiat
}

type Storage interface {
	Update(res map[string]Fiat)
	Get() map[string]Fiat
}

func NewInMemoryCRCStore() Storage {
	return &storedList{
		mu:     new(sync.Mutex),
		Stored: make(map[string]Fiat),
	}
}

func (r *storedList) Update(res map[string]Fiat) {
	r.mu.Lock()
	r.Stored = res
	r.mu.Unlock()
}

func (r *storedList) Get() map[string]Fiat {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Stored
}
