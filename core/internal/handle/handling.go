package handle

import (
	"strings"
	"sync"

	"github.com/button-tech/utils-price-tool/pkg/storage"
	"github.com/button-tech/utils-price-tool/pkg/typeconv"
	"github.com/pkg/errors"
)

type FiatMap interface {
	Get(a storage.Api) (f storage.FiatMap, err error)
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

func Reply(u *UniqueData, v string, f FiatMap) ([]response, error) {
	supportAPIs := chooseVersion(v)
	if _, ok := supportAPIs[u.API]; !ok {
		return nil, errors.New("API: no matches")
	}
	return mapping(u, f)
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

func mapping(u *UniqueData, f FiatMap) ([]response, error) {
	result := make([]response, 0, len(u.Currencies))
	api := u.API
	stored, err := f.Get(typeconv.StorageApi(api))
	if err != nil {
		return nil, err
	}
	for c := range u.Currencies {
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
				}
			}
		}
		if price.Currency != "" {
			result = append(result, price)
		}
	}
	return result, nil
}

func changesControl(m map[string]string, d *storage.Details, c string) error {
	switch c {
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

func storageCC(api, t string) (c storage.CryptoCurrency) {
	if api == "ntrust" {
		c = typeconv.StorageCC(t)
		return
	}
	c = typeconv.StorageCC(strings.ToLower(t))
	return
}
