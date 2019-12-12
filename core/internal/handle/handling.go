package handle

import (
	"log"
	"sync"

	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/services"
	"github.com/pkg/errors"
)

type Cache interface {
	Get(k string, d cache.Details)
	Set(k string, d cache.Details)
}

var supportedAPIv1 = map[string]struct{}{
	"crc":   {},
	"cmc":   {},
	"huobi": {},
}

var supportedAPIv2 = map[string]struct{}{
	"otrust": {},
	"pcmc":   {},
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

func Reply(u *UniqueData, v string, store *cache.Cache, s *services.GetPrices) ([]Response, error) {
	supportAPIs := chooseVersion(v)
	if _, ok := supportAPIs[u.API]; !ok {
		return nil, errors.New("API: no matches")
	}
	return mapping(u, store, s)
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

func mapping(u *UniqueData, store *cache.Cache, s *services.GetPrices) ([]Response, error) {
	result := make([]Response, 0, len(u.Currencies))

	var wg sync.WaitGroup
	for c := range u.Currencies {
		tokens := services.TokensWithCurrency{Currency: c}
		var price Response
		for t := range u.Tokens {
			k := cache.GenKey(u.API, c, t)
			details, ok := store.Get(k)
			if ok {
				contract := map[string]string{t: details.Price}
				if err := changesControl(contract, &details, u.Change); err != nil {
					return nil, err
				}
				price.Rates = append(price.Rates, contract)
			} else {
				if s != nil {
					tokens.Tokens = append(tokens.Tokens, services.Token{Contract: t})
				}
			}
		}
		if len(price.Rates) > 0 {
			price.Currency = c
			result = append(result, price)
		}
		if s != nil && len(tokens.Tokens) > 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup, store *cache.Cache) {
				if err := s.GetPricesCMC(tokens); err != nil {
					log.Println(err)
				}
				wg.Done()
			}(&wg, store)
		}
	}

	if s == nil || len(result) > 0 {
		return result, nil
	}
	wg.Wait()
	for c := range u.Currencies {
		price := Response{Currency: c}
		for t := range u.Tokens {
			k := cache.GenKey("cmc", c, t)
			details, ok := store.Get(k)
			if ok {
				contract := map[string]string{t: details.Price}
				if err := changesControl(contract, &details, u.Change); err != nil {
					return nil, err
				}
				price.Rates = append(price.Rates, contract)
			}
		}
		if price.Currency != "" {
			result = append(result, price)
		}
	}

	return result, nil
}

func changesControl(m map[string]string, d *cache.Details, c string) (err error) {
	err = errors.New("API changes: no matches")
	switch c {
	case "0", "":
	case "1":
		if d.ChangePCTHour != "" {
			m["percent_change"] = d.ChangePCTHour
		} else {
			return
		}
	case "24":
		if d.ChangePCT24Hour != "" {
			m["percent_change"] = d.ChangePCT24Hour
		} else {
			return
		}
	case "7d":
		if d.ChangePCT7Day != "" {
			m["percent_change"] = d.ChangePCT7Day
		} else {
			return
		}
	default:
		return
	}
	return nil
}
