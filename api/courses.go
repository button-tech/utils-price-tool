package api

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/storage"
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

//type listApi struct {
//	API []api `json:"api"`
//	//Time             struct {
//	//	Start int `json:"start"`
//	//	End   int `json:"end"`
//	//} `json:"time"`
//}

type api struct {
	Name             string   `json:"name"`
	SupportedChanges []string `json:"supported_changes"`
}

var supportAPIs = map[string]struct{}{
	"crc":   {},
	"cmc":   {},
	"huobi": {},
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
		result, err := ac.converter(&req, a)
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

func (ac *apiController) mapping(req *uniqueRequest, api string) []*response {
	result := make([]*response, 0)
	stored := ac.store.Get()[storage.Api(api)]
	if stored == nil {
		return nil
	}

	for c := range req.currencies {
		price := response{}

		if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
			price.Currency = c

			for t := range req.tokens {
				if val, ok := fiatVal[storage.CryptoCurrency(strings.ToLower(t))]; ok {
					contract := map[string]string{t: val.Price}
					if contract = changesControl(contract, val, req.change); len(contract) == 0 {
						return nil
					} else {
						price.Rates = append(price.Rates, contract)
					}
				}
			}
		}
		if price.Currency != "" {
			result = append(result, &price)
		}
	}
	return result
}

func changesControl(m map[string]string, s *storage.Details, c string) map[string]string {
	switch c {
	case "1":
		if s.ChangePCTHour != "" {
			m["percent_change"] = s.ChangePCTHour
			return m
		}
		return nil
	case "24":
		if s.ChangePCT24Hour != "" {
			m["percent_change"] = s.ChangePCT24Hour
			return m
		}
		return nil
	default:
		return m
	}
}

func (ac *apiController) converter(req *request, api string) ([]*response, error) {
	if _, ok := supportAPIs[api]; !ok {
		return nil, errors.New("no matches API")
	}

	resp := ac.mapping(unique(req), api)
	if resp == nil {
		return nil, errors.New("no matches support changes API")
	}
	return resp, nil
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
		"BTC":   []string{"0x0000000000000000000000000000000000000000", "Bitcoin"},
		"ETH":   []string{"0x000000000000000000000000000000000000003c", "Ethereum"},
		"ETC":   []string{"0x000000000000000000000000000000000000003d", "Ethereum Classic"},
		"BCH":   []string{"0x0000000000000000000000000000000000000091", "Bitcoin Cash"},
		"LTC":   []string{"0x0000000000000000000000000000000000000002", "Litecoin"},
		"XLM":   []string{"0x0000000000000000000000000000000000000094", "Stellar"},
		"WAVES": []string{"0x0000000000000000000000000000000000579bfc", "Waves"},
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
	stored := ac.store.Get()["coinMarketCap"]
	for _, symbol := range r.Currencies {
		currDetail := ac.privateCurrencies[symbol]

		bip := currDetail[0]
		name := currDetail[1]

		val := stored[storage.Fiat("USD")]
		details := val[storage.CryptoCurrency(bip)]
		priceInfo, err := coinMarketPricesInfo(details.Price, details.ChangePCT24Hour, details.ChangePCT7Day)
		if err != nil {
			return errors.Wrap(err, "privatePrices")
		}

		//a := privateCMC{
		//	Name: c,
		//	Symbol: sybmol,
		//	Quote: struct {
		//		USD struct {
		//			Price            float64 `json:"price"`;
		//			PercentChange24H float64 `json:"percent_change_24h"`;
		//			PercentChange7D  float64 `json:"percent_change_7d"`
		//		}
		//	}{USD: struct {
		//		Price:
		//	}{}}
		//}

		u := usd{
			Price:            priceInfo.price,
			PercentChange24H: priceInfo.change24Hour,
			PercentChange7D:  priceInfo.change7Day,
		}

		q := quote{
			USD: u,
		}

		currencies = append(currencies, privateCMC{
			Name:  name ,
			Symbol: symbol,
			Quote:  q,
		})
	}

	respondWithJSON(ctx, 200, map[string]interface{}{"data": &currencies})
	return nil
}

type coinMarketPrices struct {
	price        float64
	change24Hour float64
	change7Day   float64
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
		return nil, errors.Wrap(err, "priceConversion")
	}

	return &coinMarketPrices{
		price:        convPrice,
		change24Hour: change24Hour,
		change7Day:   change7Day,
	}, nil
}

func (s *Server) initCoursesAPI() {
	s.G.Post("/prices", s.ac.getCourses)
	s.G.Post("/change", s.ac.privatePrices)
	s.G.Get("/list", s.ac.apiInfo)
}
