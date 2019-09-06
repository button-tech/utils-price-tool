package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/utils-tool_prices/storage"
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

	res := cr.store.Get()
	c.JSON(200, gin.H{"res": res.Prices})
}

func (cr *controller) Mount(r *gin.Engine) {
	v1 := r.Group("/api")
	{
		v1.GET("/prices", cr.getCourses)
	}
}
