package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"sync"
)

var TrustWalletURL = os.Getenv("TRUST_URL")

func CmcUpdateWorker(wg *sync.WaitGroup, p *Prices) {
	defer wg.Done()
	tokens := p.CreateCMCRequestData()

	var tokensWG sync.WaitGroup
	for _, t := range tokens {
		tokensWG.Add(1)
		go func(token TokensWithCurrency, tWG *sync.WaitGroup) {
			defer tWG.Done()
			if err := p.SetPricesCMC(token); err != nil {
				logger.Error("cmcWorker", err)
				return
			}
		}(t, &tokensWG)
	}
	tokensWG.Wait()
}

func (p *Prices) SetPricesCMC(tokens TokensWithCurrency) error {
	var cmc CoinMarketCap

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
		p.store.Set(k, details)
	}
	return nil
}

func (p *Prices) CreateCMCRequestData() []TokensWithCurrency {
	var tokensOneCurrency TokensWithCurrency
	tokensMultiCurrencies := make([]TokensWithCurrency, 0, len(currencies))
	tokens := make([]Token, 0, len(p.List))

	for _, c := range p.List {
		tokens = append(tokens, Token{Contract: c})
	}

	for _, t := range p.Tokens {
		tokens = append(tokens, Token{Contract: t})
	}
	tokensOneCurrency.Tokens = tokens

	for _, c := range currencies {
		tokensOneCurrency.Currency = c
		tokensMultiCurrencies = append(tokensMultiCurrencies, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}
