package services

import (
	"time"
)

type TokensWithCurrency struct {
	Currency string  `json:"currency"`
	Tokens   []Token `json:"tokens"`
}

type Token struct {
	Contract string `json:"contract"`
}

// data to get prices to trust-wallet
type TokensWithCurrencies struct {
	Tokens []TokensWithCurrency
}

// DTO for coin-market-cap response
type gotPrices struct {
	Status   bool         `json:"status"`
	Docs     []docsPrices `json:"docs"`
	Currency string       `json:"currency"`
}

type docsPrices struct {
	Price            string `json:"price"`
	Contract         string `json:"contract"`
	PercentChange24H string `json:"percent_change_24h"`
}

// DTO for crypto-compare response
type currency struct {
	ToSymbol        string  `json:"TOSYMBOL"`
	FromSymbol      string  `json:"FROMSYMBOL"`
	Price           float64 `json:"PRICE"`
	ChangePCT24Hour float64 `json:"CHANGEPCT24HOUR"`
	ChangePCTHour   float64 `json:"CHANGEPCTHOUR"`
}

// DTO for top-list response
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
