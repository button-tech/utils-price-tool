package services

import (
	"encoding/json"
	"fmt"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"github.com/button-tech/utils-price-tool/storage/storetrustwallet"
	"github.com/imroc/req"
	"github.com/valyala/fastjson"
	"log"
	"os"
)

var currencies = []string{"USD", "EUR", "RUB"}

var convertedCurrencies = map[string]string{
	"0x0000000000000000000000000000000000000000": "BTC",
	"0x000000000000000000000000000000000000003C": "ETH",
	"0x0000000000000000000000000000000000000002": "LTC",
	"0x000000000000000000000000000000000000003D": "ETC",
	"0x0000000000000000000000000000000000000091": "BCH",
	"0x0000000000000000000000000000000000579BFC": "WAVES",
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

// Data to get prices to trust-wallet
type TokensWithCurrencies struct {
	Tokens []TokensWithCurrency
}

type Service interface {
	GetPricesCMC(tokens *TokensWithCurrency) (storetrustwallet.GotPrices, error)
	GetCRCPrices() (map[string]storecrc.Cr, error)
	GetTopList() (*storetoplist.TopList, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

// Create from hard-code tokens request data
func InitRequestData() TokensWithCurrencies {
	tokensMultiCurrencies := TokensWithCurrencies{}
	tokens := make([]Token, 0)
	tokensOneCurrency := TokensWithCurrency{}

	for k := range convertedCurrencies {
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

// Get prices from trust-wallet
func (s *service) GetPricesCMC(tokens *TokensWithCurrency) (storetrustwallet.GotPrices, error) {
	url := os.Getenv("TRUST_URL")

	rq, err := req.Post(url, req.BodyJSON(tokens))
	if err != nil {
		return storetrustwallet.GotPrices{}, fmt.Errorf("can not make a request: %v", err)
	}

	gotPrices := storetrustwallet.GotPrices{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return storetrustwallet.GotPrices{}, fmt.Errorf("can not marshal: %v", err)
	}

	return gotPrices, nil
}

func (s *service) GetTopList() (*storetoplist.TopList, error) {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=10&convert=USD"

	rq, err := req.Get(url, req.Header{"X-CMC_PRO_API_KEY": os.Getenv("API_KEY")})
	if err != nil {
		return nil, fmt.Errorf("can not make a request: %v", err)
	}

	list := storetoplist.TopList{}
	if err = rq.ToJSON(&list); err != nil {
		return nil, fmt.Errorf("can not marshal: %v", err)
	}

	return &list, nil
}

// Get prices from crypto-compare
func (s *service) GetCRCPrices() (map[string]storecrc.Cr, error) {
	url := "https://min-api.cryptocompare.com/data/pricemultifull?tsyms=USD,EUR,RUB"

	var forParams string
	for _, k := range convertedCurrencies {
		forParams += k + ","
	}

	rq, err := req.Get(url, req.Param{"fsyms": forParams})
	if err != nil {
		return nil, fmt.Errorf("can not make req: %v", err)
	}

	byteRq := rq.Bytes()
	m, err := crcFastJson(byteRq)
	if err != nil {
		return nil, fmt.Errorf("can not do fastJson: %v", err)
	}

	return m, nil
}

func crcFastJson(byteRq []byte) (map[string]storecrc.Cr, error) {
	m := make(map[string]storecrc.Cr)
	var cr storecrc.Cr

	var p fastjson.Parser
	parsed, err := p.ParseBytes(byteRq)
	if err != nil {
		return nil, fmt.Errorf("can not parseBytes: %v", err)
	}

	o := parsed.GetObject("RAW")
	o.Visit(func(k []byte, v *fastjson.Value) {

		if err := json.Unmarshal([]byte(v.String()), &cr); err != nil {
			log.Printf("can not unmarshal elem: %v", v.String())
		}

		for t, c := range convertedCurrencies {
			if c == string(k) {
				m[t] = cr
			}
		}
	})

	return m, nil
}
