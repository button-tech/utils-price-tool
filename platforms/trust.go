package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"sync"
)

var trustWalletV2URL = os.Getenv("TRUST_URL_V2")

func TrustUpdateWorker(wg *sync.WaitGroup, p *Prices) {
	defer wg.Done()
	var inWG sync.WaitGroup
	for _, v := range p.TrustV2Coins {
		inWG.Add(1)
		go func(inWg *sync.WaitGroup, price PricesTrustV2) {
			defer inWG.Done()
			if err := p.pricesTrust(price); err != nil {
				logger.Error("trustV2Worker", err)
				return
			}
		}(&inWG, v)
	}
	inWG.Wait()
}

func (p *Prices) pricesTrust(prices PricesTrustV2) error {
	res, err := req.Post(trustWalletV2URL, req.BodyJSON(&prices))

	if err != nil {
		return errors.Wrap(err, "PricesTrustV2")
	}

	if res.Response().StatusCode != 200 {
		return errors.Wrap(errors.New("error"), "PricesTrustV2")
	}

	var trustRes trustV2Response

	if err := res.ToJSON(&trustRes); err != nil {
		return errors.Wrap(err, "PricesTrustV2toJSON")
	}

	trustMapping(&trustRes, p.store)

	return nil
}

func trustMapping(r *trustV2Response, store *cache.Cache) {
	var wg sync.WaitGroup
	wg.Add(len(r.Docs))
	for _, doc := range r.Docs {
		go func(doc trustDoc, wg *sync.WaitGroup) {
			defer wg.Done()
			coin := strconv.Itoa(doc.Coin)
			sd := detailsConversion(doc.Price.Value, 0, doc.Price.Change24H, 0)
			k := cache.GenKey("ntrust", r.Currency, coin)
			store.Set(k, sd)
		}(doc, &wg)
	}
	wg.Wait()
}

func createTrustV2RequestData() []PricesTrustV2 {
	prices := make([]PricesTrustV2, 0, len(currencies))
	for _, c := range currencies {
		price := PricesTrustV2{Currency: c}
		allAssets := make([]Assets, 0, len(TrustV2Coins))
		for _, v := range TrustV2Coins {
			allAssets = append(allAssets, Assets{Coin: v, Type: coin})
		}
		price.Assets = allAssets
		prices = append(prices, price)
	}
	return prices
}
