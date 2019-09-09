package storetoplist

import "time"

type Top10List struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
	} `json:"status"`
	Data []struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         int         `json:"max_supply"`
		CirculatingSupply int         `json:"circulating_supply"`
		TotalSupply       int         `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             USD         `json:"quote"`
	} `json:"data"`
}

//type Currencies struct {
//	USD USD `json:"USD"`
//	EUR EUR `json:"EUR"`
//	RUB RUB `json:"RUB"`
//}
//
//type RUB struct {
//	Price            float64   `json:"price"`
//	Volume24H        float64   `json:"volume_24h"`
//	PercentChange1H  float64   `json:"percent_change_1h"`
//	PercentChange24H float64   `json:"percent_change_24h"`
//	PercentChange7D  float64   `json:"percent_change_7d"`
//	MarketCap        float64   `json:"market_cap"`
//	LastUpdated      time.Time `json:"last_updated"`
//}

type USD struct {
	Price            float64   `json:"price"`
	Volume24H        float64   `json:"volume_24h"`
	PercentChange1H  float64   `json:"percent_change_1h"`
	PercentChange24H float64   `json:"percent_change_24h"`
	PercentChange7D  float64   `json:"percent_change_7d"`
	MarketCap        float64   `json:"market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

//type EUR struct {
//	Price            float64   `json:"price"`
//	Volume24H        float64   `json:"volume_24h"`
//	PercentChange1H  float64   `json:"percent_change_1h"`
//	PercentChange24H float64   `json:"percent_change_24h"`
//	PercentChange7D  float64   `json:"percent_change_7d"`
//	MarketCap        float64   `json:"market_cap"`
//	LastUpdated      time.Time `json:"last_updated"`
//}
//
//
//type TopListRUB struct {
//	Status struct {
//		Timestamp    time.Time   `json:"timestamp"`
//		ErrorCode    int         `json:"error_code"`
//		ErrorMessage interface{} `json:"error_message"`
//		Elapsed      int         `json:"elapsed"`
//		CreditCount  int         `json:"credit_count"`
//	} `json:"status"`
//	Data []struct {
//		ID                int         `json:"id"`
//		Name              string      `json:"name"`
//		Symbol            string      `json:"symbol"`
//		Slug              string      `json:"slug"`
//		NumMarketPairs    int         `json:"num_market_pairs"`
//		DateAdded         time.Time   `json:"date_added"`
//		Tags              []string    `json:"tags"`
//		MaxSupply         int         `json:"max_supply"`
//		CirculatingSupply int         `json:"circulating_supply"`
//		TotalSupply       int         `json:"total_supply"`
//		Platform          interface{} `json:"platform"`
//		CmcRank           int         `json:"cmc_rank"`
//		LastUpdated       time.Time   `json:"last_updated"`
//		Quote             struct {
//			RUB RUB `json:"RUB"`
//		} `json:"quote"`
//	} `json:"data"`
//}
//
//type TopListEUR struct {
//	Status struct {
//		Timestamp    time.Time   `json:"timestamp"`
//		ErrorCode    int         `json:"error_code"`
//		ErrorMessage interface{} `json:"error_message"`
//		Elapsed      int         `json:"elapsed"`
//		CreditCount  int         `json:"credit_count"`
//	} `json:"status"`
//	Data []struct {
//		ID                int         `json:"id"`
//		Name              string      `json:"name"`
//		Symbol            string      `json:"symbol"`
//		Slug              string      `json:"slug"`
//		NumMarketPairs    int         `json:"num_market_pairs"`
//		DateAdded         time.Time   `json:"date_added"`
//		Tags              []string    `json:"tags"`
//		MaxSupply         int         `json:"max_supply"`
//		CirculatingSupply int         `json:"circulating_supply"`
//		TotalSupply       int         `json:"total_supply"`
//		Platform          interface{} `json:"platform"`
//		CmcRank           int         `json:"cmc_rank"`
//		LastUpdated       time.Time   `json:"last_updated"`
//		Quote             struct {
//			EUR EUR `json:"EUR"`
//		} `json:"quote"`
//	} `json:"data"`
//}
//
//type TopListUSD struct {
//	Status struct {
//		Timestamp    time.Time   `json:"timestamp"`
//		ErrorCode    int         `json:"error_code"`
//		ErrorMessage interface{} `json:"error_message"`
//		Elapsed      int         `json:"elapsed"`
//		CreditCount  int         `json:"credit_count"`
//	} `json:"status"`
//	Data []struct {
//		ID                int         `json:"id"`
//		Name              string      `json:"name"`
//		Symbol            string      `json:"symbol"`
//		Slug              string      `json:"slug"`
//		NumMarketPairs    int         `json:"num_market_pairs"`
//		DateAdded         time.Time   `json:"date_added"`
//		Tags              []string    `json:"tags"`
//		MaxSupply         int         `json:"max_supply"`
//		CirculatingSupply int         `json:"circulating_supply"`
//		TotalSupply       int         `json:"total_supply"`
//		Platform          interface{} `json:"platform"`
//		CmcRank           int         `json:"cmc_rank"`
//		LastUpdated       time.Time   `json:"last_updated"`
//		Quote             struct {
//			USD USD `json:"USD"`
//		} `json:"quote"`
//	} `json:"data"`
//}
