package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/button-tech/utils-price-tool/storage"
	"github.com/button-tech/utils-price-tool/storage/storecrc"
	"net/http"
	"strings"
)

type controller struct {
	store    storage.Storage
	storeCRC storecrc.Storage
}

func NewController(store storage.Storage, storeCRC storecrc.Storage) *controller {
	return &controller{store, storeCRC}
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
	Currency      string               `json:"currency"`
	Rates         []*map[string]string `json:"rates"`
	PercentChange string               `json:"percent_change,omitempty"`
}

// make Response list API
type listApi struct {
	API []struct {
		Name             string   `json:"name"`
		SupportedChanges []string `json:"supported_changes"`
		Time             struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"time"`
	} `json:"api"`
}

func (cr *controller) getCourses(c *gin.Context) {
	req := dataTokensAndCurrencies{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	result := cr.converter(&req)
	storeCRC := cr.storeCRC.Get()

	c.JSON(200, gin.H{"data": &result, "CRC": storeCRC})
}

func (cr *controller) list(c *gin.Context) {
	//
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/api/v1/")
	{
		v1.POST("/prices", cr.getCourses)
	}
}

func (cr *controller) converter(req *dataTokensAndCurrencies) *[]prices {
	var result []prices

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
							price.Rates = append(price.Rates, &contract)
						}
					}
				}
			}
		}

		result = append(result, price)
	}

	return &result
}
