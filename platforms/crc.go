package platforms

import (
	"encoding/json"
	"github.com/button-tech/logger"
	"github.com/button-tech/utils-price-tool/core/currencies"
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"github.com/button-tech/utils-price-tool/types"
	"github.com/imroc/req"
	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
	"strconv"
	"strings"
	"sync"
)

const urlCRC = "https://min-api.cryptocompare.com/data/pricemultifull"

func CrcUpdateWorker(wg *sync.WaitGroup, p *cache.Cache) {
	defer wg.Done()
	SetCRC(p)
}

func SetCRC(c *cache.Cache) {
	var fsyms string
	for value := range c.List {
		fsyms += value + ","
	}

	sortedCurrencies := createCRCRequestData()

	cmc := make(chan map[string][]types.CryptoCompare, len(sortedCurrencies))

	var wg sync.WaitGroup
	wg.Add(len(sortedCurrencies))
	for _, tsyms := range sortedCurrencies {
		go crcPricesRequest(tsyms, fsyms, cmc, &wg, c)
	}
	wg.Wait()

	close(cmc)

	fiatMapping(cmc, c)
}

func crcPricesRequest(tsyms, fsyms string, cmc chan<- map[string][]types.CryptoCompare, wg *sync.WaitGroup, c *cache.Cache) {
	defer wg.Done()

	res, err := req.Get(urlCRC, req.Param{
		"fsyms": fsyms,
		"tsyms": tsyms,
	})
	if err != nil {
		logger.Error("crcPricesRequest", err)
		return
	}

	if res.Response().StatusCode != 200 {
		logger.Error("crcPricesRequest", err)
		return
	}

	result, err := crcFastJson(res.Bytes(), c)
	if err != nil {
		logger.Error("crcPricesRequest", err)
		return
	}

	cmc <- result
}

func crcFastJson(byteRq []byte, c *cache.Cache) (map[string][]types.CryptoCompare, error) {
	var parser fastjson.Parser

	parsed, err := parser.ParseBytes(byteRq)
	if err != nil {
		return nil, errors.Wrap(err, "crcFastJson")
	}

	cryptoCompareDict := make(map[string][]types.CryptoCompare)

	rawObject := parsed.GetObject("RAW")

	rawObject.Visit(func(key []byte, value *fastjson.Value) {
		if obj, ok := c.List[string(key)]; ok {
			crypto := value.GetObject()

			crypto.Visit(func(key []byte, value *fastjson.Value) {

				var c types.CryptoCompare

				if err := json.Unmarshal([]byte(value.String()), &c); err != nil {
					logger.Error("o.Visit", err)
					return
				}

				result, ok := cryptoCompareDict[c.ToSymbol]
				if !ok {
					cryptoCompareDict[c.ToSymbol] = make([]types.CryptoCompare, 0)
				}

				result = append(result, c)

				c.FromSymbol = obj

				result = append(result, c)

				cryptoCompareDict[c.ToSymbol] = result
			})
		}
	})

	return cryptoCompareDict, nil
}

func fiatMapping(cmc chan map[string][]types.CryptoCompare, c *cache.Cache) {
	for {
		m, ok := <-cmc
		if !ok {
			break
		}
		for k, v := range m {
			for _, i := range v {
				details := cache.Details{}
				details.Price = strconv.FormatFloat(i.Price, 'f', -1, 64)
				details.ChangePCT24Hour = strconv.FormatFloat(i.ChangePCT24Hour, 'f', 2, 64)
				details.ChangePCTHour = strconv.FormatFloat(i.ChangePCTHour, 'f', 2, 64)
				k := cache.GenKey("crc", k, i.FromSymbol)
				c.Set(k, details)
			}
		}
	}
}

func createCRCRequestData() []string {
	sortedCurrencies := make([]string, 0, 7)
	n := 0
	step := 25
	for i := 0; i < 6; i++ {
		c := strings.Join(currencies.SupportedCurrenciesList[n:step], ",")
		sortedCurrencies = append(sortedCurrencies, c)
		n += 25
		step += 25
	}
	c := strings.Join(currencies.SupportedCurrenciesList[150:], ",")
	sortedCurrencies = append(sortedCurrencies, c)

	return sortedCurrencies
}
