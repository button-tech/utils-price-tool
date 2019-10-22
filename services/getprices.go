package services

import (
	"encoding/json"
	"fmt"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/imroc/req"
	"github.com/pkg/errors"
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

const (
	urlHuobi   = "https://api.hbdm.com/api/v1/contract_index"
	urlTopList = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=100&convert=USD"
	urlCRC     = "https://min-api.cryptocompare.com/data/pricemultifull"
)

var (
	urlTrustWallet = os.Getenv("TRUST_URL")
	topListAPIKey  = os.Getenv("API_KEY")
)

type Service struct {
	list map[string]string
}

func New() *Service {
	return &Service{
		list: make(map[string]string),
	}
}

func (s *Service) CreateCMCRequestData() RequestCoinMarketCap {
	tokensMultiCurrencies := RequestCoinMarketCap{}
	tokensOneCurrency := TokensWithCurrency{}
	tokens := make([]Token, 0)

	for _, c := range s.list {
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
func (s *Service) GetTopList(c map[string]string) error {
	rq, err := req.Get(urlTopList, req.Header{"X-CMC_PRO_API_KEY": topListAPIKey})
	if err != nil {
		return fmt.Errorf("can not make a request: %v", err)
	}

	list := topList{}
	if err = rq.ToJSON(&list); err != nil {
		return fmt.Errorf("can not marshal: %v", err)
	}

	if list.Status.ErrorCode != 0 {
		return fmt.Errorf("bad request: %v", list.Status.ErrorCode)
	}

	topListMap := make(map[string]string)
	for _, item := range list.Data {
		if val, ok := c[item.Symbol]; ok {
			topListMap[item.Symbol] = val
		}
	}

	s.list = topListMap
	return nil
}

type maps struct {
	FiatMap  storage.FiatMap
	PriceMap map[storage.CryptoCurrency]*storage.Details
}

func storeMapsConstructor() maps {
	return maps{
		FiatMap:  make(storage.FiatMap),
		PriceMap: make(map[storage.CryptoCurrency]*storage.Details),
	}
}

func (s *Service) GetPricesCMC(tokens TokensWithCurrency) (storage.FiatMap, error) {
	rq, err := req.Post(urlTrustWallet, req.BodyJSON(tokens))
	if err != nil {
		return nil, fmt.Errorf("can not make a request: %v", err)
	}

	gotPrices := coinMarketCap{}
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

func (s *Service) GetPricesCRC() storage.FiatMap {
	var fsyms string
	for k := range s.list {
		fsyms += k + ","
	}

	sortedCurrencies := CreateCRCRequestData()
	c := make(chan map[string][]*cryptoCompare, len(sortedCurrencies))

	wg := sync.WaitGroup{}
	for _, tsyms := range sortedCurrencies {
		wg.Add(1)
		go crcPricesRequest(tsyms, fsyms, s.list, c, &wg)
	}
	wg.Wait()
	close(c)
	return fiatMapping(c)
}

func crcFastJson(byteRq []byte, list map[string]string) (map[string][]*cryptoCompare, error) {
	var p fastjson.Parser
	parsed, err := p.ParseBytes(byteRq)
	if err != nil {
		return nil, fmt.Errorf("can not parseBytes: %v", err)
	}

	m := make(map[string][]*cryptoCompare)

	o := parsed.GetObject("RAW")
	o.Visit(func(k []byte, v *fastjson.Value) {
		if val, ok := list[string(k)]; ok {
			crypto := v.GetObject()
			crypto.Visit(func(key []byte, value *fastjson.Value) {

				c := cryptoCompare{}
				if err := json.Unmarshal([]byte(value.String()), &c); err != nil {
					log.Printf("can not unmarshal elem: %v", value.String())
					return
				}

				c.FromSymbol = val

				valM, okM := m[c.ToSymbol]
				if !okM {
					m[c.ToSymbol] = make([]*cryptoCompare, 0)
				}
				valM = append(valM, &c)
				m[c.ToSymbol] = valM
			})
		}
	})
	return m, nil
}

func crcPricesRequest(tsyms, fsyms string, list map[string]string, c chan<- map[string][]*cryptoCompare, wg *sync.WaitGroup) {
	rq, err := req.Get(urlCRC, req.Param{
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

func fiatMapping(c chan map[string][]*cryptoCompare) storage.FiatMap {
	fiatMap := make(storage.FiatMap)

	done := false
	for !done {
		select {
		case m, ok := <-c:
			if !ok {
				done = true
				break
			}
			for k, v := range m {
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
	}

	return fiatMap
}

func (s *Service) GetPricesHUOBI() storage.FiatMap {
	rq, err := req.Get(urlHuobi)
	if err != nil {
		log.Println(errors.Wrap(err, "huobi"))
		return nil
	}

	var h huobi
	if err := rq.ToJSON(&h); err != nil {
		log.Println(errors.Wrap(err, "toJSON huobi"))
		return nil
	}
	return huobiMapping(&h, s.list)
}

func huobiMapping(h *huobi, list map[string]string) storage.FiatMap {
	fiatMap := make(storage.FiatMap)
	priceMap := make(map[storage.CryptoCurrency]*storage.Details)

	for _, i := range h.Data {
		if val, ok := list[i.Symbol]; ok {
			var details storage.Details
			details.Price = strconv.FormatFloat(i.IndexPrice, 'f', -1, 64)
			priceMap[storage.CryptoCurrency(strings.ToLower(val))] = &details
			fiatMap[storage.Fiat("USD")] = priceMap
		}
	}
	return fiatMap
}
