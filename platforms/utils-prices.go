package platforms

import (
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
	"os"
	"strconv"

	"github.com/imroc/req"
	"github.com/pkg/errors"
)

type Prices struct {
	TrustV2Coins []PricesTrustV2
	List         map[string]string
	Tokens       map[string]string
	store        *cache.Cache
}

func NewPrices(store *cache.Cache) *Prices {
	return &Prices{
		TrustV2Coins: createTrustV2RequestData(),
		List:         make(map[string]string),
		store:        store,
	}
}

var topListAPIKey = os.Getenv("API_KEY")

const (
	urlTopList = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=USD"
)

const coin = "coin"

// Get top List of crypto-currencies from coin-market
func (p *Prices) GetTopList(c map[string]string) error {
	var topList pureCoinMarketCap

	res, err := req.Get(urlTopList, req.Header{"X-CMC_PRO_API_KEY": topListAPIKey})
	if err != nil {
		return errors.Wrap(err, "getTopList")
	}

	if res.Response().StatusCode != 200 {
		return errors.Wrap(errors.New("error"), "getTopList")
	}

	if err = res.ToJSON(&topList); err != nil {
		return errors.Wrap(err, "getTopList")
	}

	if topList.Status.ErrorCode != 0 {
		return errors.New("responseHTTPStatus: NotOk")
	}

	top100 := topList.Data[:100]
	topListMap := make(map[string]string)
	for _, item := range top100 {
		if val, ok := c[item.Symbol]; ok {
			topListMap[item.Symbol] = val
			pricesData := detailsConversion(
				item.Quote.USD.Price,
				item.Quote.USD.PercentChange1H,
				item.Quote.USD.PercentChange24H,
				item.Quote.USD.PercentChange7D,
			)

			p.store.Set(cache.GenKey("coinMarketCap", "usd", item.Symbol), pricesData)
		}
	}

	cmcMapping(topList, p.store)

	p.List = topListMap

	return nil
}

func detailsConversion(price, hour, hour24, sevenDay float64) cache.Details {
	d := cache.Details{Price: strconv.FormatFloat(price, 'f', 10, 64)}
	if floatValid(hour) {
		d.ChangePCTHour = strconv.FormatFloat(hour, 'f', 6, 64)
	}
	if floatValid(hour24) {
		d.ChangePCT24Hour = strconv.FormatFloat(hour24, 'f', 6, 64)
	}
	if floatValid(sevenDay) {
		d.ChangePCT7Day = strconv.FormatFloat(sevenDay, 'f', 6, 64)
	}
	return d
}

func floatValid(s float64) bool {
	return s != 0
}

var PureCMCCoins = map[string]int{
	"AE":    457,
	"ALGO":  283,
	"ATOM":  118,
	"BCH":   145,
	"BNB":   714,
	"BTC":   0,
	"DASH":  5,
	"DCR":   42,
	"DGB":   20,
	"DOGE":  3,
	"ETC":   61,
	"ETH":   60,
	"ICX":   74,
	"LTC":   2,
	"NANO":  165,
	"ONT":   1024,
	"QTUM":  2301,
	"RVN":   175,
	"THETA": 500,
	"TRX":   195,
	"VET":   818,
	"WAVES": 5741564,
	"XLM":   148,
	"XRP":   144,
	"XTZ":   1729,
	"ZEC":   133,
	"ZIL":   313,
}

var TrustV2Coins = map[string]int{
	"ETH":   60,
	"ETC":   61,
	"ICX":   74,
	"ATOM":  118,
	"XRP":   144,
	"XLM":   148,
	"POA":   178,
	"TRX":   195,
	"FIO":   235,
	"NIM":   242,
	"IOTX":  304,
	"ZIL":   313,
	"AION":  425,
	"AE":    457,
	"THETA": 500,
	"BNB":   714,
	"VET":   818,
	"CLO":   820,
	"TOMO":  889,
	"TT":    1001,
	"ONT":   1024,
	"XTZ":   1729,
	"KIN":   2017,
	"NAS":   2718,
	"GO":    6060,
	"WAN":   5718350,
	"WAVES": 5741564,
	"SEM":   7562605,
	"BTC":   0,
	"LTC":   2,
	"DOGE":  3,
	"DASH":  5,
	"VIA":   14,
	"GRS":   17,
	"ZEC":   133,
	"XZC":   136,
	"BCH":   145,
	"RVN":   175,
	"QTUM":  2301,
	"ZEL":   19167,
	"DCR":   42,
	"ALGO":  283,
	"NANO":  165,
	"DGB":   20,
}

var currencies = []string{
	"AED",
	"AFN",
	"ALL",
	"AMD",
	"ANG",
	"AOA",
	"ARS",
	"AUD",
	"AWG",
	"AZN",
	"BAM",
	"BBD",
	"BDT",
	"BGN",
	"BHD",
	"BIF",
	"BMD",
	"BND",
	"BOB",
	"BRL",
	"BSD",
	"BTC",
	"BTN",
	"BWP",
	"BYN",
	"BYR",
	"BZD",
	"CAD",
	"CDF",
	"CHF",
	"CLF",
	"CLP",
	"CNY",
	"COP",
	"CRC",
	"CUC",
	"CUP",
	"CVE",
	"CZK",
	"DJF",
	"DKK",
	"DOP",
	"DZD",
	"EGP",
	"ERN",
	"ETB",
	"EUR",
	"FJD",
	"FKP",
	"GBP",
	"GEL",
	"GGP",
	"GHS",
	"GIP",
	"GMD",
	"GNF",
	"GTQ",
	"GYD",
	"HKD",
	"HNL",
	"HRK",
	"HTG",
	"HUF",
	"IDR",
	"ILS",
	"IMP",
	"INR",
	"IQD",
	"IRR",
	"ISK",
	"JEP",
	"JMD",
	"JOD",
	"JPY",
	"KES",
	"KGS",
	"KHR",
	"KMF",
	"KPW",
	"KRW",
	"KWD",
	"KYD",
	"KZT",
	"LAK",
	"LBP",
	"LKR",
	"LRD",
	"LSL",
	"LTL",
	"LVL",
	"LYD",
	"MAD",
	"MDL",
	"MGA",
	"MKD",
	"MMK",
	"MNT",
	"MOP",
	"MRO",
	"MUR",
	"MVR",
	"MWK",
	"MXN",
	"MYR",
	"MZN",
	"NAD",
	"NGN",
	"NIO",
	"NOK",
	"NPR",
	"NZD",
	"OMR",
	"PAB",
	"PEN",
	"PGK",
	"PHP",
	"PKR",
	"PLN",
	"PYG",
	"QAR",
	"RON",
	"RUB",
	"RWF",
	"SAR",
	"SBD",
	"SCR",
	"SDG",
	"SEK",
	"SGD",
	"SHP",
	"SLL",
	"SOS",
	"SRD",
	"STD",
	"SVC",
	"SYP",
	"SZL",
	"THB",
	"TJS",
	"TMT",
	"TND",
	"TOP",
	"TRY",
	"TTD",
	"TWD",
	"TZS",
	"UAH",
	"UGX",
	"USD",
	"UYU",
	"UZS",
	"VEF",
	"VND",
	"VUV",
	"WST",
	"XAF",
	"XAG",
	"XAU",
	"XCD",
	"XDR",
	"XOF",
	"XPF",
	"YER",
	"ZAR",
	"ZMK",
	"ZMW",
	"ZWL",
}
