package handlers

import (
	"errors"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type controller struct {
	store    storage.Storage
	storeCRC storecrc.Storage
	list     storetoplist.Storage
}

func NewController(store storage.Storage, storeCRC storecrc.Storage, list storetoplist.Storage) *controller {
	return &controller{store, storeCRC, list}
}

// data what to get
type dataTokensAndCurrencies struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

// make Response for get prices
type prices struct {
	Currency      string              `json:"currency"`
	Rates         []map[string]string `json:"rates"`
	PercentChange string              `json:"percent_change,omitempty"`
}

// make Response list API
//type listApi struct {
//	API []struct {
//		Name             string   `json:"name"`
//		SupportedChanges []string `json:"supported_changes"`
//		Time             struct {
//			Start int `json:"start"`
//			End   int `json:"end"`
//		} `json:"time"`
//	} `json:"api"`
//}

func (cr *controller) getCourses(c *gin.Context) {
	req := dataTokensAndCurrencies{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	//result := cr.storeCRC.Get()
	//c.JSON(200, gin.H{"data": &res})

	result, err := cr.converter(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "shutdown server"})
		return
	}

	c.JSON(200, gin.H{"data": &result})
}

//func (cr *controller) list(c *gin.Context) {
//	//
//}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/courses/v1/")
	{
		v1.POST("/prices", cr.getCourses)
	}
}

func (cr controller) converter(req *dataTokensAndCurrencies) (*[]prices, error) {
	result := make([]prices, 0)

	switch req.API {
	case "cmc":
		stored := cr.store.Get()

		for _, rq := range req.Currencies {
			price := prices{}

			for _, st := range stored {
				if rq == st.Currency {
					price.Currency = rq

					for _, t := range req.Tokens {
						for _, st := range st.Docs {
							if strings.ToLower(t) == st.Contract {
								contract := map[string]string{t: st.Price}
								price.Rates = append(price.Rates, contract)
							}
						}
					}
				}
			}

			result = append(result, price)
		}

		return &result, nil

	case "crc":
		//res, _ := cr.converterCRCWithChanges(req)
		stored := cr.storeCRC.Get()

		for _, curr := range req.Currencies {
			price := prices{}

			for t, pr := range stored {
				switch curr {
				case pr.USD.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.USD.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}

				case pr.EUR.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.EUR.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}

				case pr.RUB.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.RUB.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}
				}
			}
			result = append(result, price)
		}

		return &result, nil

	default:
		return nil, errors.New("no matches API")
	}
}

func (cr *controller) converterCMCWithChanges(req *dataTokensAndCurrencies) (*[]prices, error) {
	result := make([]prices, 0)

	stored := cr.store.Get()

	for _, rq := range req.Currencies {
		price := prices{}

		for _, st := range stored {
			if rq == st.Currency {
				price.Currency = rq

				for _, t := range req.Tokens {
					for _, st := range st.Docs {
						if strings.ToLower(t) == st.Contract {
							contract := map[string]string{t: st.Price}
							price.Rates = append(price.Rates, contract)
						}
					}
				}
			}
		}

		result = append(result, price)
	}

	return &result, nil

}

func (cr *controller) converterCRCWithChanges(req *dataTokensAndCurrencies) (*[]prices, error) {
	result := make([]prices, 0)

	switch req.Change {
	case "1":
		stored := cr.storeCRC.Get()

		for _, curr := range req.Currencies {
			price := prices{}

			for t, pr := range stored {
				switch curr {
				case pr.USD.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							price.PercentChange = strconv.FormatFloat(pr.USD.CHANGEPCTHOUR, 'f', 2, 64)
							strPrice := strconv.FormatFloat(pr.USD.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}

				case pr.EUR.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							price.PercentChange = strconv.FormatFloat(pr.EUR.CHANGEPCTHOUR, 'f', 2, 64)
							strPrice := strconv.FormatFloat(pr.EUR.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}

				case pr.RUB.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							price.PercentChange = strconv.FormatFloat(pr.RUB.CHANGEPCTHOUR, 'f', 2, 64)
							strPrice := strconv.FormatFloat(pr.RUB.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}
				}
			}
			result = append(result, price)
		}

		return &result, nil

	case "24":
		stored := cr.storeCRC.Get()

		for _, curr := range req.Currencies {
			price := prices{}

			for t, pr := range stored {
				switch curr {
				case pr.USD.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							price.PercentChange = strconv.FormatFloat(pr.USD.CHANGEPCT24HOUR, 'f', 2, 64)
							strPrice := strconv.FormatFloat(pr.USD.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}

				case pr.EUR.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							price.PercentChange = strconv.FormatFloat(pr.EUR.CHANGEPCT24HOUR, 'f', 2, 64)
							strPrice := strconv.FormatFloat(pr.EUR.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}

				case pr.RUB.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						if rqCrypto == t {
							price.PercentChange = strconv.FormatFloat(pr.RUB.CHANGEPCT24HOUR, 'f', 2, 64)
							strPrice := strconv.FormatFloat(pr.RUB.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							price.Rates = append(price.Rates, rate)
						}
					}
				}
			}
			result = append(result, price)
		}

		return &result, nil

	default:
		return nil, errors.New("no matches API")
	}
}
