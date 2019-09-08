package storecrc

import "sync"

type Result struct {
	CryptoCurr string
	Curr Currencies
}

type Currencies struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	RUB float64 `json:"RUB"`
}

type storedList struct {
	mu     sync.Mutex
	Stored *[]Result
}

type Storage interface {
	Update(res *[]Result)
	Get() []Result
}

func NewInMemoryCRCStore() Storage {
	return &storedList{
		mu:     sync.Mutex{},
		Stored: new([]Result),
	}
}

func (r *storedList) Update(res *[]Result) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Stored = res
}

func (r *storedList) Get() []Result {
	r.mu.Lock()
	defer r.mu.Unlock()
	return *r.Stored
}

