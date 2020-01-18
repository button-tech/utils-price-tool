package services

import (
	"time"
)

type RequestCoinMarketCap struct {
	Tokens []TokensWithCurrency
}

type TokensWithCurrency struct {
	Currency string  `json:"currency"`
	Tokens   []Token `json:"tokens"`
}

type Token struct {
	Contract string `json:"contract"`
}

type coinMarketCap struct {
	Status   bool         `json:"status"`
	Docs     []docsPrices `json:"docs"`
	Currency string       `json:"currency"`
}

type docsPrices struct {
	Price            string `json:"price"`
	Contract         string `json:"contract"`
	PercentChange24H string `json:"percent_change_24h"`
}

type cryptoCompare struct {
	ToSymbol        string  `json:"TOSYMBOL"`
	FromSymbol      string  `json:"FROMSYMBOL"`
	Price           float64 `json:"PRICE"`
	ChangePCT24Hour float64 `json:"CHANGEPCT24HOUR"`
	ChangePCTHour   float64 `json:"CHANGEPCTHOUR"`
}

type huobi struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol     string  `json:"symbol"`
		IndexPrice float64 `json:"index_price"`
		IndexTs    int64   `json:"index_ts"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type pureCoinMarketCap struct {
	Status struct {
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
	} `json:"status"`
	Data []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  struct {
			USD struct {
				Price            float64 `json:"price"`
				PercentChange1H  float64 `json:"percent_change_1h"`
				PercentChange24H float64 `json:"percent_change_24h"`
				PercentChange7D  float64 `json:"percent_change_7d"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

type PricesTrustV2 struct {
	Currency string   `json:"currency"`
	Assets   []Assets `json:"Assets"`
}

type Assets struct {
	Coin    int    `json:"coin"`
	Type    string `json:"type"`
	TokenID string `json:"token_id"`
}

type trustV2Response struct {
	Currency string `json:"currency"`
	Docs     []struct {
		Coin  int    `json:"coin"`
		Type  string `json:"type"`
		Price struct {
			Value     float64 `json:"value"`
			Change24H float64 `json:"change_24h"`
		} `json:"price"`
		LastUpdate time.Time `json:"last_update"`
	} `json:"docs"`
}
