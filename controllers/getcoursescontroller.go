package controllers

import (
	"errors"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type controller struct {
	store    storage.Cached
}

func NewController(store storage.Cached) *controller {
	return &controller{store: store}
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
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}


// Mapping stored prices with request tokens
// Args have params:
// [0]=currency, [1]=flag - to play with changes


type percentChanges struct {
	Rates map[string]string `json:"rate"`
	PercentChange string `json:"percent_change,omitempty"`
}


// make Response list API
type listApi struct {
	API []api `json:"api"`
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
		result, err := cr.converterCMC(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no matches API changes"})
			return
		}
		c.JSON(200, gin.H{"data": &result})

	case "crc":
		result, err := cr.converterCRC(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no matches API changes"})
			return
		}

		c.JSON(200, gin.H{"data": &result})

	default:
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
		list := listApi{API: API}

		c.JSON(http.StatusBadRequest, gin.H{"error": &list, "description": "please, use these API"})
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
	list := listApi{API: API}

	c.JSON(200, &list)
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/courses/v1/")
	{
		v1.POST("/prices", cr.getCourses)
		v1.GET("/list", cr.apiInfo)
	}
}

func (cr *controller) converterCMC(req *dataTokensAndCurrencies) ([]prices, error) {
	result := make([]prices, 0)
	stored := cr.store.Get()["cmc"]

	switch req.Change {
	case "24":

		for _, c := range req.Currencies {
			price := prices{}

			if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
				price.Currency = c

				for _, t := range req.Tokens {
					if val, ok := fiatVal[storage.CryptoCurrency(strings.ToLower(t))]; ok {
						contract := map[string]string{t: val.Price}
						contract["percent_change"] = val.ChangePCT24Hour
						price.Rates = append(price.Rates, contract)
					}
				}
			}
			result = append(result, price)
		}
		return result, nil

		case "0", "":
		for _, c := range req.Currencies {
			price := prices{}

			if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
				price.Currency = c

				for _, t := range req.Tokens {
					if val, ok := fiatVal[storage.CryptoCurrency(strings.ToLower(t))]; ok {
						contract := map[string]string{t: val.Price}
						price.Rates = append(price.Rates, contract)
					}
				}
			}
			result = append(result, price)
		}
		return result, nil

	default:
		return nil, errors.New("no matches API changes")
	}
}

func (cr *controller) converterCRC(req *dataTokensAndCurrencies) ([]*prices, error) {
	result := make([]*prices, 0)

	switch req.Change {
	case "0", "":
		stored := cr.store.Get()["crc"]

		for _, c := range req.Currencies {
			price := prices{}

			if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
				price.Currency = c

				for _, t := range req.Tokens {
					if val, ok := fiatVal[storage.CryptoCurrency(t)]; ok {
						contract := map[string]string{t: val.Price}
						price.Rates = append(price.Rates, contract)
					}
				}
			}
			result = append(result, &price)
		}
		return result, nil

	case "1":
		stored := cr.store.Get()["crc"]

		for _, c := range req.Currencies {
			price := prices{}

			if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
				price.Currency = c

				for _, t := range req.Tokens {
					if val, ok := fiatVal[storage.CryptoCurrency(t)]; ok {
						contract := map[string]string{t: val.Price}
						contract["percent_change"] = val.ChangePCTHour
						price.Rates = append(price.Rates, contract)
					}
				}
			}
			result = append(result, &price)
		}
		return result, nil

	case "24":
		stored := cr.store.Get()["crc"]

		for _, c := range req.Currencies {
			price := prices{}

			if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
				price.Currency = c

				for _, t := range req.Tokens {
					if val, ok := fiatVal[storage.CryptoCurrency(t)]; ok {
						contract := map[string]string{t: val.Price}
						contract["percent_change"] = val.ChangePCT24Hour
						price.Rates = append(price.Rates, contract)
					}
				}
			}
			result = append(result, &price)
		}
		return result, nil

	default:
		return nil, errors.New("no matches API")
	}
}
