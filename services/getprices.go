package services

import (
	"encoding/json"
	"github.com/button-tech/logger"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/button-tech/utils-price-tool/storage"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

const (
	urlHuobi   = "https://api.hbdm.com/api/v1/contract_index"
	urlTopList = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=100&convert=USD"
	urlCRC     = "https://min-api.cryptocompare.com/data/pricemultifull"
)

const coin = "coin"

var (
	urlTrustWallet   = os.Getenv("TRUST_URL")
	topListAPIKey    = os.Getenv("API_KEY")
	urlTrustWalletV2 = os.Getenv("TRUST_URL_V2")
)

type Service struct {
	TrustV2Coins []PricesTrustV2

	store *storage.Cache
	list  map[string]string
}

type maps struct {
	FiatMap  storage.FiatMap
	PriceMap map[storage.CryptoCurrency]*storage.Details
}

var trustV2Coins = map[string]int{
	"ETH":   60,
	"ETC":   61,
	"ICX":   74,
	"ATOM":  118,
	"XRP":   144,
	"XLM":   148,
	"POA":   178,
	"TRX":   195,
	"FIO":   235,
	"NIM":   242,
	"IOTX":  304,
	"ZIL":   313,
	"AION":  425,
	"AE":    457,
	"THETA": 500,
	"BNB":   714,
	"VET":   818,
	"CLO":   820,
	"TOMO":  889,
	"TT":    1001,
	"ONT":   1024,
	"XTZ":   1729,
	"KIN":   2017,
	"NAS":   2718,
	"GO":    6060,
	"WAN":   5718350,
	"WAVES": 5741564,
	"SEM":   7562605,
	"BTC":   0,
	"LTC":   2,
	"DOGE":  3,
	"DASH":  5,
	"VIA":   14,
	"GRS":   17,
	"ZEC":   133,
	"XZC":   136,
	"BCH":   145,
	"RVN":   175,
	"QTUM":  2301,
	"ZEL":   19167,
	"DCR":   42,
	"ALGO":  283,
	"NANO":  165,
	"DGB":   20,
}

func New(store *storage.Cache) *Service {
	return &Service{
		TrustV2Coins: CreateTrustV2RequestData(),
		list:         make(map[string]string),
		store:        store,
	}
}

func CreateTrustV2RequestData() []PricesTrustV2 {
	prices := make([]PricesTrustV2, 0, len(currencies))
	for _, c := range currencies {
		price := PricesTrustV2{
			Currency: c,
		}
		allAssets := make([]assets, 0, len(trustV2Coins))
		for _, v := range trustV2Coins {
			allAssets = append(allAssets, assets{
				Coin: v,
				Type: coin,
			})
		}
		price.Assets = allAssets
		prices = append(prices, price)
	}
	return prices
}

func (s *Service) CreateCMCRequestData() []TokensWithCurrency {
	var tokensOneCurrency TokensWithCurrency
	tokensMultiCurrencies := make([]TokensWithCurrency, len(currencies))
	tokens := make([]Token, len(s.list))

	for _, c := range s.list {
		tokens = append(tokens, Token{Contract: c})
	}
	tokensOneCurrency.Tokens = tokens

	for _, c := range currencies {
		tokensOneCurrency.Currency = c
		tokensMultiCurrencies = append(tokensMultiCurrencies, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}

// Get top list of crypto-currencies from coin-market
func (s *Service) GetTopList(c map[string]string) error {
	rq, err := req.Get(urlTopList, req.Header{"X-CMC_PRO_API_KEY": topListAPIKey})
	if err != nil {
		return errors.Wrap(err, "getTopList")
	}

	var topList realCoinMarketCap
	if err = rq.ToJSON(&topList); err != nil {
		return errors.Wrap(err, "getTopList")
	}

	if topList.Status.ErrorCode != 0 {
		return errors.New("responseErrorCodeNotOk")
	}

	ms := storeMapsConstructor()
	topListMap := make(map[string]string)
	for _, item := range topList.Data {
		if val, ok := c[item.Symbol]; ok {
			topListMap[item.Symbol] = val

			pricesData := storageDetails(
				item.Quote.USD.Price,
				item.Quote.USD.PercentChange1H,
				item.Quote.USD.PercentChange24H,
				item.Quote.USD.PercentChange7D,
			)
			ms.PriceMap[storage.CryptoCurrency(val)] = &pricesData
		}
	}

	ms.FiatMap[storage.Fiat("USD")] = ms.PriceMap
	s.store.Set("coinMarketCap", ms.FiatMap)
	s.list = topListMap

	return nil
}

func storageDetails(price, hour, hour24, sevenDay float64) storage.Details {
	return storage.Details{
		Price:           strconv.FormatFloat(price, 'f', 10, 64),
		ChangePCTHour:   strconv.FormatFloat(hour, 'f', 6, 64),
		ChangePCT24Hour: strconv.FormatFloat(hour24, 'f', 6, 64),
		ChangePCT7Day:   strconv.FormatFloat(sevenDay, 'f', 6, 64),
	}
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
		return nil, errors.Wrap(err, "GetPricesCMC")
	}

	gotPrices := coinMarketCap{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return nil, errors.Wrap(err, "GetPricesCMC")
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
	sortedCurrencies := make([]string, 7)

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

	var wg sync.WaitGroup
	for _, tsyms := range sortedCurrencies {
		wg.Add(1)
		go s.crcPricesRequest(tsyms, fsyms, c, &wg)
	}
	wg.Wait()
	close(c)

	return fiatMapping(c)
}

func (s *Service) crcFastJson(byteRq []byte) (map[string][]*cryptoCompare, error) {
	var p fastjson.Parser
	parsed, err := p.ParseBytes(byteRq)
	if err != nil {
		return nil, errors.Wrap(err, "crcFastJson")
	}

	m := make(map[string][]*cryptoCompare)

	o := parsed.GetObject("RAW")
	o.Visit(func(k []byte, v *fastjson.Value) {
		if val, ok := s.list[string(k)]; ok {
			crypto := v.GetObject()
			crypto.Visit(func(key []byte, value *fastjson.Value) {

				var c cryptoCompare
				if err := json.Unmarshal([]byte(value.String()), &c); err != nil {
					logger.Error("o.Visit", err)
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

func (s *Service) crcPricesRequest(tsyms, fsyms string, c chan<- map[string][]*cryptoCompare, wg *sync.WaitGroup) {
	rq, err := req.Get(urlCRC, req.Param{
		"fsyms": fsyms,
		"tsyms": tsyms,
	})
	if err != nil {
		logger.Error("crcPricesRequest", err)
		return
	}

	byteRq := rq.Bytes()
	m, err := s.crcFastJson(byteRq)
	if err != nil {
		logger.Error("crcPricesRequest", err)
		return
	}

	c <- m
	defer wg.Done()
}

func fiatMapping(c chan map[string][]*cryptoCompare) storage.FiatMap {
	if c == nil {
		return nil
	}
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

func (s *Service) GetPricesHUOBI() (storage.FiatMap, error) {
	rq, err := req.Get(urlHuobi)
	if err != nil {
		return nil, errors.Wrap(err, "huobi")
	}

	var h huobi
	if err := rq.ToJSON(&h); err != nil {
		return nil, errors.Wrap(err, "toJSON huobi")
	}
	return huobiMapping(&h, s.list), nil
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

func (s *Service) GetPricesTrustV2(prices PricesTrustV2) (storage.FiatMap, error) {
	rq := req.New()
	resp, err := rq.Post(urlTrustWalletV2, req.BodyJSON(&prices))
	if err != nil {
		return nil, errors.Wrap(err, "GetPricesTrustV2")
	}

	var r trustV2Response
	if err := resp.ToJSON(&r); err != nil {
		return nil, errors.Wrap(err, "GetPricesTrustV2toJSON")
	}

	m := trustV2FiatMap(&r)
	return m, nil
}

func trustV2FiatMap(r *trustV2Response) storage.FiatMap {
	m := make(storage.FiatMap)
	currency := storage.Fiat(r.Currency)
	m[currency] = map[storage.CryptoCurrency]*storage.Details{}
	for _, doc := range r.Docs {
		coin := storage.CryptoCurrency(strconv.Itoa(doc.Coin))

		sd := storageDetails(doc.Price.Value, 0, doc.Price.Change24H, 0)
		m[currency][coin] = &sd
	}

	return m
}
