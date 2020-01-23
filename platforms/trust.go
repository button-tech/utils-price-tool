package platforms

import (
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/currencies"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"sync"
)

var (
	trustWalletV2URL = os.Getenv("TRUST_URL_V2")
	TrustV2Coins     []types.PricesTrustV2
)

func init() {
	TrustV2Coins = make([]types.PricesTrustV2, 0, len(currencies.SupportedCurrenciesList))
	for _, c := range currencies.SupportedCurrenciesList {
		price := types.PricesTrustV2{Currency: c}
		allAssets := make([]types.Assets, 0, len(currencies.TrustV2Coins))
		for _, v := range currencies.TrustV2Coins {
			allAssets = append(allAssets, types.Assets{Coin: v, Type: "coin"})
		}
		price.Assets = allAssets
		TrustV2Coins = append(TrustV2Coins, price)
	}
}

func TrustUpdateWorker(wg *sync.WaitGroup, p *cache.Cache) {
	defer wg.Done()
	var inWG sync.WaitGroup
	for _, v := range TrustV2Coins {
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

func SetPricesTrust(prices types.PricesTrustV2, p *cache.Cache) error {
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
			sd := cache.DetailsConversion(doc.Price.Value, 0, doc.Price.Change24H, 0)
			k := cache.GenKey("ntrust", trustRes.Currency, coin)
			p.Set(k, sd)
		}(doc, &wg)
	}
	wg.Wait()

	return nil
}
