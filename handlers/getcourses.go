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
	res := cr.store.Get()
	for _, rc := req.Currencies {
		for _, st := range res {

		}
	}




	c.JSON(200, gin.H{"res": res})
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/api")
	{
		v1.GET("/prices", cr.getCourses)
	}
}
