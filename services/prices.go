package services

import (
	"encoding/json"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/button-tech/logger"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

const (
	urlHuobi   = "https://api.hbdm.com/api/v1/contract_index"
	urlTopList = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=USD"
	urlCRC     = "https://min-api.cryptocompare.com/data/pricemultifull"
)

const coin = "coin"

var (
	topListAPIKey    = os.Getenv("API_KEY")
	trustWalletURL   = os.Getenv("TRUST_URL")
	trustWalletV2URL = os.Getenv("TRUST_URL_V2")
)

type GetPrices struct {
	mu           sync.Mutex
	TrustV2Coins []PricesTrustV2
	List         map[string]string
	Tokens       map[string]string

	store *cache.Cache
}

var PureCMCCoins = map[string]int{
	"AE":    457,
	"ALGO":  283,
	"ATOM":  118,
	"BCH":   145,
	"BNB":   714,
	"BTC":   0,
	"DASH":  5,
	"DCR":   42,
	"DGB":   20,
	"DOGE":  3,
	"ETC":   61,
	"ETH":   60,
	"ICX":   74,
	"LTC":   2,
	"NANO":  165,
	"ONT":   1024,
	"QTUM":  2301,
	"RVN":   175,
	"THETA": 500,
	"TRX":   195,
	"VET":   818,
	"WAVES": 5741564,
	"XLM":   148,
	"XRP":   144,
	"XTZ":   1729,
	"ZEC":   133,
	"ZIL":   313,
}

var TrustV2Coins = map[string]int{
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

func New(store *cache.Cache) *GetPrices {
	return &GetPrices{
		TrustV2Coins: CreateTrustV2RequestData(),
		List:         make(map[string]string),
		store:        store,
	}
}

func CreateTrustV2RequestData() []PricesTrustV2 {
	prices := make([]PricesTrustV2, 0, len(currencies))
	for _, c := range currencies {
		price := PricesTrustV2{Currency: c}
		allAssets := make([]Assets, 0, len(TrustV2Coins))
		for _, v := range TrustV2Coins {
			allAssets = append(allAssets, Assets{Coin: v, Type: coin})
		}
		price.Assets = allAssets
		prices = append(prices, price)
	}
	return prices
}

func (s *GetPrices) CreateCMCRequestData() []TokensWithCurrency {
	var tokensOneCurrency TokensWithCurrency
	tokensMultiCurrencies := make([]TokensWithCurrency, 0, len(currencies))
	tokens := make([]Token, 0, len(s.List))

	for _, c := range s.List {
		tokens = append(tokens, Token{Contract: c})
	}

	for _, t := range s.Tokens {
		tokens = append(tokens, Token{Contract: t})
	}
	tokensOneCurrency.Tokens = tokens

	for _, c := range currencies {
		tokensOneCurrency.Currency = c
		tokensMultiCurrencies = append(tokensMultiCurrencies, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}

// Get top List of crypto-currencies from coin-market
func (s *GetPrices) GetTopList(c map[string]string) error {
	rq, err := req.Get(urlTopList, req.Header{"X-CMC_PRO_API_KEY": topListAPIKey})
	if err != nil {
		return errors.Wrap(err, "getTopList")
	}

	var topList pureCoinMarketCap
	if err = rq.ToJSON(&topList); err != nil {
		return errors.Wrap(err, "getTopList")
	}

	if topList.Status.ErrorCode != 0 {
		return errors.New("responseHTTPStatus: NotOk")
	}

	top100 := topList.Data[:100]
	topListMap := make(map[string]string)
	for _, item := range top100 {
		if val, ok := c[item.Symbol]; ok {
			topListMap[item.Symbol] = val
			pricesData := detailsConversion(
				item.Quote.USD.Price,
				item.Quote.USD.PercentChange1H,
				item.Quote.USD.PercentChange24H,
				item.Quote.USD.PercentChange7D,
			)
			k := cache.GenKey("coinMarketCap", "usd", item.Symbol)
			s.store.Set(k, pricesData)
		}
	}
	pureCMCMapping(topList, s.store)
	s.List = topListMap
	return nil
}

func detailsConversion(price, hour, hour24, sevenDay float64) cache.Details {
	d := cache.Details{Price: strconv.FormatFloat(price, 'f', 10, 64)}
	if floatValid(hour) {
		d.ChangePCTHour = strconv.FormatFloat(hour, 'f', 6, 64)
	}
	if floatValid(hour24) {
		d.ChangePCT24Hour = strconv.FormatFloat(hour24, 'f', 6, 64)
	}
	if floatValid(sevenDay) {
		d.ChangePCT7Day = strconv.FormatFloat(sevenDay, 'f', 6, 64)
	}
	return d
}

func floatValid(s float64) bool {
	return s != 0
}

func (s *GetPrices) GetPricesCMC(tokens TokensWithCurrency) error {
	rq, err := req.Post(trustWalletURL, req.BodyJSON(tokens))
	if err != nil {
		return errors.Wrap(err, "GetPricesCMC")
	}

	gotPrices := coinMarketCap{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return errors.Wrap(err, "GetPricesCMC")
	}

	for _, v := range gotPrices.Docs {
		details := cache.Details{}
		details.Price = v.Price
		details.ChangePCT24Hour = v.PercentChange24H

		k := cache.GenKey("cmc", gotPrices.Currency, v.Contract)
		s.store.Set(k, details)
	}
	return nil
}

func (s *GetPrices) GetTokenPriceCMC(token TokensWithCurrency) (string, error) {
	rq, err := req.Post(trustWalletURL, req.BodyJSON(token))
	if err != nil {
		return "", errors.Wrap(err, "GetPricesCMC")
	}

	gotPrices := coinMarketCap{}
	if err = rq.ToJSON(&gotPrices); err != nil {
		return "", errors.Wrap(err, "GetPricesCMC")
	}

	doc := gotPrices.Docs[0]
	return doc.Price, nil
}

func CreateCRCRequestData() []string {
	sortedCurrencies := make([]string, 0, 7)
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

func (s *GetPrices) GetPricesCRC() {
	var fsyms string
	for k := range s.List {
		fsyms += k + ","
	}
	sortedCurrencies := CreateCRCRequestData()
	c := make(chan map[string][]cryptoCompare, len(sortedCurrencies))

	var wg sync.WaitGroup
	for _, tsyms := range sortedCurrencies {
		wg.Add(1)
		go s.crcPricesRequest(tsyms, fsyms, c, &wg)
	}
	wg.Wait()
	close(c)
	fiatMapping(c, s.store)
}

func (s *GetPrices) crcFastJson(byteRq []byte) (map[string][]cryptoCompare, error) {
	var p fastjson.Parser
	parsed, err := p.ParseBytes(byteRq)
	if err != nil {
		return nil, errors.Wrap(err, "crcFastJson")
	}

	m := make(map[string][]cryptoCompare)
	o := parsed.GetObject("RAW")
	o.Visit(func(k []byte, v *fastjson.Value) {
		if val, ok := s.List[string(k)]; ok {
			crypto := v.GetObject()
			crypto.Visit(func(key []byte, value *fastjson.Value) {
				var c cryptoCompare
				if err := json.Unmarshal([]byte(value.String()), &c); err != nil {
					logger.Error("o.Visit", err)
					return
				}
				valM, okM := m[c.ToSymbol]
				if !okM {
					m[c.ToSymbol] = make([]cryptoCompare, 0)
				}
				valM = append(valM, c)
				formatCryptoCompare(&c, val)
				valM = append(valM, c)
				m[c.ToSymbol] = valM
			})
		}
	})
	return m, nil
}

func formatCryptoCompare(c *cryptoCompare, from string) {
	c.FromSymbol = from
}

func (s *GetPrices) crcPricesRequest(tsyms, fsyms string, c chan<- map[string][]cryptoCompare, wg *sync.WaitGroup) {
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

func fiatMapping(c chan map[string][]cryptoCompare, store *cache.Cache) {
	done := false
	for !done {
		select {
		case m, ok := <-c:
			if !ok {
				done = true
				break
			}
			for k, v := range m {
				for _, i := range v {
					details := cache.Details{}
					details.Price = strconv.FormatFloat(i.Price, 'f', -1, 64)
					details.ChangePCT24Hour = strconv.FormatFloat(i.ChangePCT24Hour, 'f', 2, 64)
					details.ChangePCTHour = strconv.FormatFloat(i.ChangePCTHour, 'f', 2, 64)
					k := cache.GenKey("crc", k, i.FromSymbol)
					store.Set(k, details)
				}
			}
		}
	}
}

func (s *GetPrices) GetPricesHUOBI() error {
	rq, err := req.Get(urlHuobi)
	if err != nil {
		return errors.Wrap(err, "huobi")
	}

	var h huobi
	if err := rq.ToJSON(&h); err != nil {
		return errors.Wrap(err, "toJSON huobi")
	}
	huobiMapping(&h, s.List, s.store)
	return nil
}

func huobiMapping(h *huobi, list map[string]string, store *cache.Cache) {
	for _, i := range h.Data {
		if val, ok := list[i.Symbol]; ok {
			var details cache.Details
			details.Price = strconv.FormatFloat(i.IndexPrice, 'f', -1, 64)
			k := cache.GenKey("huobi", "usd", val)
			store.Set(k, details)
		}
	}
}

func (s *GetPrices) GetPricesTrustV2(prices PricesTrustV2) error {
	rq := req.New()
	resp, err := rq.Post(trustWalletV2URL, req.BodyJSON(&prices))
	if err != nil {
		return errors.Wrap(err, "GetPricesTrustV2")
	}

	var r trustV2Response
	if err := resp.ToJSON(&r); err != nil {
		return errors.Wrap(err, "GetPricesTrustV2toJSON")
	}
	trustV2FiatMap(&r, s.store)
	return nil
}

func trustV2FiatMap(r *trustV2Response, store *cache.Cache) {
	for _, doc := range r.Docs {
		coin := strconv.Itoa(doc.Coin)
		sd := detailsConversion(doc.Price.Value, 0, doc.Price.Change24H, 0)
		k := cache.GenKey("ntrust", r.Currency, coin)
		store.Set(k, sd)
	}
}

func pureCMCMapping(pure pureCoinMarketCap, store *cache.Cache) {
	for _, v := range pure.Data {
		if coinID, ok := TrustV2Coins[v.Symbol]; ok {
			pricesData := detailsConversion(
				v.Quote.USD.Price,
				v.Quote.USD.PercentChange1H,
				v.Quote.USD.PercentChange24H,
				v.Quote.USD.PercentChange7D,
			)
			convCoinID := strconv.Itoa(coinID)
			kt := cache.GenKey("pcmc", "usd", convCoinID)
			store.Set(kt, pricesData)
		}
	}
}
