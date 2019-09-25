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

//var currencies = []string{"USD", "EUR", "RUB"}

var currencies = []string{
	"AED",
	"AFN",
	"ALL",
	"AMD",
	"ANG",
	"AOA",
	"ARS",
	"AUD",
	"AWG",
	"AZN",
	"BAM",
	"BBD",
	"BDT",
	"BGN",
	"BHD",
	"BIF",
	"BMD",
	"BND",
	"BOB",
	"BRL",
	"BSD",
	"BTC",
	"BTN",
	"BWP",
	"BYN",
	"BYR",
	"BZD",
	"CAD",
	"CDF",
	"CHF",
	"CLF",
	"CLP",
	"CNY",
	"COP",
	"CRC",
	"CUC",
	"CUP",
	"CVE",
	"CZK",
	"DJF",
	"DKK",
	"DOP",
	"DZD",
	"EGP",
	"ERN",
	"ETB",
	"EUR",
	"FJD",
	"FKP",
	"GBP",
	"GEL",
	"GGP",
	"GHS",
	"GIP",
	"GMD",
	"GNF",
	"GTQ",
	"GYD",
	"HKD",
	"HNL",
	"HRK",
	"HTG",
	"HUF",
	"IDR",
	"ILS",
	"IMP",
	"INR",
	"IQD",
	"IRR",
	"ISK",
	"JEP",
	"JMD",
	"JOD",
	"JPY",
	"KES",
	"KGS",
	"KHR",
	"KMF",
	"KPW",
	"KRW",
	"KWD",
	"KYD",
	"KZT",
	"LAK",
	"LBP",
	"LKR",
	"LRD",
	"LSL",
	"LTL",
	"LVL",
	"LYD",
	"MAD",
	"MDL",
	"MGA",
	"MKD",
	"MMK",
	"MNT",
	"MOP",
	"MRO",
	"MUR",
	"MVR",
	"MWK",
	"MXN",
	"MYR",
	"MZN",
	"NAD",
	"NGN",
	"NIO",
	"NOK",
	"NPR",
	"NZD",
	"OMR",
	"PAB",
	"PEN",
	"PGK",
	"PHP",
	"PKR",
	"PLN",
	"PYG",
	"QAR",
	"RON",
	"RUB",
	"RWF",
	"SAR",
	"SBD",
	"SCR",
	"SDG",
	"SEK",
	"SGD",
	"SHP",
	"SLL",
	"SOS",
	"SRD",
	"STD",
	"SVC",
	"SYP",
	"SZL",
	"THB",
	"TJS",
	"TMT",
	"TND",
	"TOP",
	"TRY",
	"TTD",
	"TWD",
	"TZS",
	"UAH",
	"UGX",
	"USD",
	"UYU",
	"UZS",
	"VEF",
	"VND",
	"VUV",
	"WST",
	"XAF",
	"XAG",
	"XAU",
	"XCD",
	"XDR",
	"XOF",
	"XPF",
	"YER",
	"ZAR",
	"ZMK",
	"ZMW",
	"ZWL",
}

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
	GetPricesCMC(tokens TokensWithCurrency) (storage.FiatMap, error)
	GetPricesCRC() (storage.FiatMap, error)
	GetTopList(c map[string]string) (map[string]string, error)
}

type service struct{}

func New() Service {
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

	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=100&convert=USD"

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

	topListMap := make(map[string]string)
	for _, item := range list.Data {
		if val, ok := c[item.Symbol]; ok {
			topListMap[val] = item.Symbol
		}
	}

	return topListMap, nil
}

type maps struct {
	FiatMap  storage.FiatMap
	PriceMap map[storage.CryptoCurrency]*storage.Details
}

// Make maps for storage
func storeMapsConstructor() maps {
	return maps{
		FiatMap:  make(storage.FiatMap),
		PriceMap: make(map[storage.CryptoCurrency]*storage.Details),
	}
}

// Get prices from trust-wallet
func (s *service) GetPricesCMC(tokens TokensWithCurrency) (storage.FiatMap, error) {

	url := os.Getenv("TRUST_URL")

	rq, err := req.Post(url, req.BodyJSON(tokens))
	if err != nil {
		return nil, fmt.Errorf("can not make a request: %v", err)
	}

	gotPrices := GotPrices{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return nil, fmt.Errorf("can not marshal: %v", err)
	}

	maps := storeMapsConstructor()
	for _, v := range gotPrices.Docs {
		details := storage.Details{}

		details.Price = v.Price
		details.ChangePCT24Hour = v.PercentChange24H

		maps.PriceMap[storage.CryptoCurrency(strings.ToLower(v.Contract))] = &details
	}
	maps.FiatMap[storage.Fiat(gotPrices.Currency)] = maps.PriceMap

	return maps.FiatMap, nil
}

// Get prices from crypt-compare
func (s *service) GetPricesCRC() (storage.FiatMap, error) {

	var fiats string
	for _, v := range currencies {
		fiats += v + ","
	}

	url := "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=" + fiats

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

	fiatMap := make(storage.FiatMap)
	for k, v := range m {
		priceMap := make(map[storage.CryptoCurrency]*storage.Details)

		for _, i := range v {
			details := storage.Details{}
			details.Price = strconv.FormatFloat(1/i.PRICE, 'f', -1, 64)
			details.ChangePCT24Hour = strconv.FormatFloat(i.CHANGEPCT24HOUR, 'f', 2, 64)
			details.ChangePCTHour = strconv.FormatFloat(i.CHANGEPCTHOUR, 'f', 2, 64)

			priceMap[storage.CryptoCurrency(strings.ToLower(i.TOSYMBOL))] = &details
		}

		fiatMap[storage.Fiat(k)] = priceMap
	}

	return fiatMap, nil
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
