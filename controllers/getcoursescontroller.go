package controllers

import (
	"errors"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/gin-gonic/gin"
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
		c.JSON(400, gin.H{"err": err})
		return
	}

	switch req.API {
	case "cmc":
		result, err := cr.converter(&req, "cmc")
		if err != nil {
			c.JSON(400, gin.H{"error": "no matches API changes"})
			return
		}
		c.JSON(200, gin.H{"data": result})

	case "crc":
		result, err := cr.converter(&req, "crc")
		if err != nil {
			c.JSON(400, gin.H{"error": "no matches API changes"})
			return
		}

		c.JSON(200, gin.H{"data": result})

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

		c.JSON(400, gin.H{"error": &list, "description": "please, use these API"})
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

func(cr *controller) mapping(req *dataTokensAndCurrencies, api, ch string) []*prices {
	result := make([]*prices, 0)
	stored := cr.store.Get()[storage.Api(api)]

	for _, c := range req.Currencies {
		price := prices{}

		if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
			price.Currency = c

			for _, t := range req.Tokens {
				if val, ok := fiatVal[storage.CryptoCurrency(t)]; ok {
					contract := map[string]string{t: val.Price}
					price.Rates = append(price.Rates, changesControl(contract, val, ch))
				}
			}
		}
		result = append(result, &price)
	}
	return result
}

func changesControl(m map[string]string, s *storage.Details , c string) map[string]string {
	switch c {
	case "1":
		m["percent_change"] = s.ChangePCTHour
		return m
	case "24":
		m["percent_change"] = s.ChangePCT24Hour
		return m
	default:
		return m
	}
}

func (cr *controller) converter(req *dataTokensAndCurrencies, api string) ([]*prices, error) {
	switch api {
	case "cmc":
		switch req.Change {
		case "24":
			return cr.mapping(req, "cmc", "24"), nil
		case "0", "":
			return cr.mapping(req, "cmc", "0"), nil
		default:
			return nil, errors.New("no matches API changes")
		}

	case "crc":
		switch req.Change {
		case "24":
			return cr.mapping(req, "crc", "24"), nil
		case "1":
			return cr.mapping(req, "crc", "1"), nil
		case "0", "":
			return cr.mapping(req, "crc", "0"), nil
		default:
			return nil, errors.New("no matches API changes")
		}

	default:
		return nil, errors.New("no matches API")
	}
}
