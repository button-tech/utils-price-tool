package storecrc

import "sync"

type Result struct {
	CryptoCurr string
}

type Currencies struct {
	USD string `json:"USD"`
	EUR string `json:"EUR"`
	RUB string `json:"RUB"`
}

type storedList struct {
	mu     sync.Mutex
	Stored *[]map[string]Currencies
}

type Storage interface {
	Update(res *[]map[string]Currencies)
	Get() []map[string]Currencies
}

func NewInMemoryCRCStore() Storage {
	return &storedList{
		mu:     sync.Mutex{},
		Stored: new([]map[string]Currencies),
	}
}

func (r *storedList) Update(res *[]map[string]Currencies) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Stored = res
}

func (r *storedList) Get() []map[string]Currencies {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *r.Stored
}

