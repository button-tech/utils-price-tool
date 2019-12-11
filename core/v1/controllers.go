package v1

// Utils-price-tool v1 API.
//
// This project is included in Button-Wallet Utils-price-tool project
//
//     Schemes: http
//     Host: localhost
//     BasePath: /courses/v1/
//     Version: 0.0.1
//     Contact: Frolov Ivan <if@buttonwallet.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta

import (
	"encoding/json"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"strconv"

	"github.com/button-tech/utils-price-tool/core/internal/handle"
	"github.com/button-tech/utils-price-tool/core/internal/respond"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const (
	v1 = "v1"
)

type privateInputCurrencies struct {
	Currencies []string `json:"currencies"`
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

type coinMarketPrices struct {
	price        float64
	change24Hour float64
	change7Day   float64
}

func (c *controller) courses(ctx *routing.Context) error {
	const funcName = "courses"
	var r handle.Data
	if err := json.Unmarshal(ctx.PostBody(), &r); err != nil {
		respond.WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, respond.Error{
			API:     v1,
			Func:    funcName,
			Err:     err,
			Payload: respond.Payload("request", "json.Unmarshal"),
		}, nil)
		return nil
	}

	unique := handle.Unify(&r)
	resp, err := handle.Reply(&unique, v1, c.store, c.getPrices)
	if err != nil {
		respond.WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, respond.Error{
			API:     v1,
			Func:    funcName,
			Err:     err,
			Payload: respond.Payload("response", "handle.Reply"),
		}, map[string]interface{}{"api": supportInfo(), "error": "please, use these API"})
		return nil
	}

	respond.WithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"data": resp})
	return nil
}

func (c *controller) info(ctx *routing.Context) error {
	respond.WithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"api": supportInfo()})
	return nil
}

func (c *controller) privatePrices(ctx *routing.Context) error {
	const funcName = "privatePrices"
	var r privateInputCurrencies
	if err := json.Unmarshal(ctx.PostBody(), &r); err != nil {
		respond.WithWrapErrJSON(ctx, fasthttp.StatusBadRequest, respond.Error{
			API:     v1,
			Func:    funcName,
			Err:     err,
			Payload: respond.Payload("request", "json.Unmarshal"),
		}, nil)
		return nil
	}

	currencies := make([]privateCMC, 0, len(r.Currencies))
	for _, symbol := range r.Currencies {
		currDetail := c.privateCurrencies[symbol]

		bip := currDetail[0]
		name := currDetail[1]

		k := cache.GenKey("coinMarketCap", "usd", bip)
		d, ok := c.store.Get(k)
		if ok {
			priceInfo, err := coinMarketPricesInfo(d.Price, d.ChangePCT24Hour, d.ChangePCT7Day)
			if err != nil {
				return errors.Wrap(err, "privatePrices")
			}

			q := priceQuote(&priceInfo)
			currencies = append(currencies, privateCMC{
				Name:   name,
				Symbol: symbol,
				Quote:  q,
			})
		}
	}

	respond.WithJSON(ctx, fasthttp.StatusOK, map[string]interface{}{"data": currencies})
	return nil
}

func supportInfo() []handle.APIs {
	supportedCRC := []string{"0", "1", "24"}
	crc := handle.APIs{
		Name:             "crc",
		SupportedChanges: supportedCRC,
	}

	supportedCMC := []string{"24"}
	cmc := handle.APIs{
		Name:             "cmc",
		SupportedChanges: supportedCMC,
	}

	supportedHuobi := []string{"0"}
	huobi := handle.APIs{
		Name:             "huobi",
		SupportedChanges: supportedHuobi,
	}

	return []handle.APIs{crc, cmc, huobi}
}

func coinMarketPricesInfo(price, hour24, sevenDay string) (coinMarketPrices, error) {
	convPrice, err := strconv.ParseFloat(price, 10)
	if err != nil {
		return coinMarketPrices{}, errors.Wrap(err, "priceConversion")
	}

	change24Hour, err := strconv.ParseFloat(hour24, 10)
	if err != nil {
		return coinMarketPrices{}, errors.Wrap(err, "change24HourConversion")
	}

	change7Day, err := strconv.ParseFloat(sevenDay, 10)
	if err != nil {
		return coinMarketPrices{}, errors.Wrap(err, "change7DayConversion")
	}

	return coinMarketPrices{
		price:        convPrice,
		change24Hour: change24Hour,
		change7Day:   change7Day,
	}, nil
}

func priceQuote(info *coinMarketPrices) quote {
	return quote{USD: usd{
		Price:            info.price,
		PercentChange24H: info.change24Hour,
		PercentChange7D:  info.change7Day,
	}}
}
