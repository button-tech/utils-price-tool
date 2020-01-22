package prices

import (
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"os"
	"strconv"

	"github.com/button-tech/utils-price-tool/core/currencies"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"sync"
)

type PricesData struct {
	TrustV2Coins []types.PricesTrustV2
	List         map[string]string
	Tokens       map[string]string
	Store        *cache.Cache
}

func New(store *cache.Cache) *PricesData {
	prices := make([]types.PricesTrustV2, 0, len(currencies.SupportedCurrenciesList))
	for _, c := range currencies.SupportedCurrenciesList {
		price := types.PricesTrustV2{Currency: c}
		allAssets := make([]types.Assets, 0, len(currencies.TrustV2Coins))
		for _, v := range currencies.TrustV2Coins {
			allAssets = append(allAssets, types.Assets{Coin: v, Type: coin})
		}
		price.Assets = allAssets
		prices = append(prices, price)
	}

	return &PricesData{
		TrustV2Coins: prices,
		List:         make(map[string]string),
		Store:        store,
	}
}

var topListAPIKey = os.Getenv("API_KEY")

const (
	urlTopList = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=USD"
)

const coin = "coin"

// Get top List of crypto-currencies from coin-market
func (p *PricesData) SetTopList(c map[string]string) error {
	var topList types.PureCoinMarketCap

	res, err := req.Get(urlTopList, req.Header{"X-CMC_PRO_API_KEY": topListAPIKey})
	if err != nil {
		return errors.Wrap(err, "getTopList")
	}

	if res.Response().StatusCode != 200 {
		return errors.Wrap(errors.New("error"), "SetTopList")
	}

	if err = res.ToJSON(&topList); err != nil {
		return errors.Wrap(err, "SetTopList")
	}

	if topList.Status.ErrorCode != 0 {
		return errors.New("responseHTTPStatus: NotOk")
	}

	top100 := topList.Data[:100]
	topListMap := make(map[string]string)
	for _, item := range top100 {
		if val, ok := c[item.Symbol]; ok {
			topListMap[item.Symbol] = val
			pricesData := p.DetailsConversion(
				item.Quote.USD.Price,
				item.Quote.USD.PercentChange1H,
				item.Quote.USD.PercentChange24H,
				item.Quote.USD.PercentChange7D,
			)

			p.Store.Set(cache.GenKey("coinMarketCap", "usd", item.Symbol), pricesData)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(topList.Data))
	for _, v := range topList.Data {
		go func(v types.CmcData, wg *sync.WaitGroup) {
			defer wg.Done()
			if coinID, ok := currencies.TrustV2Coins[v.Symbol]; ok {
				pricesData := p.DetailsConversion(
					v.Quote.USD.Price,
					v.Quote.USD.PercentChange1H,
					v.Quote.USD.PercentChange24H,
					v.Quote.USD.PercentChange7D,
				)
				convCoinID := strconv.Itoa(coinID)
				kt := cache.GenKey("pcmc", "usd", convCoinID)
				p.Store.Set(kt, pricesData)
			}
		}(v, &wg)
	}
	wg.Wait()

	p.List = topListMap

	return nil
}

func (_ *PricesData) DetailsConversion(price, hour, hour24, sevenDay float64) cache.Details {
	d := cache.Details{Price: strconv.FormatFloat(price, 'f', 10, 64)}
	if hour != 0 {
		d.ChangePCTHour = strconv.FormatFloat(hour, 'f', 6, 64)
	}
	if hour24 != 0 {
		d.ChangePCT24Hour = strconv.FormatFloat(hour24, 'f', 6, 64)
	}
	if sevenDay != 0 {
		d.ChangePCT7Day = strconv.FormatFloat(sevenDay, 'f', 6, 64)
	}
	return d
}
