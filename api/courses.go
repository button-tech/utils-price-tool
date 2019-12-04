package api

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/storage"
	"github.com/button-tech/utils-price-tool/pkg/typeconv"
	"github.com/pkg/errors"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type request struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

type uniqueRequest struct {
	tokens     map[string]struct{}
	currencies map[string]struct{}
	change     string
	api        string
}

type response struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

type privateCMC struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Quote  quote  `json:"quote"`
}

type coinMarketPrices struct {
	price        float64
	change24Hour float64
	change7Day   float64
}

type quote struct {
	USD usd
}

type usd struct {
	Price            float64 `json:"price"`
	PercentChange24H float64 `json:"percent_change_24h"`
	PercentChange7D  float64 `json:"percent_change_7d"`
}

type privateInputCurrencies struct {
	Currencies []string `json:"currencies"`
}

type api struct {
	Name             string   `json:"name"`
	SupportedChanges []string `json:"supported_changes"`
}

var supportAPIs = map[string]struct{}{
	"crc":    {},
	"cmc":    {},
	"huobi":  {},
	"ntrust": {},
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

func (s *Server) initCoursesAPI() {
	controller := apiController{
		privateCurrencies: s.privateCurrencies,
		store:             s.store,
	}
	s.G.Post("/prices", controller.getCourses)
	s.G.Post("/change", controller.privatePrices)
	s.G.Get("/list", controller.apiInfo)
}

func (s *Server) initCoursesAPIv2() {
	controller := apiController{
		store: s.store,
	}
	s.Gv2.Post("/prices", controller.getCoursesV2)
	s.Gv2.Get("/info", controller.getInfoV2)
}

func (ac *apiController) getCourses(ctx *routing.Context) error {
	var req request
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		logger.Error("getCourses", err, logger.Params{
			"from": "json.Unmarshal",
		})
		return err
	}

	a := req.API
	switch a {
	case "cmc", "crc", "huobi":
		result, err := ac.converter(&req)
		if err != nil {
			respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			logger.Error("getCourses", err.Error(), logger.Params{
				"from": "ac.converter",
			})
			return nil
		}
		respondWithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{
			"data": result,
		})
		return nil

	default:
		supportedCRC := []string{"0", "1", "24"}
		crc := api{
			Name:             "crc",
			SupportedChanges: supportedCRC,
		}

		supportedCMC := []string{"24"}
		cmc := api{
			Name:             "cmc",
			SupportedChanges: supportedCMC,
		}

		API := []api{crc, cmc}

		respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			"api":   API,
			"error": "please, use these API",
		})
		return nil
	}
}

func (ac *apiController) apiInfo(ctx *routing.Context) error {
	supportedCRC := []string{"0", "1", "24"}
	crc := api{
		Name:             "crc",
		SupportedChanges: supportedCRC,
	}

	supportedCMC := []string{"24"}
	cmc := api{
		Name:             "cmc",
		SupportedChanges: supportedCMC,
	}

	supportedHuobi := []string{"0"}
	huobi := api{
		Name:             "huobi",
		SupportedChanges: supportedHuobi,
	}

	API := []api{crc, cmc, huobi}
	respondWithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{
		"api": API,
	})

	return nil
}

func (ac *apiController) mapping(req *uniqueRequest) ([]response, error) {
	result := make([]response, 0, len(req.currencies))
	api := req.api
	stored, err := ac.store.Get(typeconv.StorageApi(api))
	if err != nil {
		return nil, err
	}

	for c := range req.currencies {
		price := response{}
		if fiatVal, fiatOk := stored[typeconv.StorageFiat(c)]; fiatOk {
			price.Currency = c
			for t := range req.tokens {
				currency := tokenToStorageCC(req.api, t)
				if details, ok := fiatVal[currency]; ok {
					contract := map[string]string{t: details.Price}
					if err := changesControl(contract, details, req.change); err != nil {
						return nil, err
					}
					price.Rates = append(price.Rates, contract)
				}
			}
		}
		if price.Currency != "" {
			result = append(result, price)
		}
	}
	return result, nil
}

func tokenToStorageCC(api, t string) (c storage.CryptoCurrency) {
	if api == "ntrust" {
		c = typeconv.StorageCC(t)
		return
	}
	c = typeconv.StorageCC(strings.ToLower(t))
	return
}

func changesControl(m map[string]string, d *storage.Details, c string) error {
	switch c {
	case "1":
		if d.ChangePCTHour != "" {
			m["percent_change"] = d.ChangePCTHour
			return nil
		}
	case "24":
		if d.ChangePCT24Hour != "" {
			m["percent_change"] = d.ChangePCT24Hour
			return nil
		}
	default:
		return nil
	}
	return errors.New("API changes: no matches")
}

func (ac *apiController) converter(req *request) ([]response, error) {
	api := req.API
	if _, ok := supportAPIs[api]; !ok {
		return nil, errors.New("API: no matches")
	}
	u := unique(req)
	return ac.mapping(u)
}

