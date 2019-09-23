package services

import (
	"encoding/json"
	"fmt"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/imroc/req"
	"github.com/valyala/fastjson"
	"log"
	"os"
	"strconv"
	"strings"
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

type Service interface {
	GetPricesCMC(tokens TokensWithCurrency) (map[storage.Fiat]map[storage.CryptoCurrency]*storage.Details, error)
	GetPricesCRC() (map[storage.Fiat]map[storage.CryptoCurrency]*storage.Details, error)
	GetTopList(c map[string]string) (map[string]string, error)
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

// Get top list of crypto-currencies from coin-market
func (s *service) GetTopList(c map[string]string) (map[string]string, error) {

	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=10&convert=USD"

	rq, err := req.Get(url, req.Header{"X-CMC_PRO_API_KEY": os.Getenv("API_KEY")})
	if err != nil {
		return nil, fmt.Errorf("can not make a request: %v", err)
	}

	list := TopList{}
	if err = rq.ToJSON(&list); err != nil {
		return nil, fmt.Errorf("can not marshal: %v", err)
	}

	if list.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("bad request: %v", list.Status.ErrorCode)
	}

	topListMap := make(map[string]string, 0)
	for _, item := range list.Data {
		if val, ok := c[item.Symbol]; ok {
			topListMap[val] = item.Symbol
		}
	}

	return topListMap, nil
}

type maps struct {
	FiatMap map[storage.Fiat]map[storage.CryptoCurrency]*storage.Details
	PriceMap map[storage.CryptoCurrency]*storage.Details
}

// Make maps for storage
func storeConstructor() maps {
	return maps{
		FiatMap:  make(map[storage.Fiat]map[storage.CryptoCurrency]*storage.Details),
		PriceMap: make(map[storage.CryptoCurrency]*storage.Details),
	}
}

// Get prices from trust-wallet
func (s *service) GetPricesCMC(tokens TokensWithCurrency) (map[storage.Fiat]map[storage.CryptoCurrency]*storage.Details, error) {

	url := os.Getenv("TRUST_URL")

	rq, err := req.Post(url, req.BodyJSON(tokens))
	if err != nil {
		return nil, fmt.Errorf("can not make a request: %v", err)
	}

	gotPrices := GotPrices{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return nil, fmt.Errorf("can not marshal: %v", err)
	}

	maps := storeConstructor()
	for _, v := range gotPrices.Docs {
		details := storage.Details{}

		details.Price = v.Price
		details.ChangePCT24Hour = v.PercentChange24H

		maps.PriceMap[storage.CryptoCurrency(strings.ToUpper(v.Contract))] = &details
	}
	maps.FiatMap[storage.Fiat(gotPrices.Currency)] = maps.PriceMap

	return maps.FiatMap, nil
}

// Get prices from crypt-compare
func(s *service) GetPricesCRC() (map[storage.Fiat]map[storage.CryptoCurrency]*storage.Details, error) {

	url := "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=USD,EUR,RUB"

	var forParams string
	for _, k := range convertedCurrencies {
		forParams += k + ","
	}

	rq, err := req.Get(url, req.Param{"tsyms": forParams})
	if err != nil {
		return nil, fmt.Errorf("can not make req: %v", err)
	}

	byteRq := rq.Bytes()
	m, err := crcFastJson(byteRq)
	if err != nil {
		return nil, fmt.Errorf("can not do fastJson: %v", err)
	}

	maps := storeConstructor()
	for k, v := range m {
		for _, i := range v {
			details := storage.Details{}

			details.Price = strconv.FormatFloat(1/i.PRICE, 'f', 2, 64)
			details.ChangePCT24Hour = strconv.FormatFloat(i.CHANGEPCT24HOUR, 'f', 2, 64)
			details.ChangePCTHour = strconv.FormatFloat(i.CHANGEPCTHOUR, 'f', 2, 64)

			maps.PriceMap[storage.CryptoCurrency(i.TOSYMBOL)] = &details
		}
		maps.FiatMap[storage.Fiat(k)] = maps.PriceMap
	}

	return maps.FiatMap, nil
}

func crcFastJson(byteRq []byte) (map[string][]*Currency, error) {
	m := make(map[string][]*Currency)

	var p fastjson.Parser
	parsed, err := p.ParseBytes(byteRq)
	if err != nil {
		return nil, fmt.Errorf("can not parseBytes: %v", err)
	}

	o := parsed.GetObject("RAW")
	o.Visit(func(k []byte, v *fastjson.Value) {

		currencies := make([]*Currency, 0)

		fiats := v.GetObject()
		fiats.Visit(func(key []byte, value *fastjson.Value) {
			currency := Currency{}

			if err := json.Unmarshal([]byte(value.String()), &currency); err != nil {
				log.Printf("can not unmarshal elem: %v", value.String())
				return
			}

			for t, c := range convertedCurrencies {
				if c == currency.TOSYMBOL {
					currency.TOSYMBOL = t
				}
			}
			currencies = append(currencies, &currency)

			m[string(k)] = currencies
		})
	})

	return m, nil
}