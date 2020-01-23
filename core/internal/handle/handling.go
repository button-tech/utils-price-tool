package handle

import (
	"log"
	"sync"

	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/platforms"
	t "github.com/button-tech/utils-price-tool/types"
	"github.com/pkg/errors"
)

type Data struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

type UniqueData struct {
	Tokens     t.Set
	Currencies t.Set
	Change     string
	API        string
}

type Response struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

type APIs struct {
	Name             string         `json:"name"`
	SupportedChanges []string       `json:"supported_changes"`
	SupportedFiats   map[string]int `json:"supported_fiats"`
}

var supportedAPIv1 = t.Set{
	"crc":   {},
	"cmc":   {},
	"huobi": {},
}

var supportedAPIv2 = t.Set{
	"otrust": {},
	"pcmc":   {},
	"ntrust": {},
}

func Unify(r *Data) UniqueData {
	uniqueTokens := make(t.Set)
	uniqueCurrencies := make(t.Set)

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

func Reply(u *UniqueData, v string, c *cache.Cache) ([]Response, error) {
	supportAPIs := chooseVersion(v)
	if _, ok := supportAPIs[u.API]; !ok {
		return nil, errors.New("API: no matches")
	}

	result := make([]Response, 0, len(u.Currencies))

	var wg sync.WaitGroup
	for currency := range u.Currencies {
		tokens := t.TokensWithCurrency{Currency: currency}
		var response Response

		for token := range u.Tokens {
			key := cache.GenKey(u.API, currency, token)

			details, ok := c.Get(key)
			if ok {
				contract := map[string]string{token: details.Price}
				if err := changesControl(contract, &details, u.Change); err != nil {
					return nil, err
				}
				response.Rates = append(response.Rates, contract)

			} else {
				if c != nil && u.API == "cmc" {
					tokens.Tokens = append(tokens.Tokens, t.Token{Contract: token})
				}
			}
		}

		if len(response.Rates) > 0 {
			response.Currency = currency
			result = append(result, response)
		}
		if c != nil && len(tokens.Tokens) > 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup, store *cache.Cache) {
				if err := platforms.SetCMC(tokens, store); err != nil {
					log.Println(err)
				}
				wg.Done()
			}(&wg, c)
		}
	}
	wg.Wait()

	if c == nil || len(result) > 0 {
		return result, nil
	}

	for currency := range u.Currencies {
		price := Response{Currency: currency}
		for token := range u.Tokens {
			k := cache.GenKey("cmc", currency, token)
			details, ok := c.Get(k)
			if ok {
				contract := map[string]string{token: details.Price}
				if err := changesControl(contract, &details, u.Change); err != nil {
					return nil, err
				}
				price.Rates = append(price.Rates, contract)

				// test variant
				c.Delete(k)
			}
		}
		if price.Currency != "" {
			result = append(result, price)
		}
	}

	return result, nil
}

func unify(wg *sync.WaitGroup, u t.Set, subject []string) {
	for _, s := range subject {
		if _, ok := u[s]; !ok {
			u[s] = struct{}{}
		}
	}
	wg.Done()
}

func chooseVersion(v string) t.Set {
	if v == "v1" {
		return supportedAPIv1
	}
	return supportedAPIv2
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
