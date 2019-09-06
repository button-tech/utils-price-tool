package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/utils-tool_prices/services"
	"github.com/utils-tool_prices/storage"
	"log"
	"sync"
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

func (cr *controller) getCourses(c *gin.Context) {
	//resp := dataTokensAndCurrencies{}
	//if err := c.ShouldBindJSON(&resp); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"err": err})
	//	return
	//}

	tokens := services.InitRequestData()
	wg := sync.WaitGroup{}

	var stors storage.ResultPrices

	for _, t := range tokens.Tokens {
		wg.Add(1)

		go func(wg *sync.WaitGroup, st *storage.ResultPrices, t *services.TokensWithCurrency) {
			defer wg.Done()

			got, err := services.GetPricesCMC(t)
			if err != nil {
				log.Println(err)
				return
			}

			var store storage.Prices
			for _, i := range got.Docs {
				store.Currency = got.Currency
				contract := map[string]string{i.Contract: i.Price}
				store.Rates = append(store.Rates, &contract)
			}

			st.Update(store)

		}(&wg, &stors, &t)
		wg.Wait()
	}


	c.JSON(200, gin.H{"res": stors.Prices})
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/api")
	{
		v1.GET("/prices", cr.getCourses)
	}
}
