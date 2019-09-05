package getcourses

import (
	"fmt"
	"github.com/imroc/req"
	"os"
	"time"
)

type TopList struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
	} `json:"status"`
	Data []Data  `json:"data"`
}

type Data struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
	Quote Quote `json:"quote"`
}

type Quote struct {
	USD CurrencyData `json:"USD"`
}

type CurrencyData struct {
	Price            float64   `json:"price"`
	Volume24H        float64   `json:"volume_24h"`
	PercentChange1H  float64   `json:"percent_change_1h"`
	PercentChange24H float64   `json:"percent_change_24h"`
	PercentChange7D  float64   `json:"percent_change_7d"`
	MarketCap        float64   `json:"market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

// get top list limit 100
func(list *TopList) GetTopList() error {
	url := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=6&convert=USD"
	authHeader := req.Header{
		"X-CMC_PRO_API_KEY": os.Getenv("X-CMC_PRO_API_KEY"),
		"Accept": "application/json",
	}
	rq, err := req.Get(url, authHeader)
	if err != nil {
		return fmt.Errorf("can not make a request", err)
	}

	if err = rq.ToJSON(list); err != nil {
		return fmt.Errorf("can not unmarshal")
	}

	return nil
}

func(list *TopList) Convert() error {
	return nil
}