func unique(req *request) *uniqueRequest {
	uniqueTokens := make(map[string]struct{})
	uniqueCurrencies := make(map[string]struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup, ut map[string]struct{}) {
		for _, t := range req.Tokens {
			if _, ok := ut[t]; !ok {
				ut[t] = struct{}{}
			}
		}
		wg.Done()
	}(&wg, uniqueTokens)

	wg.Add(1)
	go func(wg *sync.WaitGroup, uc map[string]struct{}) {
		for _, c := range req.Currencies {
			if _, ok := uc[c]; !ok {
				uc[c] = struct{}{}
			}
		}
		wg.Done()
	}(&wg, uniqueCurrencies)
	wg.Wait()

	return &uniqueRequest{
		tokens:     uniqueTokens,
		currencies: uniqueCurrencies,
		change:     req.Change,
		api:        req.API,
	}
}

func privateCurrencies() map[string][]string {
	return map[string][]string{
		"BTC":   {"0x0000000000000000000000000000000000000000", "Bitcoin"},
		"ETH":   {"0x000000000000000000000000000000000000003c", "Ethereum"},
		"ETC":   {"0x000000000000000000000000000000000000003d", "Ethereum Classic"},
		"BCH":   {"0x0000000000000000000000000000000000000091", "Bitcoin Cash"},
		"LTC":   {"0x0000000000000000000000000000000000000002", "Litecoin"},
		"XLM":   {"0x0000000000000000000000000000000000000094", "Stellar"},
		"WAVES": {"0x0000000000000000000000000000000000579bfc", "Waves"},
	}
}

func (ac *apiController) privatePrices(ctx *routing.Context) error {
	var r privateInputCurrencies
	if err := json.Unmarshal(ctx.PostBody(), &r); err != nil {
		logger.Error("privatePrices", err)
		respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{"err": "can't unmarshal body"})
		return nil
	}

	currencies := make([]privateCMC, 0, len(r.Currencies))
	stored, err := ac.store.Get(typeconv.StorageApi("coinMarketCap"))
	if err != nil {
		return errors.Wrap(err, "privatePrices")
	}
	for _, symbol := range r.Currencies {
		currDetail := ac.privateCurrencies[symbol]

		bip := currDetail[0]
		name := currDetail[1]

		val := stored[typeconv.StorageFiat("USD")]
		details := val[typeconv.StorageCC(bip)]
		priceInfo, err := coinMarketPricesInfo(details.Price, details.ChangePCT24Hour, details.ChangePCT7Day)
		if err != nil {
			return errors.Wrap(err, "privatePrices")
		}

		q := priceQuote(priceInfo)
		currencies = append(currencies, privateCMC{
			Name:   name,
			Symbol: symbol,
			Quote:  q,
		})
	}

	respondWithJSON(ctx, 200, map[string]interface{}{"data": &currencies})
	return nil
}

func priceQuote(info *coinMarketPrices) quote {
	return quote{USD: usd{
		Price:            info.price,
		PercentChange24H: info.change24Hour,
		PercentChange7D:  info.change7Day,
	}}
}

func coinMarketPricesInfo(price, hour24, sevenDay string) (*coinMarketPrices, error) {
	convPrice, err := strconv.ParseFloat(price, 10)
	if err != nil {
		return nil, errors.Wrap(err, "priceConversion")
	}

	change24Hour, err := strconv.ParseFloat(hour24, 10)
	if err != nil {
		return nil, errors.Wrap(err, "change24HourConversion")
	}

	change7Day, err := strconv.ParseFloat(sevenDay, 10)
	if err != nil {
		return nil, errors.Wrap(err, "change7DayConversion")
	}

	return &coinMarketPrices{
		price:        convPrice,
		change24Hour: change24Hour,
		change7Day:   change7Day,
	}, nil
}

func (ac *apiController) getInfoV2(ctx *routing.Context) error {
	respondWithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"support": trustV2Coins})
	return nil
}

func (ac *apiController) getCoursesV2(ctx *routing.Context) error {
	var req request
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		logger.Error("getCourses", err, logger.Params{
			"from": "json.Unmarshal",
		})
		return err
	}

	a := req.API
	switch a {
	case "ntrust":
		result, err := ac.converter(&req)
		if err != nil {
			respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			logger.Error("getCourses", err.Error(), logger.Params{
				"from": "ac.converter",
			})
			return nil
		}
		respondWithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{
			"data": result,
		})
	default:
		supported := v2SupportedInfo()
		respondWithJSON(ctx, fasthttp.StatusBadRequest, map[string]interface{}{
			"api":   supported,
			"error": "please, use these API",
		})
	}
	return nil
}

func v2SupportedInfo() api {
	supportedNewTrust := []string{"0", "24"}
	return api{
		Name:             "ntrust",
		SupportedChanges: supportedNewTrust,
	}
}
