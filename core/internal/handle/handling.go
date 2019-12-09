package handle

import (
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"strings"
	"sync"

	"github.com/button-tech/utils-price-tool/pkg/typeconv"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/pkg/errors"
)

const trust = "cmc"

type Cache interface {
	Get(a cache.Api) (f cache.FiatMap, err error)
	Set(a cache.Api, f cache.FiatMap)
}

var supportedAPIv1 = map[string]struct{}{
	"crc":   {},
	"cmc":   {},
	"huobi": {},
}

var supportedAPIv2 = map[string]struct{}{
	"otrust": {},
	"cmc":    {},
	"ntrust": {},
}

func Unify(r *Data) UniqueData {
	uniqueTokens := make(map[string]struct{})
	uniqueCurrencies := make(map[string]struct{})

	var wg sync.WaitGroup
	wg.Add(2)
	go unify(&wg, uniqueTokens, r.Tokens)
	go unify(&wg, uniqueCurrencies, r.Currencies)
	wg.Wait()
	return UniqueData{
		Tokens:     uniqueTokens,
		Currencies: uniqueCurrencies,
		Change:     r.Change,
		API:        r.API,
	}
}

func Reply(u *UniqueData, v string, f Cache, s *services.Service) ([]response, error) {
	supportAPIs := chooseVersion(v)
	if _, ok := supportAPIs[u.API]; !ok {
		return nil, errors.New("API: no matches")
	}
	return mapping(u, f, s)
}

func unify(wg *sync.WaitGroup, u map[string]struct{}, subject []string) {
	for _, s := range subject {
		if _, ok := u[s]; !ok {
			u[s] = struct{}{}
		}
	}
	wg.Done()
}

func chooseVersion(v string) map[string]struct{} {
	if v == "v1" {
		return supportedAPIv1
	}
	return supportedAPIv2
}

func mapping(u *UniqueData, c Cache, s *services.Service) ([]response, error) {
	result := make([]response, 0, len(u.Currencies))
	api := u.API
	stored, err := c.Get(typeconv.StorageApi(api))
	if err != nil {
		return nil, err
	}

	var notExistTokens []services.TokensWithCurrency
	for c := range u.Currencies {
		tokens := services.TokensWithCurrency{Currency: c}
		price := response{}
		if fiatVal, fiatOk := stored[typeconv.StorageFiat(c)]; fiatOk {
			price.Currency = c
			for t := range u.Tokens {
				currency := storageCC(u.API, t)
				if details, ok := fiatVal[currency]; ok {
					contract := map[string]string{t: details.Price}
					if err := changesControl(contract, details, u.Change); err != nil {
						return nil, err
					}
					price.Rates = append(price.Rates, contract)
				} else {
					if s != nil {
						tokens.Tokens = append(tokens.Tokens, services.Token{Contract: t})
					}
				}
			}
		}

		if price.Currency != "" {
			result = append(result, price)
		}
		if s != nil {
			notExistTokens = append(notExistTokens, tokens)
		}
	}
	if s == nil {
		return result, nil
	}

	var wg sync.WaitGroup
	wg.Add(len(notExistTokens))
	notExistTokensChan := make(chan cache.FiatMap, len(notExistTokens))
	for _, t := range notExistTokens {
		go func(wg *sync.WaitGroup, t services.TokensWithCurrency, c chan cache.FiatMap) {
			f, err := s.GetPricesCMC(t)
			if err != nil {

			}
			// cache.Set("cmcTokens", f)
			c <- f
			wg.Done()
		}(&wg, t, notExistTokensChan)
	}
	wg.Wait()
	close(notExistTokensChan)

	for f := range notExistTokensChan {
		for i, v := range result {
			if ccMap, ok := f[typeconv.StorageFiat(v.Currency)]; ok {
				for c, d := range ccMap {
					contract := map[string]string{string(c): d.Price}
					if err := changesControl(contract, d, u.Change); err != nil {
						return nil, err
					}
					result[i].Rates = append(result[i].Rates, contract)
				}
			}
		}
	}

	return result, nil
}

func changesControl(m map[string]string, d *cache.Details, c string) error {
	switch c {
	case "0", "":
	case "1":
		if d.ChangePCTHour != "" {
			m["percent_change"] = d.ChangePCTHour
		}
	case "24":
		if d.ChangePCT24Hour != "" {
			m["percent_change"] = d.ChangePCT24Hour
		}
	default:
		return errors.New("API changes: no matches")
	}
	return nil
}

func storageCC(api, t string) (c cache.CryptoCurrency) {
	if api == "ntrust" {
		c = typeconv.StorageCC(t)
		return
	}
	c = typeconv.StorageCC(strings.ToLower(t))
	return
}
