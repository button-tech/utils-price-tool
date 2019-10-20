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

type topList struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
	} `json:"status"`
	Data []data `json:"data"`
}

type data struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
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
