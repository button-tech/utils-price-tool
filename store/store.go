package store

import (
	"github.com/imroc/req"
	"log"
	"os"
	"sync"
)



type Storage interface {

}

// data to trustwallet
type dataGetPrices struct {
	Currency string  `json:"currency"`
	Tokens    []Token `json:"tokens"`
}

type resultPrices struct {
	mu   sync.Mutex
	Data *[]Prices `json:"data"`
}

type Prices struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

type Token struct {
	Contract string `json:"contract"`
}

type service struct {

}

func newService() Storage {
	return &service{}
}

func (r *gotPrices) getPrices(d dataGetPrices) {
	url := os.Getenv("TRUST_URL")
	rq, err := req.Post(url, req.BodyJSON(d))
	if err != nil {
		log.Println(err)
		return
	}
}

func (res *resultPrices) Set() {
	res.mu.Lock()

	res.mu.Unlock()
}



