package controllers

import (
	"errors"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/gin-gonic/gin"
	"strings"
)

type controller struct {
	store    storage.Cached
}

func NewController(store storage.Cached) *controller {
	return &controller{store: store}
}

type request struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

type response struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

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
	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"err": err})
		return
	}

	a := req.API; switch a {
	case "cmc":
		result, err := cr.converter(&req, a)
		if err != nil {
			c.JSON(400, gin.H{"error": "no matches API changes"})
			return
		}
		c.JSON(200, gin.H{"data": result})

	case "crc":
		result, err := cr.converter(&req, a)
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

func(cr *controller) mapping(req *request, api, ch string) []*response {
	result := make([]*response, 0)
	stored := cr.store.Get()[storage.Api(api)]

	for _, c := range req.Currencies {
		price := response{}

		if fiatVal, fiatOk := stored[storage.Fiat(c)]; fiatOk {
			price.Currency = c

			for _, t := range req.Tokens {
				if val, ok := fiatVal[storage.CryptoCurrency(strings.ToLower(t))]; ok {
					contract := map[string]string{t: val.Price}
					if contract = changesControl(contract, val, ch); contract == nil {
						return nil
					} else {
						price.Rates = append(price.Rates, contract)
					}
				}
			}
		}
		if price.Currency != "" {
			result = append(result, &price)
		}
	}
	return result
}

func changesControl(m map[string]string, s *storage.Details , c string) map[string]string {
	switch c {
	case "1":
		if s.ChangePCTHour != "" {
			m["percent_change"] = s.ChangePCTHour
			return m
		}
		return nil
	case "24":
		if s.ChangePCT24Hour != "" {
			m["percent_change"] = s.ChangePCT24Hour
			return m
		}
		return nil
	default:
		return m
	}
}

func (cr *controller) converter(req *request, api string) ([]*response, error) {
	a := api; switch a {
	case "cmc", "crc":
		resp := cr.switcher(req, a)
		if resp == nil {
			return nil, errors.New("no matches API")
		}
		return resp, nil
	default:
		return nil, errors.New("no matches API")
	}
}

func(cr *controller) switcher(req *request, api string) []*response {
	return cr.mapping(req, api, req.Change)
}
