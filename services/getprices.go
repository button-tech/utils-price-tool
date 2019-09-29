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
	"sync"
)

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

type Service interface {
	GetPricesCMC(tokens TokensWithCurrency) (storage.FiatMap, error)
	GetPricesCRC(list map[string]string) storage.FiatMap
	GetTopList(c map[string]string) (map[string]string, error)
}

type service struct{}

func New() Service {
	return &service{}
}

// Create from hard-code tokens request data
func CreateCMCRequestData(list map[string]string) TokensWithCurrencies {
	tokensMultiCurrencies := TokensWithCurrencies{}
	tokensOneCurrency := TokensWithCurrency{}
	tokens := make([]Token, 0)

	for _, c := range list {
		token := Token{}
		token.Contract = c
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

	list := topList{}
	if err = rq.ToJSON(&list); err != nil {
		return nil, fmt.Errorf("can not marshal: %v", err)
	}

	if list.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("bad request: %v", list.Status.ErrorCode)
	}

	topListMap := make(map[string]string)
	for _, item := range list.Data {
		if val, ok := c[item.Symbol]; ok {
			topListMap[item.Symbol] = val
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

	gotPrices := gotPrices{}
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
//func (s *service) GetPricesCRC(list map[string]string) (storage.FiatMap, error) {
//
//	var forParams string
//	for k := range list {
//		forParams += k + ","
//	}
//
//	url := "https://min-api.cryptocompare.com/data/pricemultifull"
//	rq, err := req.Get(url, req.Param{
//		"fsyms": forParams,
//		"tsyms": "USD,EUR,RUB",
//	})
//	if err != nil {
//		return nil, fmt.Errorf("can not make req: %v", err)
//	}
//
//	byteRq := rq.Bytes()
//	m, err := crcFastJson(byteRq, list)
//	if err != nil {
//		return nil, fmt.Errorf("can not do fastJson: %v", err)
//	}
//
//	fiatMap := make(storage.FiatMap)
//	for k, v := range m {
//		priceMap := make(map[storage.CryptoCurrency]*storage.Details)
//
//		for _, i := range v {
//			details := storage.Details{}
//			details.Price = strconv.FormatFloat(i.Price, 'f', -1, 64)
//			details.ChangePCT24Hour = strconv.FormatFloat(i.ChangePCT24Hour, 'f', 2, 64)
//			details.ChangePCTHour = strconv.FormatFloat(i.ChangePCTHour, 'f', 2, 64)
//
//			priceMap[storage.CryptoCurrency(strings.ToLower(i.FromSymbol))] = &details
//		}
//
//		fiatMap[storage.Fiat(k)] = priceMap
//	}
//	return fiatMap, nil
//}

func CreateCRCRequestData() []string {
	sortedCurrencies := make([]string, 0)

	n := 0
	step := 25
	for i := 0; i < 6; i++ {
		c := strings.Join(currencies[n:step], ",")
		sortedCurrencies = append(sortedCurrencies, c)
		n += 25
		step += 25
	}
	c := strings.Join(currencies[150:], ",")
	sortedCurrencies = append(sortedCurrencies, c)

	return sortedCurrencies
}

func (s *service) GetPricesCRC(list map[string]string) storage.FiatMap {
	var fsyms string
	for k := range list {
		fsyms += k + ","
	}

	sortedCurrencies := CreateCRCRequestData()
	c := make(chan map[string][]*currency, len(sortedCurrencies))

	wg := sync.WaitGroup{}
	for _, tsyms := range sortedCurrencies {
		wg.Add(1)
		go crcPricesRequests(tsyms, fsyms, list, c, &wg)
	}
	wg.Wait()
	close(c)

	return fiatMapping(c)
}

func crcFastJson(byteRq []byte, list map[string]string) (map[string][]*currency, error) {
	var p fastjson.Parser
	parsed, err := p.ParseBytes(byteRq)
	if err != nil {
		return nil, fmt.Errorf("can not parseBytes: %v", err)
	}

	m := make(map[string][]*currency)

	o := parsed.GetObject("RAW")
	o.Visit(func(k []byte, v *fastjson.Value) {
		if val, ok := list[string(k)]; ok {
			crypto := v.GetObject()
			crypto.Visit(func(key []byte, value *fastjson.Value) {

				c := currency{}
				if err := json.Unmarshal([]byte(value.String()), &c); err != nil {
					log.Printf("can not unmarshal elem: %v", value.String())
					return
				}

				c.FromSymbol = val

				valM, okM := m[c.ToSymbol]
				if !okM {
					m[c.ToSymbol] = make([]*currency, 0)
				}
				valM = append(valM, &c)
				m[c.ToSymbol] = valM
			})
		}
	})
	return m, nil
}

func crcPricesRequests(tsyms, fsyms string, list map[string]string, c chan <- map[string][]*currency, wg *sync.WaitGroup)  {
	url := "https://min-api.cryptocompare.com/data/pricemultifull"
	rq, err := req.Get(url, req.Param{
		"fsyms": fsyms,
		"tsyms": tsyms,
	})
	if err != nil {
		log.Printf("can not make req: %v", err)
	}

	byteRq := rq.Bytes()
	m, err := crcFastJson(byteRq, list)
	if err != nil {
		log.Printf("can not do fastJson: %v", err)
	}

	c <- m
	wg.Done()
}

func fiatMapping(c <- chan map[string][]*currency) storage.FiatMap {
	fiatMap := make(storage.FiatMap)

	for range c {
		for k, v := range <- c {
			priceMap := make(map[storage.CryptoCurrency]*storage.Details)

			for _, i := range v {
				details := storage.Details{}
				details.Price = strconv.FormatFloat(i.Price, 'f', -1, 64)
				details.ChangePCT24Hour = strconv.FormatFloat(i.ChangePCT24Hour, 'f', 2, 64)
				details.ChangePCTHour = strconv.FormatFloat(i.ChangePCTHour, 'f', 2, 64)

				priceMap[storage.CryptoCurrency(strings.ToLower(i.FromSymbol))] = &details
			}
			if _, ok := fiatMap[storage.Fiat(k)]; !ok {
				fiatMap[storage.Fiat(k)] = map[storage.CryptoCurrency]*storage.Details{}
			}
			fiatMap[storage.Fiat(k)] = priceMap
		}
	}
	return fiatMap
}
