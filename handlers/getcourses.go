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

	//if len(req.Currencies) == 3 {
	//	var price Prices
	//
	//	for _, i := range req.Tokens {
	//		for _, st := range stored {
	//			price.Currency = st.Currency
	//
	//			for _, iSt := range st.Docs {
	//				if iSt.Contract == i {
	//					contract := map[string]string{i:iSt.Price}
	//					price.Rates = append(price.Rates, &contract)
	//				}
	//			}
	//		}
	//	}
	//
	//	result = append(result, price)
	//	c.JSON(200, gin.H{"data": &result})
	//	return
	//}
	//

	// if currencies 3
	//if len(req.Currencies) == 3 {
	//	for _, r := range req.Tokens {
	//		for _, st := range stored {
	//			var price Prices
	//
	//			for _, doc := range st.Docs {
	//				price.Currency = st.Currency
	//
	//				if doc.Contract == r {
	//					contract := map[string]string{r: doc.Price}
	//					price.Rates = append(price.Rates, &contract)
	//				}
	//			}
	//			result = append(result, price)
	//		}
	//	}
	//}
	//c.JSON(200, &result)
	//return

	for _, rc := range req.Tokens {
		for _, st := range stored {
			var price Prices
			price.Currency = st.Currency
			for _, iSt := range st.Docs {
				if iSt.Contract == rc {
					contract := map[string]string{rc: iSt.Price}
					price.Rates = append(price.Rates, &contract)
				}

			}
			result = append(result, price)
		}

	}

	c.JSON(200, gin.H{"data": &result})
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/api/v1/")
	{
		v1.POST("/prices", cr.getCourses)
	}
}
