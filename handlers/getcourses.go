package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/utils-price-tool/storage"
	"net/http"
)

type controller struct {
	store storage.Storage
}

func NewController(store storage.Storage) *controller {
	return &controller{store}
}

// data what to get
type dataTokensAndCurrencies struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

// make Response, when no params
type Prices struct {
	Currency string               `json:"currency"`
	Rates    []*map[string]string `json:"rates"`
}

func (cr *controller) getCourses(c *gin.Context) {
	req := dataTokensAndCurrencies{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	var result []Prices
	stored := cr.store.Get()

	for _, rq := range req.Currencies {
		price := Prices{}

		for _, st := range stored {
			if rq == st.Currency {
				price.Currency = rq

				for _, t := range req.Tokens {
					for _, st := range st.Docs {
						if t == st.Contract {
							contract := map[string]string{t:st.Price}
							price.Rates = append(price.Rates, &contract)
						}
					}
				}
			}
		}

		result = append(result, price)
	}


	c.JSON(200, gin.H{"data": &result})
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/api/v1/")
	{
		v1.POST("/prices", cr.getCourses)
	}
}
