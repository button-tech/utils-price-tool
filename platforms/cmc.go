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
			if err := SetCMC(token, c); err != nil {
				logger.Error("cmcWorker", err)
				return
			}
		}(t, &tokensWG)
	}
	tokensWG.Wait()
}

func SetCMC(tokens types.TokensWithCurrency, c *cache.Cache) error {
	cmc, err := RequestCMC(tokens)
	if err != nil {
		return err
	}

	for _, v := range cmc.Docs {
		details := cache.Details{}
		details.Price = v.Price
		details.ChangePCT24Hour = v.PercentChange24H

		k := cache.GenKey("cmc", cmc.Currency, v.Contract)
		c.Set(k, details)
	}
	return nil
}

func SetWithGettingDetails(tokens types.TokensWithCurrency, c *cache.Cache) ([]cache.Details, error) {
	cmc, err := RequestCMC(tokens)
	if err != nil {
		return nil, err
	}

	detailsList := make([]cache.Details, len(cmc.Docs))

	for i, v := range cmc.Docs {
		details := cache.Details{}
		details.Price = v.Price
		details.ChangePCT24Hour = v.PercentChange24H
		details.Token = v.Contract
		details.Fiat = tokens.Currency
		detailsList[i] = details
		k := cache.GenKey("cmc", cmc.Currency, v.Contract)
		c.Set(k, details)
	}

	return detailsList, nil
}

func RequestCMC(tokens types.TokensWithCurrency) (*types.CmcResponse, error) {
	var cmc types.CmcResponse

	res, err := req.Post(TrustWalletURL, req.BodyJSON(tokens))
	if err != nil {
		return nil, errors.Wrap(err, "PricesCMC")
	}

	if res.Response().StatusCode != 200 {
		return nil, errors.Wrap(errors.New("error"), "PricesCMC")
	}

	if err = res.ToJSON(&cmc); err != nil {
		return nil, errors.Wrap(err, "PricesCMC")
	}

	return &cmc, nil
}

func CreateCMCRequestData(c *cache.Cache) []types.TokensWithCurrency {
	var tokensOneCurrency types.TokensWithCurrency

	tokensMultiCurrencies := make([]types.TokensWithCurrency, 0, len(currencies.SupportedCurrenciesList))

	tokens := make([]types.Token, 0, len(c.List))

	for _, v := range c.List {
		tokens = append(tokens, types.Token{Contract: v})
	}

	tokensOneCurrency.Tokens = tokens

	for _, v := range currencies.SupportedCurrenciesList {
		tokensOneCurrency.Currency = v
		tokensMultiCurrencies = append(tokensMultiCurrencies, tokensOneCurrency)
	}

	return tokensMultiCurrencies
}

func SingleERC20Course(fiat, crypto string) (string, error) {

	var cmc types.CmcResponse

	token := make([]types.Token, 0, 1)
	token = append(token, types.Token{Contract: crypto})

	singleToken := types.TokensWithCurrency{
		Currency: fiat,
		Tokens:   token,
	}

	res, err := req.Post(TrustWalletURL, req.BodyJSON(singleToken))
	if err != nil {
		return "", errors.Wrap(err, "PricesCMC")
	}

	if res.Response().StatusCode != 200 {
		return "", errors.Wrap(errors.New("error"), "PricesCMC")
	}

	if err = res.ToJSON(&cmc); err != nil {
		return "", errors.Wrap(err, "PricesCMC")
	}

	doc := cmc.Docs[0]

	return doc.Price, nil
}
