package controllers

import (
	"errors"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"github.com/button-tech/utils-price-tool/storage/storetoplist"
	"github.com/button-tech/utils-price-tool/storage/storetrustwallet"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type controller struct {
	store    storetrustwallet.Storage
	storeCRC storecrc.Storage
	list     storetoplist.Storage
}

func NewController(store storetrustwallet.Storage, storeCRC storecrc.Storage, list storetoplist.Storage) *controller {
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
	Rates         []percentChanges `json:"rates"`
	//PercentChange string              `json:"percent_change,omitempty"`
}

type percentChanges struct {
	Rates map[string]string
	PercentChange string `json:"percent_change,omitempty"`
}

// make Response list API
type listApi struct {
	API []api  `json:"api"`
		//Time             struct {
		//	Start int `json:"start"`
		//	End   int `json:"end"`
		//} `json:"time"`
}

type api struct {
	Name             string   `json:"name"`
	SupportedChanges []string `json:"supported_changes"`
}

func (cr *controller) getCourses(c *gin.Context) {
	req := dataTokensAndCurrencies{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	switch req.API {
	case "cmc":
		result, err := cr.converterCMCW(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no matches API changes"})
			return
		}
		c.JSON(200, gin.H{"data": &result})
		return

	case "crc":
		result, err := cr.converterCRC(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no matches API changes"})
			return
		}

		c.JSON(200, gin.H{"data": &result})
		return
	}

}

func (cr *controller) apiInfo(c *gin.Context) {
	supportedCRC := []string{"0", "1", "24"}
	crc := api{
		Name:             "crc",
		SupportedChanges: supportedCRC,
	}

	supportedCMC := []string{"24"}
	cmc := api{
		Name:             "cmc",
		SupportedChanges: supportedCMC,
	}

	API := []api{crc, cmc}
	list := listApi{API:API}

	c.JSON(200, &list)
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/courses/v1/")
	{
		v1.POST("/prices", cr.getCourses)
		v1.GET("/list", cr.apiInfo)
	}
}

func (cr *controller) converterCMCW(req *dataTokensAndCurrencies) (*[]prices, error) {
	result := make([]prices, 0)
	stored := cr.store.Get()

	if req.Change == "24" || req.Change == "" {
		for _, rq := range req.Currencies {
			price := prices{}

			for _, st := range stored {
				if rq == st.Currency {
					price.Currency = rq

					for _, t := range req.Tokens {
						for _, st := range st.Docs {
							pctChanges := percentChanges{}

							if strings.ToLower(t) == st.Contract {
								contract := map[string]string{t: st.Price}
								pctChanges.Rates = contract
								pctChanges.PercentChange = st.PercentChange24H
								price.Rates = append(price.Rates, pctChanges)
							}
						}
					}
				}
			}
			result = append(result, price)
		}

		return &result, nil
	}

	return nil, errors.New("no matches API changes")
}


func (cr *controller) converterCRC(req *dataTokensAndCurrencies) (*[]prices, error) {
	result := make([]prices, 0)

	switch req.Change {
	case "0", "":
		stored := cr.storeCRC.Get()

		for _, curr := range req.Currencies {
			price := prices{}

			for t, pr := range stored {
				switch curr {
				case pr.USD.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.USD.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							price.Rates = append(price.Rates, pctChanges)
						}
					}

				case pr.EUR.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.EUR.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							price.Rates = append(price.Rates, pctChanges)
						}
					}

				case pr.RUB.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.RUB.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							price.Rates = append(price.Rates, pctChanges)
						}
					}
				}
			}
			result = append(result, price)
		}

		return &result, nil

	case "1":
		stored := cr.storeCRC.Get()

		for _, curr := range req.Currencies {
			price := prices{}

			for t, pr := range stored {
				switch curr {
				case pr.USD.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.USD.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							pctChanges.PercentChange = strconv.FormatFloat(pr.USD.CHANGEPCTHOUR, 'f', 2, 64)
							price.Rates = append(price.Rates, pctChanges)
						}
					}

				case pr.EUR.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.EUR.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							pctChanges.PercentChange = strconv.FormatFloat(pr.EUR.CHANGEPCTHOUR, 'f', 2, 64)
							price.Rates = append(price.Rates, pctChanges)
						}
					}

				case pr.RUB.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.RUB.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							pctChanges.PercentChange = strconv.FormatFloat(pr.RUB.CHANGEPCTHOUR, 'f', 2, 64)
							price.Rates = append(price.Rates, pctChanges)
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
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.USD.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							pctChanges.PercentChange = strconv.FormatFloat(pr.USD.CHANGEPCT24HOUR, 'f', 2, 64)
							price.Rates = append(price.Rates, pctChanges)
						}
					}

				case pr.EUR.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.EUR.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							pctChanges.PercentChange = strconv.FormatFloat(pr.EUR.CHANGEPCT24HOUR, 'f', 2, 64)
							price.Rates = append(price.Rates, pctChanges)
						}
					}

				case pr.RUB.TOSYMBOL:
					price.Currency = curr
					for _, rqCrypto := range req.Tokens {
						pctChanges := percentChanges{}

						if rqCrypto == t {
							strPrice := strconv.FormatFloat(pr.RUB.PRICE, 'f', 2, 64)
							rate := map[string]string{t: strPrice}
							pctChanges.Rates = rate
							pctChanges.PercentChange = strconv.FormatFloat(pr.RUB.CHANGEPCT24HOUR, 'f', 2, 64)
							price.Rates = append(price.Rates, pctChanges)
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
