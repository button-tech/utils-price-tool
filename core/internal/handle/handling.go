package handle

import (
	"log"
	"sync"

	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/platforms"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
)

type Cache interface {
	Get(k string, d cache.Details)
	Set(k string, d cache.Details)
}

type Data struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

type UniqueData struct {
	Tokens     map[string]struct{}
	Currencies map[string]struct{}
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

func Reply(u *UniqueData, v string, store *cache.Cache, s *cache.Cache) ([]Response, error) {
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

func mapping(u *UniqueData, store *cache.Cache, p *cache.Cache) ([]Response, error) {
	result := make([]Response, 0, len(u.Currencies))

	var wg sync.WaitGroup
	for c := range u.Currencies {
		tokens := types.TokensWithCurrency{Currency: c}
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
				if p != nil && u.API == "cmc" {
					tokens.Tokens = append(tokens.Tokens, types.Token{Contract: t})
				}
			}
		}
		if len(price.Rates) > 0 {
			price.Currency = c
			result = append(result, price)
		}
		if p != nil && len(tokens.Tokens) > 0 {
			wg.Add(1)
			go func(wg *sync.WaitGroup, store *cache.Cache) {
				if err := platforms.SetPricesCMC(tokens, p); err != nil {
					log.Println(err)
				}
				wg.Done()
			}(&wg, store)
		}
	}

	if p == nil || len(result) > 0 {
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

				// test variant
				store.Delete(k)
			}
		}
		if price.Currency != "" {
			result = append(result, price)
		}
	}

	return result, nil
}

func SingleERC20Course(fiat, crypto string, s *cache.Cache) (string, error) {

	var cmc types.CmcResponse

	token := make([]types.Token, 0, 1)
	token = append(token, types.Token{Contract: crypto})

	singleToken := types.TokensWithCurrency{
		Currency: fiat,
		Tokens:   token,
	}

	res, err := req.Post(platforms.TrustWalletURL, req.BodyJSON(singleToken))
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
