package storecrc

import "sync"

//type Fiat struct {
//	USD Currency `json:"USD"`
//	EUR Currency `json:"EUR"`
//	RUB Currency `json:"RUB"`
//}

type Currency struct {
	TOSYMBOL 	      string `json:"TOSYMBOL"`
	FROMSYMBOL        string  `json:"FROMSYMBOL"`
	PRICE           float64 `json:"PRICE"`
	CHANGEPCT24HOUR float64 `json:"CHANGEPCT24HOUR"`
	CHANGEPCTHOUR   float64 `json:"CHANGEPCTHOUR"`
}

type storedList struct {
	mu     *sync.Mutex
	Stored map[string][]*Currency
}

type Storage interface {
	Update(key string, value []*Currency)
	Get() map[string][]*Currency
}

func NewInMemoryCRCStore() Storage {
	return &storedList{
		mu:     new(sync.Mutex),
		Stored: make(map[string][]*Currency, 0),
	}
}

func (r *storedList) Update(key string, value []*Currency) {
	r.mu.Lock()
	r.Stored[key] = value
	r.mu.Unlock()
}

func (r *storedList) Get() map[string][]*Currency {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Stored
}
