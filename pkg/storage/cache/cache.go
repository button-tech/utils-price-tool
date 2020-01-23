package cache

import (
	"fmt"
	"github.com/button-tech/utils-price-tool/core/currencies"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
	"sync"
)

// key is "api_usd_currency"
type stored map[string]Details

type Details struct {
	Token           string
	Fiat            string
	Price           string
	ChangePCTHour   string
	ChangePCT24Hour string
	ChangePCT7Day   string
}

type Cache struct {
	sync.RWMutex
	items stored
	List  map[string]string
}

type Key struct {
	API      string
	Fiat     string
	Currency string
}

func NewCache() *Cache {
	return &Cache{
		items: make(stored),
		List:  make(map[string]string),
	}
}

func (c *Cache) Set(k Key, d Details) {
	c.Lock()
	c.items[key(k)] = d
	c.Unlock()
}

func GenKey(a, f, c string) Key {
	return Key{API: strings.ToLower(a), Fiat: strings.ToLower(f), Currency: strings.ToLower(c)}
}

func key(k Key) string {
	return fmt.Sprintf("%s_%s_%s", k.API, k.Fiat, k.Currency)
}

func (c *Cache) Get(k Key) (d Details, ok bool) {
	c.RLock()
	defer c.RUnlock()
	d, ok = c.items[key(k)]
	if !ok {
		return Details{}, false
	}
	return
}

func (c *Cache) Delete(k Key) {
	c.Lock()
	delete(c.items, key(k))
	c.Unlock()
}

func (s *Cache) SetTopList(c map[string]string) error {
	var topList types.PureCoinMarketCap

	urlTopList := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=USD"

	topListAPIKey := os.Getenv("API_KEY")

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
			pricesData := DetailsConversion(
				item.Quote.USD.Price,
				item.Quote.USD.PercentChange1H,
				item.Quote.USD.PercentChange24H,
				item.Quote.USD.PercentChange7D,
			)

			s.Set(GenKey("coinMarketCap", "usd", item.Symbol), pricesData)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(topList.Data))
	for _, v := range topList.Data {
		go func(v types.CmcData, wg *sync.WaitGroup) {
			defer wg.Done()
			if coinID, ok := currencies.TrustV2Coins[v.Symbol]; ok {
				pricesData := DetailsConversion(
					v.Quote.USD.Price,
					v.Quote.USD.PercentChange1H,
					v.Quote.USD.PercentChange24H,
					v.Quote.USD.PercentChange7D,
				)
				convCoinID := strconv.Itoa(coinID)
				kt := GenKey("pcmc", "usd", convCoinID)
				s.Set(kt, pricesData)
			}
		}(v, &wg)
	}
	wg.Wait()

	s.List = topListMap

	return nil
}

func DetailsConversion(price, hour, hour24, sevenDay float64) Details {
	d := Details{Price: strconv.FormatFloat(price, 'f', 10, 64)}
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
