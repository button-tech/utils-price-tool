package services

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/utils-price-tool/storage"
	"os"
)

var currencies = []string{"USD", "EUR", "RUB"}

var convertedCurrencies = map[string]string{
	"0x0000000000000000000000000000000000000000": "BTC",
	"0x000000000000000000000000000000000000003C": "ETH",
	"0x0000000000000000000000000000000000000002": "LTC",
	"0x000000000000000000000000000000000000003D": "ETC",
	"0x0000000000000000000000000000000000000091": "BCH",
	"0x0000000000000000000000000000000000579BFC": "Waves",
	"0x0000000000000000000000000000000000000094": "XLM",
	"0x00000000000000000000000000000000000000C2": "EOS",
	"0x00000000000000000000000000000000000002CA": "BNB",
	"0x0000000000000000000000000000000000000090": "XRP",
}

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
	GetCRCPrices() (storage.CryptoCurrencies, error)
	GetTopList() (storage.Top10List, error)
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

// trust-wallet
func (s *service) GetPricesCMC(tokens *TokensWithCurrency) (storage.GotPrices, error) {
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

// top list from coin-market-cap
func (s *service) GetTopList() (storage.Top10List, error) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=10"
	rq, err := req.Get(url, req.Header{"X-CMC_PRO_API_KEY": os.Getenv("API_KEY")})
	if err != nil {
		return storage.Top10List{}, fmt.Errorf("can not make a request: %v", err)
	}

	list := storage.Top10List{}
	if err = rq.ToJSON(&list); err != nil {
		return storage.Top10List{}, fmt.Errorf("can not marshal: %v", err)
	}

	return list, nil
}

// get prices from crypto-compare
func (s *service) GetCRCPrices() (storage.CryptoCurrencies, error) {
	url := "https://min-api.cryptocompare.com/data/pricemulti?tsyms=USD,EUR,RUB"
	forParams := ""
	for _, k := range convertedCurrencies {
		forParams += k + ","
	}

	rq, err := req.Get(url, req.Param{"fsyms": forParams})
	if err != nil {
		return storage.CryptoCurrencies{}, fmt.Errorf("can not make a request: %v", err)
	}

	crypto := storage.CryptoCurrencies{}
	if err = rq.ToJSON(&crypto); err != nil {
		return storage.CryptoCurrencies{}, fmt.Errorf("can not marshal: %v", err)
	}

	return crypto, nil
}

func converter(s storage.Top10List) {

}
