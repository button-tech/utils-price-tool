package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/currencies"
	"github.com/button-tech/utils-price-tool/core/prices"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"sync"
)

var TrustWalletURL = os.Getenv("TRUST_URL")

func CmcUpdateWorker(wg *sync.WaitGroup, p *prices.PricesData) {
	defer wg.Done()
	tokens := CreateCMCRequestData(p)

	var tokensWG sync.WaitGroup
	for _, t := range tokens {
		tokensWG.Add(1)
		go func(token types.TokensWithCurrency, tWG *sync.WaitGroup) {
			defer tWG.Done()
			if err := SetPricesCMC(token, p); err != nil {
				logger.Error("cmcWorker", err)
				return
			}
		}(t, &tokensWG)
	}
	tokensWG.Wait()
}

func SetPricesCMC(tokens types.TokensWithCurrency, p *prices.PricesData) error {
	var cmc types.CoinMarketCap

	res, err := req.Post(TrustWalletURL, req.BodyJSON(tokens))
	if err != nil {
		return errors.Wrap(err, "PricesCMC")
	}

	if res.Response().StatusCode != 200 {
		return errors.Wrap(errors.New("error"), "PricesCMC")
	}

	if err = res.ToJSON(&cmc); err != nil {
		return errors.Wrap(err, "PricesCMC")
	}

	for _, v := range cmc.Docs {
		details := cache.Details{}
		details.Price = v.Price
		details.ChangePCT24Hour = v.PercentChange24H

		k := cache.GenKey("cmc", cmc.Currency, v.Contract)
		p.Store.Set(k, details)
	}
	return nil
}

func CreateCMCRequestData(p *prices.PricesData) []types.TokensWithCurrency {
	var tokensOneCurrency types.TokensWithCurrency

	tokensMultiCurrencies := make([]types.TokensWithCurrency, 0, len(currencies.SupportedCurrenciesList))

	tokens := make([]types.Token, 0, len(p.List))

	for _, c := range p.List {
		tokens = append(tokens, types.Token{Contract: c})
	}

	for _, t := range p.Tokens {
		tokens = append(tokens, types.Token{Contract: t})
	}
	tokensOneCurrency.Tokens = tokens

	for _, c := range currencies.SupportedCurrenciesList {
		tokensOneCurrency.Currency = c
		tokensMultiCurrencies = append(tokensMultiCurrencies, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}
