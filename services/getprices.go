package services

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/utils-tool_prices/storage"
	"log"
	"os"
)

var convertedCurrencies = map[string]string{
	"0x0000000000000000000000000000000000000000": "BTC",
	"0x000000000000000000000000000000000000003C": "ETH",
	"0x0000000000000000000000000000000000000002": "LTC",
	"0x000000000000000000000000000000000000003D": "ETC",
	"0x0000000000000000000000000000000000000091": "BCH",
	"0x0000000000000000000000000000000000579BFC": "Waves",
	"0x0000000000000000000000000000000000000094": "XLM",
}

var currencies = []string{"USD", "EUR", "RUB"}

type TokensWithCurrency struct {
	Currency string  `json:"currency"`
	Tokens   []Token `json:"tokens"`
}

type Token struct {
	Contract string `json:"contract"`
}

// data to get prices to trust-wallet
type TokensWithCurrencies struct {
	Tokens []TokensWithCurrency
}

// data from trust-wallet
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

type Service interface {
	GetAllPricesCMC() *[]storage.Prices
}

type service struct{}

//func NewService() Service {
//	return &service{}
//}

func InitRequestData() TokensWithCurrencies {
	tokensMultiCurrencies := TokensWithCurrencies{}
	var tokens []Token
	tokensOneCurrency := TokensWithCurrency{}

	for k, _ := range convertedCurrencies {
		token := Token{}
		token.Contract = k
		tokens = append(tokens, token)
	}

	tokensOneCurrency.Tokens = tokens

	for _, c := range currencies {
		tokensOneCurrency.Currency = c
		tokensMultiCurrencies.Tokens = append(tokensMultiCurrencies.Tokens, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}

func(s *service) GetAllPricesCMC() *storage.Prices {
	tokens := InitRequestData()
	got, err := GetPricesCMC(&tokens.Tokens[0])
	if err != nil {
		log.Println(err)
		return nil
	}
	var store storage.Prices
	for _, i := range got.Docs {
		store.Currency = got.Currency
		contract := map[string]string{i.Contract: i.Price}
		store.Rates = append(store.Rates, &contract)
	}


	return &store
}

//func (s *service) GetAllPricesCMC() *[]storage.Prices {
//	tokens := initRequestData()
//	fmt.Println(len(tokens.Tokens))
//
//	var stored []storage.Prices
//
//	wg := sync.WaitGroup{}
//
//	for _, t := range tokens.Tokens {
//		wg.Add(1)
//
//		go func(wg *sync.WaitGroup) {
//			defer wg.Done()
//
//			var store storage.Prices
//			got, err := GetPricesCMC(&t)
//			if err != nil {
//				log.Println(err)
//			}
//
//			for _, i := range got.Docs {
//				store.Currency = got.Currency
//				contract := map[string]string{i.Contract: i.Price}
//				store.Rates = append(store.Rates, &contract)
//			}
//
//			stored = append(stored, store)
//		}(&wg)
//
//	}
//	wg.Wait()
//
//	return &stored
//}

func GetPricesCMC(tokens *TokensWithCurrency) (GotPrices, error) {
	url := os.Getenv("TRUST_URL")
	rq, err := req.Post(url, req.BodyJSON(tokens))
	if err != nil {
		return GotPrices{}, fmt.Errorf("can not make a request: %v", err)
	}

	gotPrices := GotPrices{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return GotPrices{}, fmt.Errorf("can not marshal: %v", err)
	}

	return gotPrices, nil
}
