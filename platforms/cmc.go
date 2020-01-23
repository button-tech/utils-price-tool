package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/currencies"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"sync"
)

var (
	TrustWalletURL = os.Getenv("TRUST_URL")
)

func CmcUpdateWorker(wg *sync.WaitGroup, c *cache.Cache) {
	defer wg.Done()
	tokens := CreateCMCRequestData(c)

	var tokensWG sync.WaitGroup
	for _, t := range tokens {
		tokensWG.Add(1)
		go func(token types.TokensWithCurrency, tWG *sync.WaitGroup) {
			defer tWG.Done()
			if err := SetPricesCMC(token, c); err != nil {
				logger.Error("cmcWorker", err)
				return
			}
		}(t, &tokensWG)
	}
	tokensWG.Wait()
}

func SetPricesCMC(tokens types.TokensWithCurrency, p *cache.Cache) error {
	var cmc types.CmcResponse

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
		p.Set(k, details)
	}
	return nil
}

func CreateCMCRequestData(p *cache.Cache) []types.TokensWithCurrency {
	var tokensOneCurrency types.TokensWithCurrency

	tokensMultiCurrencies := make([]types.TokensWithCurrency, 0, len(currencies.SupportedCurrenciesList))

	tokens := make([]types.Token, 0, len(p.List))

	for _, c := range p.List {
		tokens = append(tokens, types.Token{Contract: c})
	}

	tokensOneCurrency.Tokens = tokens

	for _, c := range currencies.SupportedCurrenciesList {
		tokensOneCurrency.Currency = c
		tokensMultiCurrencies = append(tokensMultiCurrencies, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}
