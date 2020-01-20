package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/prices"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"sync"
)

var trustWalletV2URL = os.Getenv("TRUST_URL_V2")

func TrustUpdateWorker(wg *sync.WaitGroup, p *prices.PricesData) {
	defer wg.Done()
	var inWG sync.WaitGroup
	for _, v := range p.TrustV2Coins {
		inWG.Add(1)
		go func(inWg *sync.WaitGroup, price types.PricesTrustV2) {
			defer inWG.Done()
			if err := SetPricesTrust(price, p); err != nil {
				logger.Error("trustV2Worker", err)
				return
			}
		}(&inWG, v)
	}
	inWG.Wait()
}

func SetPricesTrust(prices types.PricesTrustV2, p *prices.PricesData) error {
	var wg sync.WaitGroup

	res, err := req.Post(trustWalletV2URL, req.BodyJSON(&prices))

	if err != nil {
		return errors.Wrap(err, "PricesTrustV2")
	}

	if res.Response().StatusCode != 200 {
		return errors.Wrap(errors.New("error"), "PricesTrustV2")
	}

	var trustRes types.TrustV2Response

	if err := res.ToJSON(&trustRes); err != nil {
		return errors.Wrap(err, "PricesTrustV2toJSON")
	}

	wg.Add(len(trustRes.Docs))
	for _, doc := range trustRes.Docs {
		go func(doc types.TrustDoc, wg *sync.WaitGroup) {
			defer wg.Done()
			coin := strconv.Itoa(doc.Coin)
			sd := p.DetailsConversion(doc.Price.Value, 0, doc.Price.Change24H, 0)
			k := cache.GenKey("ntrust", trustRes.Currency, coin)
			p.Store.Set(k, sd)
		}(doc, &wg)
	}
	wg.Wait()

	return nil
}
