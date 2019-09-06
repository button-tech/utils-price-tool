package services

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/utils-price-tool/storage"
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

type Service interface {
	GetPricesCMC(tokens *TokensWithCurrency) (storage.GotPrices, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

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

func(s *service) GetPricesCMC(tokens *TokensWithCurrency) (storage.GotPrices, error) {
	url := os.Getenv("TRUST_URL")
	rq, err := req.Post(url, req.BodyJSON(tokens))
	if err != nil {
		return storage.GotPrices{}, fmt.Errorf("can not make a request: %v", err)
	}

	gotPrices := storage.GotPrices{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return storage.GotPrices{}, fmt.Errorf("can not marshal: %v", err)
	}

	return gotPrices, nil
}
